package middleware

import (
	"context"
	"example.com/model"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"os"
)

func JWTMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		token, err := jwt.ParseWithClaims(tokenString, &model.Claims{}, func(t *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SecretKey")), nil
		})
		if err != nil || !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "claims", token.Claims.(*model.Claims))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
