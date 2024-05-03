package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Amaan-Khan14/RSS-Aggregator/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig) handlerCreateUsers(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithErr(w, 400, fmt.Sprintf("Error parsing JSON:", err))
		return
	}
	user, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
	})
	if err != nil {
		respondWithErr(w, 400, fmt.Sprint("Error creating user:", err))
		return
	}

	respondWithJSON(w, 200, databaseUserToUser(user))
}
