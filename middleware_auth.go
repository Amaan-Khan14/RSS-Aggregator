package main

import (
	"fmt"
	"net/http"

	"github.com/Amaan-Khan14/RSS-Aggregator/internal/auth"
	"github.com/Amaan-Khan14/RSS-Aggregator/internal/database"
)

type authHandler func(http.ResponseWriter, *http.Request, database.User)

func (apicfg *apiConfig) middlewareAuht(handler authHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetApiKey(r.Header)
		if err != nil {
			respondWithErr(w, 403, fmt.Sprint("Unauthorized:", err))
			return
		}
		user, err := apicfg.DB.GetUserByApiKey(r.Context(), apiKey)
		if err != nil {
			respondWithErr(w, 400, fmt.Sprint("Could'nt get user:", err))
			return
		}
		handler(w, r, user)
	}
}
