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
func (apiCfg *apiConfig) handlerCreateFeed(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		Name string `json:"name"`
		Url  string `json:"url"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithErr(w, 400, fmt.Sprintf("Error parsing JSON:", err))
		return
	}
	feed, err := apiCfg.DB.CreateFeed(r.Context(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      params.Name,
		Url:       params.Url,
		UserID:    user.ID,
	})
	if err != nil {
		respondWithErr(w, 400, fmt.Sprint("Error creating user:", err))
		return
	}

	respondWithJSON(w, 201, databaseFeedToFeed(feed))
}

// To GetFeed
func (apiCfg *apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {
	feeds, err := apiCfg.DB.GetFeed(r.Context())
	if err != nil {
		respondWithErr(w, 400, fmt.Sprint("Could'nt get feeds:", err))
		return
	}
	respondWithJSON(w, 200, databaseFeedsToFeeds(feeds))
}
