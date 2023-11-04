package handlers

import "github.com/julienschmidt/httprouter"

type Handler interface {
	Route(router *httprouter.Router)
}
