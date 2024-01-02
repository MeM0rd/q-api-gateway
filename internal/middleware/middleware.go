package middleware

import (
	"github.com/MeM0rd/q-api-gateway/pkg/utils/response"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
)

func Auth(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		_, err := r.Cookie(os.Getenv("COOKIE_NAME"))
		if err != nil {
			log.Printf("unauthorized in mw.Auth: %v", err)
			response.Unauthorized(w)
			return
		}

		next(w, r, params)
	}
}
