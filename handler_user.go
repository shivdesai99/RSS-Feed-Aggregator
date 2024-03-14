package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/RSS-Feed-Aggregator/internal/auth"
	"github.com/RSS-Feed-Aggregator/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithErr(w, 400, fmt.Sprintf("Error parsing JSON: %s", err))
	}

	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID: uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name: 	   params.Name,
	})
	if err != nil {
		respondWithErr(w, 400, fmt.Sprintf("Error creating user: %s", err))
		return
	}

	respondWithJSON(w, 201, databaseUsertoUser(user))
}

func (apiCfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		respondWithErr(w, 403, fmt.Sprintf("Auth error: %v", err)) // 403 -> Permissions Error
	}

	user, err := apiCfg.DB.GetUserByAPIKey(r.Context(), apiKey)
	if err != nil {
		respondWithErr(w, 400, fmt.Sprintf("Error getting user: %v", err))
		return
	}

	respondWithJSON(w, 200, databaseUsertoUser(user))
}