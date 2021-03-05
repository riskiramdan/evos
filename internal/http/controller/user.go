package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/riskiramdan/evos/internal/data"
	"github.com/riskiramdan/evos/internal/http/response"
	"github.com/riskiramdan/evos/internal/types"
	"github.com/riskiramdan/evos/internal/user"
	u "github.com/riskiramdan/evos/util"
)

// UserController represents the user controller
type UserController struct {
	userService user.ServiceInterface
	dataManager *data.Manager
	utility     *u.Utility
}

// UserList user list and count
type UserList struct {
	Data  []*user.Users `json:"data"`
	Total int           `json:"total"`
}

// GetListUser function for get list data users
func (a *UserController) GetListUser(w http.ResponseWriter, r *http.Request) {
	var err *types.Error

	queryValues := r.URL.Query()
	var limit = 10
	var errConversion error
	if queryValues.Get("limit") != "" {
		limit, errConversion = strconv.Atoi(queryValues.Get("limit"))
		if errConversion != nil {
			err = &types.Error{
				Path:    ".UserController->ListUser()",
				Message: errConversion.Error(),
				Error:   errConversion,
				Type:    "golang-error",
			}
			response.Error(w, "Bad Request", http.StatusBadRequest, *err)
			return
		}
	}

	var page = 1
	if queryValues.Get("page") != "" {
		page, errConversion = strconv.Atoi(queryValues.Get("page"))
		if errConversion != nil {
			err = &types.Error{
				Path:    ".UserController->ListUser()",
				Message: errConversion.Error(),
				Error:   errConversion,
				Type:    "golang-error",
			}
			response.Error(w, "Bad Request", http.StatusBadRequest, *err)
			return
		}
	}

	var name = queryValues.Get("name")
	var phone = queryValues.Get("phone")

	if limit < 0 {
		limit = 10
	}
	if page < 0 {
		page = 1
	}
	userList, count, err := a.userService.ListUsers(r.Context(), &user.FindAllUsersParams{
		Name:  name,
		Phone: phone,
		Limit: limit,
		Page:  page,
	})
	if err != nil {
		err.Path = ".UserController->ListUser()" + err.Path
		if err.Error != data.ErrNotFound {
			response.Error(w, "Internal Server Error", http.StatusInternalServerError, *err)
			return
		}
	}
	if userList == nil {
		userList = []*user.Users{}
	}

	response.JSON(w, http.StatusOK, UserList{
		Data:  userList,
		Total: count,
	})
}

// PostCreateUser for creating data user
func (a *UserController) PostCreateUser(w http.ResponseWriter, r *http.Request) {
	var err *types.Error

	decoder := json.NewDecoder(r.Body)

	var params *user.TransactionParams
	errDecode := decoder.Decode(&params)
	if errDecode != nil {
		err = &types.Error{
			Path:    ".UserController->CreateUser()",
			Message: errDecode.Error(),
			Error:   errDecode,
			Type:    "golang-error",
		}
		response.Error(w, "Bad Request", http.StatusBadRequest, *err)
		return
	}
	resp := &user.Users{}
	errTransaction := a.dataManager.RunInTransaction(r.Context(), func(ctx context.Context) error {
		resp, err = a.userService.CreateUser(ctx, &user.TransactionParams{
			Name:     params.Name,
			RoleID:   params.RoleID,
			Phone:    params.Phone,
			Password: a.utility.RandStringBytesMaskImprSrcSB(4),
		})
		if err != nil {
			return err.Error
		}
		return nil
	})
	if errTransaction != nil {
		err.Path = ".UserController->CreateUser()" + err.Path
		if errTransaction == user.ErrPhoneAlreadyExists {
			response.Error(w, user.ErrPhoneAlreadyExists.Error(), http.StatusUnprocessableEntity, *err)
		} else {
			response.Error(w, "Internal Server Error", http.StatusInternalServerError, *err)
		}
		return
	}

	response.JSON(w, http.StatusOK, resp)
}

// PostLogin for getting authorization ..
func (a *UserController) PostLogin(w http.ResponseWriter, r *http.Request) {
	var err *types.Error

	decoder := json.NewDecoder(r.Body)

	var params user.LoginParams
	errDecode := decoder.Decode(&params)
	if errDecode != nil {
		err = &types.Error{
			Path:    ".UserController->Login()",
			Message: errDecode.Error(),
			Error:   errDecode,
			Type:    "golang-error",
		}
		response.Error(w, "Bad Request", http.StatusBadRequest, *err)
		return
	}

	var sess *user.LoginResponse
	errTransaction := a.dataManager.RunInTransaction(r.Context(), func(ctx context.Context) error {
		sess, err = a.userService.Login(r.Context(), params.Phone, params.Password)
		if err != nil {
			return err.Error
		}
		return nil
	})
	if errTransaction != nil {
		err.Path = ".UserController->Login()" + err.Path
		if err.Error == user.ErrWrongPassword || err.Error == data.ErrNotFound || err.Error == user.ErrWrongPhone {
			response.Error(w, "Phone / password is wrong", http.StatusBadRequest, *err)
		} else {
			response.Error(w, "Internal Server Error", http.StatusInternalServerError, *err)
		}
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "sessionId",
		Value: sess.SessionID,
	})

	response.JSON(w, http.StatusOK, sess)
}

// NewUserController creates a new user controller
func NewUserController(
	userService user.ServiceInterface,
	dataManager *data.Manager,
	utility *u.Utility,
) *UserController {
	return &UserController{
		userService: userService,
		dataManager: dataManager,
		utility:     utility,
	}
}
