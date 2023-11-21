package main

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
)

func createJWT(account *Account) (string, error) {
	claims := &jwt.MapClaims{
		"expiresAt":     15_000,
		"accountNumber": account.Number,
	}

	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})

}

func permissionDenied(w http.ResponseWriter) {
	err := WriteJSON(w, http.StatusUnauthorized, ApiError{"permission denied"})
	if err != nil {
		return
	}
}

func withJWTAuth(handlerFunc http.HandlerFunc, s Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr := r.Header.Get("Authorization")
		token, err := validateJWT(tokenStr)

		if err != nil || !token.Valid {
			permissionDenied(w)
			return
		}

		userId, err := parseId(r)
		if err != nil {
			permissionDenied(w)
			return
		}

		account, err := s.GetAccountById(userId)
		if err != nil {
			permissionDenied(w)
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		if account.Number != int64(claims["accountNumber"].(float64)) {
			permissionDenied(w)
			return
		}

		handlerFunc(w, r)
	}
}
