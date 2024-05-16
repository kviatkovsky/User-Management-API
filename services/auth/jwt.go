package auth

import (
	"context"
	"fmt"
	"main/configs"
	"main/types"
	"main/utils"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserKey contextKey = "userID"
const AccessLevelAdmin = 3

func CreateJWT(secret []byte, userId uint) (string, error) {
	expiration := time.Second * time.Duration(configs.Envs.JWTExpirationInSeconds)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID": userId,
		"exp":    time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func AdminWithJWTAuth(handlerFunc http.HandlerFunc, store types.UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		tokenString := utils.GetTokenFromRequest(r)
		token, err := validateJWT(tokenString)
		if err != nil {
			permissionDenied(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		userID := claims["userID"].(float64)

		if err != nil {
			utils.WriteError(w, http.StatusForbidden, fmt.Errorf("failed to convert userID to int: %v", err))
			permissionDenied(w)
			return
		}
		u, err := store.GetUserById(int(userID))
		if err != nil {
			utils.WriteError(w, http.StatusForbidden, fmt.Errorf("failed to get user by id: %v", err))
			permissionDenied(w)
			return
		}

		superAdmin, err := isUserAdmin(u.ID, store)
		if err != nil {
			utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
		}

		if !superAdmin {
			utils.WriteError(w, http.StatusForbidden, fmt.Errorf("User doesn't have correct access level"))
			permissionDenied(w)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, u.ID)
		r = r.WithContext(ctx)

		handlerFunc(w, r)
	}
}

func isUserAdmin(uid uint, store types.UserStore) (bool, error) {
	userRole, err := store.GetUserRoleByUserId(uid)
	if err != nil {
		return false, err
	}
	fmt.Println(userRole)

	return userRole.AccessLevel == AccessLevelAdmin, nil
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(configs.Envs.JWTSecret), nil
	})
}

func permissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
}
