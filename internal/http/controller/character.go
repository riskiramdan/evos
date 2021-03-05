package controller

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/riskiramdan/evos/internal/character"
	"github.com/riskiramdan/evos/internal/data"
	"github.com/riskiramdan/evos/internal/http/response"
	"github.com/riskiramdan/evos/internal/types"
	u "github.com/riskiramdan/evos/util"
)

// CharacterController represents the character controller
type CharacterController struct {
	characterService character.ServiceInterface
	dataManager      *data.Manager
	utility          *u.Utility
}

// CharacterList character list and count
type CharacterList struct {
	Data  []*character.Characters `json:"data"`
	Total int                     `json:"total"`
}

// GetListCharacter function for get list data characters
func (a *CharacterController) GetListCharacter(w http.ResponseWriter, r *http.Request) {
	var err *types.Error

	queryValues := r.URL.Query()
	var limit = 10
	var errConversion error
	if queryValues.Get("limit") != "" {
		limit, errConversion = strconv.Atoi(queryValues.Get("limit"))
		if errConversion != nil {
			err = &types.Error{
				Path:    ".CharacterController->ListCharacter()",
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
				Path:    ".CharacterController->ListCharacter()",
				Message: errConversion.Error(),
				Error:   errConversion,
				Type:    "golang-error",
			}
			response.Error(w, "Bad Request", http.StatusBadRequest, *err)
			return
		}
	}

	var name = queryValues.Get("name")

	if limit < 0 {
		limit = 10
	}
	if page < 0 {
		page = 1
	}
	characterList, count, err := a.characterService.ListCharacters(r.Context(), &character.FindAllCharacterParams{
		Name:  name,
		Limit: limit,
		Page:  page,
	})
	if err != nil {
		err.Path = ".CharacterController->ListCharacter()" + err.Path
		if err.Error != data.ErrNotFound {
			response.Error(w, "Internal Server Error", http.StatusInternalServerError, *err)
			return
		}
	}
	if characterList == nil {
		characterList = []*character.Characters{}
	}

	response.JSON(w, http.StatusOK, CharacterList{
		Data:  characterList,
		Total: count,
	})
}

// PostCreateCharacter for creating data character
func (a *CharacterController) PostCreateCharacter(w http.ResponseWriter, r *http.Request) {
	var err *types.Error

	decoder := json.NewDecoder(r.Body)

	var params *character.TransactionParams
	errDecode := decoder.Decode(&params)
	if errDecode != nil {
		err = &types.Error{
			Path:    ".CharacterController->CreateCharacter()",
			Message: errDecode.Error(),
			Error:   errDecode,
			Type:    "golang-error",
		}
		response.Error(w, "Bad Request", http.StatusBadRequest, *err)
		return
	}
	errTransaction := a.dataManager.RunInTransaction(r.Context(), func(ctx context.Context) error {
		_, err = a.characterService.CreateCharacter(ctx, &character.TransactionParams{
			CharacterTypeID: params.CharacterTypeID,
			Name:            params.Name,
			Power:           params.Power,
		})
		if err != nil {
			return err.Error
		}
		return nil
	})
	if errTransaction != nil {
		err.Path = ".CharacterController->CreateCharacter()" + err.Path
		if errTransaction == character.ErrCharacterExists {
			response.Error(w, character.ErrCharacterExists.Error(), http.StatusUnprocessableEntity, *err)
			return
		}
		response.Error(w, "Internal Server Error", http.StatusInternalServerError, *err)
		return
	}

	response.JSON(w, http.StatusOK, "Character Created Successful")
}

// PutUpdateCharacter for update data character
func (a *CharacterController) PutUpdateCharacter(w http.ResponseWriter, r *http.Request) {
	var err *types.Error

	decoder := json.NewDecoder(r.Body)

	var params *character.TransactionParams
	errDecode := decoder.Decode(&params)
	if errDecode != nil {
		err = &types.Error{
			Path:    ".CharacterController->UpdateCharacter()",
			Message: errDecode.Error(),
			Error:   errDecode,
			Type:    "golang-error",
		}
		response.Error(w, "Bad Request", http.StatusBadRequest, *err)
		return
	}
	var sCharacterID = chi.URLParam(r, "characterId")
	characterID, errConversion := strconv.Atoi(sCharacterID)
	if errConversion != nil {
		err = &types.Error{
			Path:    ".CharacterController->UpdateCharacter()",
			Message: errConversion.Error(),
			Error:   errConversion,
			Type:    "golang-error",
		}
		response.Error(w, "Bad Request", http.StatusBadRequest, *err)
		return
	}

	errTransaction := a.dataManager.RunInTransaction(r.Context(), func(ctx context.Context) error {
		_, err = a.characterService.UpdateCharacter(ctx, characterID, params)
		if err != nil {
			return err.Error
		}
		return nil
	})
	if errTransaction != nil {
		err.Path = ".CharacterController->UpdateCharacter()" + err.Path
		if errTransaction == data.ErrAlreadyExist {
			response.Error(w, data.ErrAlreadyExist.Error(), http.StatusUnprocessableEntity, *err)
			return
		}
		response.Error(w, "Internal Server Error", http.StatusInternalServerError, *err)
		return
	}
	response.JSON(w, http.StatusOK, "Update Character Successful")

}

// NewCharacterController creates a new character controller
func NewCharacterController(
	characterService character.ServiceInterface,
	dataManager *data.Manager,
) *CharacterController {
	return &CharacterController{
		characterService: characterService,
		dataManager:      dataManager,
	}
}
