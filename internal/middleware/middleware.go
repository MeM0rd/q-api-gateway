package middleware

import (
	"github.com/MeM0rd/q-api-gateway/pkg/sessions"
	"github.com/MeM0rd/q-api-gateway/pkg/utils/response"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func Auth(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		_, err := sessions.CheckSession(r)
		if err != nil {
			log.Printf("unauthorized in mw.Auth: %v", err)
			response.Unauthorized(w)
			return
		}

		next(w, r, params)
	}
}
