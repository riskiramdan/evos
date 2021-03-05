package http

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/riskiramdan/evos/internal/appcontext"
	"github.com/riskiramdan/evos/internal/constants"
	"github.com/riskiramdan/evos/internal/data"
	"github.com/riskiramdan/evos/internal/http/response"
	"github.com/riskiramdan/evos/internal/types"
	"github.com/riskiramdan/evos/internal/user"

	"github.com/dgrijalva/jwt-go"
)

func (hs *Server) authorizedOnly(userService user.ServiceInterface) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			var tokenString string

			ctx := r.Context()
			tokenString = getBearerToken(r)
			if tokenString == "" {
				response.Error(w, "Unauthorized", http.StatusUnauthorized, types.Error{
					Path:    ".Server->authorizeOnly()",
					Message: "",
					Error:   nil,
					Type:    "",
				})
				return
			}

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Signing method invalid")
				} else if method != jwt.SigningMethodHS256 {
					return nil, fmt.Errorf("Signing method invalid")
				}

				return constants.SignatureKey, nil
			})
			if err != nil {
				response.Error(w, "Unauthorized", http.StatusUnauthorized, types.Error{
					Path:    ".Server->authorizeOnly()",
					Message: err.Error(),
					Error:   err,
					Type:    "Invalid Token",
				})
				return
			}
			singleUser, errT := userService.GetByToken(ctx, tokenString)
			if errT != nil {
				if errT.Error != data.ErrNotFound {
					response.Error(w, "Internal Server Error", http.StatusInternalServerError, *errT)
					return
				}
				response.Error(w, "Unauthorized", http.StatusUnauthorized, types.Error{
					Path:    ".Server->authorizeOnly()",
					Message: "",
					Error:   nil,
					Type:    "",
				})
				return
			}
			_, ok := token.Claims.(jwt.MapClaims)
			if !ok || !token.Valid {
				response.Error(w, "Unauthorized", http.StatusUnauthorized, types.Error{
					Path:    ".Server->authorizeOnly()",
					Message: "",
					Error:   nil,
					Type:    "",
				})
				return
			}
			ctx = context.WithValue(ctx, appcontext.KeyUserID, singleUser.ID)
			ctx = context.WithValue(ctx, appcontext.KeySessionID, *singleUser.Token)
			if singleUser.RoleID == 1 {

				ctx = context.WithValue(ctx, appcontext.KeyIsAdmin, true)
			}

			next.ServeHTTP(w, r.WithContext(ctx))
		}

		return http.HandlerFunc(fn)
	}
}

func getBearerToken(r *http.Request) string {
	token := r.Header.Get("Authorization")
	splitToken := strings.Split(token, "Bearer")

	if len(splitToken) < 2 {
		return ""
	}

	token = strings.Trim(splitToken[1], " ")
	return token
}

// AdminOnly ..
func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if !appcontext.IsAdmin(ctx) {
			response.Error(w, "Admin Only", http.StatusUnauthorized, types.Error{
				Path:    ".Server->authorizeOnly()",
				Message: "",
				Error:   nil,
				Type:    "",
			})
			return
		}
		next.ServeHTTP(w, r)
	})
}
