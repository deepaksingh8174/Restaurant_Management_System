package middleware

import (
	"context"
	"example.com/database/dbHelper"
	"example.com/model"
	"net/http"
)

func SubAdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := r.Context().Value("claims").(*model.Claims)
		userID := claims.Userid
		isSubAdmin, subAdminErr := dbHelper.IsSubAdminRole(userID)

		if subAdminErr != nil {
			//utils.RespondJSON(w, http.StatusInternalServerError, utils.Status{Message: "Failed to check db"})
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if isSubAdmin {
			ctx := context.WithValue(r.Context(), "claims", claims)
			next.ServeHTTP(w, r.WithContext(ctx))
			return
		}
			w.WriteHeader(http.StatusUnauthorized)
			return

	})
}

func AdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		claims := r.Context().Value("claims").(*model.Claims)
		userID := claims.Userid
		isAdmin, adminErr := dbHelper.IsAdminRole(userID)
		if adminErr != nil {
			//utils.RespondJSON(w, http.StatusInternalServerError, utils.Status{Message: "Failed to check db"})
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if isAdmin {
			ctx := context.WithValue(r.Context(), "claims", claims)
			next.ServeHTTP(w, r.WithContext(ctx))
			return

		}
		w.WriteHeader(http.StatusUnauthorized)
		return

	})
}
