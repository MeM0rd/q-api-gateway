package response

import (
	"net/http"
)

func Success(w http.ResponseWriter, json []byte) {
	w.Header().Set("Content-Type", "Application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(json)
}

func InternalServerError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "Application/json")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("something went wrong"))
}

func BadRequest(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "Application/json")
	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("bad request"))
}

func NotFound(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "Application/json")
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("not found"))
}

func Created(w http.ResponseWriter, json []byte) {
	w.Header().Set("Content-Type", "Application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(json)
}

func NoContent(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "Application/json")
	w.WriteHeader(http.StatusNoContent)
}

func Common(w http.ResponseWriter, code int, message string) {
	w.Header().Set("Content-Type", "Application/json")
	w.WriteHeader(code)
	w.Write([]byte(message))
}

func Unauthorized(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	w.Write([]byte("unauthorized"))
}
