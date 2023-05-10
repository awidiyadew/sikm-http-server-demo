package middleware

import (
	"context"
	"demo-app/model"
	"encoding/json"
	"net/http"
)

func ValidateMethod(method string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			http.Error(w, "method is not allowed!", http.StatusMethodNotAllowed)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func Auth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c, err := r.Cookie("username")
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(model.ErrorResp{
				Message: "unauthorized error",
			})
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "username", c.Value)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
