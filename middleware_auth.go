package main

import (
	"fmt"
	"net/http"

	"github.com/FilipLusnia/rssagg/auth"
	"github.com/FilipLusnia/rssagg/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {

	return func(handlerW http.ResponseWriter, handlerR *http.Request) {

		apiKey, err := auth.GetAPIKey(handlerR.Header)
		if err != nil {
			respondWithError(handlerW, 403, fmt.Sprintf("Auth error: %s", err))
			return
		}

		user, err := apiCfg.DB.GetUserByAPIKey(handlerR.Context(), apiKey)
		if err != nil {
			respondWithError(handlerW, 400, fmt.Sprintf("Couldn't get user: %s", err))
			return
		}

		handler(handlerW, handlerR, user)
	}
}
