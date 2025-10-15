package middleware

import (
	"context"
	"example/todolist/auth"
	"example/todolist/response"
	"log"
	"net/http"
)

type key int

const UserIdKey key = 0

func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")

		if token == "" {
			//if no token return unauthorized
			log.Printf("AuthMiddleware - missing bearer token")
			response.WriteError(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		//check token validity
		token = token[len("Bearer "):]
		mapClaims, err := auth.VerifyAccessToken(token)
		if err != nil {
			log.Printf("AuthMiddleware - error verifying token: %v", err)
			response.WriteError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		userId, ok := mapClaims["user_id"].(float64)
		if !ok {
			log.Printf("AuthMiddleware - user_id missing from claims")
			response.WriteError(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		ctx := context.WithValue(r.Context(), UserIdKey, int64(userId))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func WbsocketTokenAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.FormValue("token")

		if token == "" {
			//if no token return unauthorized
			log.Printf("Websocket AuthMiddleware - missing bearer token")
			response.WriteError(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		//check token validity
		// token = token[len("Bearer "):]
		mapClaims, err := auth.VerifyAccessToken(token)
		if err != nil {
			log.Printf("AuthMiddleware - error verifying token: %v", err)
			response.WriteError(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		userId, ok := mapClaims["user_id"].(float64)
		if !ok {
			log.Printf("AuthMiddleware - user_id missing from claims")
			response.WriteError(w, http.StatusUnauthorized, "unauthorized")
			return
		}
		ctx := context.WithValue(r.Context(), UserIdKey, int64(userId))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
