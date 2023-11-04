package profile

import (
	"github.com/MeM0rd/q-api-gateway/internal/handlers"
	"github.com/julienschmidt/httprouter"
)

type handler struct {
}

func NewHandler() handlers.Handler {
	return &handler{}
}

func (h *handler) Route(r *httprouter.Router) {

}
