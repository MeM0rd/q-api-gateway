package auth

import (
	"context"
	"encoding/json"
	"github.com/MeM0rd/q-api-gateway/internal/handlers"
	"github.com/MeM0rd/q-api-gateway/pkg/logger"
	authPbService "github.com/MeM0rd/q-api-gateway/pkg/pb/auth"
	"github.com/MeM0rd/q-api-gateway/pkg/utils/response"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

type handler struct {
	logger logger.Logger
}

func NewHandler(l *logger.Logger) handlers.Handler {
	return &handler{
		logger: *l,
	}
}

func (h *handler) Route(r *httprouter.Router) {
	r.POST("/auth/register", h.Register)
}

func (h *handler) Register(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var lm LoginModel

	err := json.NewDecoder(r.Body).Decode(&lm)
	if err != nil {
		h.logger.Fatalf("Erorr decoding: %v", err)
		response.BadRequest(w)
		return
	}

	conn := authPbService.NewConn(&h.logger)
	authClient := authPbService.NewAuthPbServiceClient(conn)
	registerResponse, err := authClient.Register(context.Background(), &authPbService.RegisterRequest{
		Email:    lm.Email,
		Surname:  lm.Surname,
		Name:     lm.Name,
		Password: lm.Password,
	})
	if err != nil {
		h.logger.Infof("error sendind to auth-service: %v", err)
		response.InternalServerError(w)
		return
	}
	if registerResponse.Err != "" {
		h.logger.Infof("error from auth-service: %v", err)
		response.InternalServerError(w)
		return
	}

	response.Created(w, []byte(registerResponse.Status))
}
