package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Amaan-Khan14/RSS-Aggregator/internal/database"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

// To Create user in the database
func (apiCfg *apiConfig) handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	type parameters struct {
		FeedID uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithErr(w, 400, fmt.Sprintf("Error parsing JSON:", err))
		return
	}
	feedfollow, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedID,
	})
	if err != nil {
		respondWithErr(w, 400, fmt.Sprint("Error creating Feed Follow:", err))
		return
	}

	respondWithJSON(w, 201, databaseFeedFollowToFeedFollow(feedfollow))
}

// To GetFeedFollows
func (apicfg *apiConfig) getFeedFollows(w http.ResponseWriter, r *http.Request, user database.User) {
	feedfollow, err := apicfg.DB.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		respondWithErr(w, 400, fmt.Sprint("Could'nt get feed follows:", err))
		return
	}
	respondWithJSON(w, 200, databaseFeedFollowsToFeedFollows(feedfollow))
}

// To GetFeedFollows
func (apicfg *apiConfig) deleteFeedFollow(w http.ResponseWriter, r *http.Request, user database.User) {
	feedFollowIdStr := chi.URLParam(r, "feedFollowId")
	feedFollowId, err := uuid.Parse(feedFollowIdStr)
	if err != nil {
		respondWithErr(w, 400, fmt.Sprint("Could'nt parse feed follow id:", err))
		return
	}
	err = apicfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedFollowId,
		UserID: user.ID,
	})
	if err != nil {
		respondWithErr(w, 400, fmt.Sprint("Could'nt delete feed follow:", err))
	}
}
