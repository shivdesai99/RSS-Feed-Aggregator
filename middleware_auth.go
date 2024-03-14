package main

import (
	"fmt"
	"net/http"

	"github.com/RSS-Feed-Aggregator/internal/auth"
	"github.com/RSS-Feed-Aggregator/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithErr(w, 403, fmt.Sprintf("Auth error: %v", err)) // 403 -> Permissions Error
		}
	
		user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			respondWithErr(w, 400, fmt.Sprintf("Error getting user: %v", err))
			return
		}

		handler(w, r , user)
	}
}