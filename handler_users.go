package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Amaan-Khan14/RSS-Aggregator/internal/database"
	"github.com/google/uuid"
)

// To Create user in the database
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

	respondWithJSON(w, 201, databaseUserToUser(user))
}

// To GetUsers Based on API Key
func (apiCfg *apiConfig) handlerGetUser(w http.ResponseWriter, r *http.Request, user database.User) {
	respondWithJSON(w, 200, databaseUserToUser(user))
}

// To GetUsers Based on name
func (apiCfg *apiConfig) handlerGetUserByName(w http.ResponseWriter, r *http.Request) {
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
	user, err := apiCfg.DB.GetUserByName(r.Context(), params.Name)
	if err != nil {
		respondWithErr(w, 400, fmt.Sprint("Could'nt get user:", err))
		return
	}

	respondWithJSON(w, 200, user)
}

func (apiCfg *apiConfig) handlerGetPostsForUser(w http.ResponseWriter, r *http.Request, user database.User) {
	posts, err := apiCfg.DB.GetPost(r.Context(), database.GetPostParams{
		UserID: user.ID,
		Limit:  10,
	})
	if err != nil {
		respondWithErr(w, 400, fmt.Sprint("Could'nt get posts:", err))
		return
	}

	respondWithJSON(w, 200, databasePostsToPosts(posts))
}
