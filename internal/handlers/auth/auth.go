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
	"time"
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
	r.POST("/auth/login", h.Login)
	r.POST("/auth/logout", h.Logout)
}

func (h *handler) Register(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var rm RegisterModel

	err := json.NewDecoder(r.Body).Decode(&rm)
	if err != nil {
		h.logger.Fatalf("Erorr decoding: %v", err)
		response.BadRequest(w)
		return
	}

	conn := authPbService.NewConn(&h.logger)
	authClient := authPbService.NewAuthPbServiceClient(conn)
	registerResponse, err := authClient.Register(context.Background(), &authPbService.RegisterRequest{
		Email:    rm.Email,
		Surname:  rm.Surname,
		Name:     rm.Name,
		Password: rm.Password,
	})
	if err != nil {
		h.logger.Infof("error sendind to auth-service: %v", err)
		response.InternalServerError(w)
		return
	}
	if registerResponse.Err != "" {
		h.logger.Infof("error from auth-service: %v", registerResponse.Err)
		response.InternalServerError(w)
		return
	}

	h.logger.Infof("msg from auth svc: %v", registerResponse.Msg)

	response.Created(w, registerResponse.Msg)
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var lm LoginModel

	err := json.NewDecoder(r.Body).Decode(&lm)
	if err != nil {
		h.logger.Fatalf("Erorr decoding: %v", err)
		response.BadRequest(w)
		return
	}

	conn := authPbService.NewConn(&h.logger)
	authClient := authPbService.NewAuthPbServiceClient(conn)
	loginResponse, err := authClient.Login(context.Background(), &authPbService.LoginRequest{
		Email:    lm.Email,
		Password: lm.Password,
	})
	if err != nil {
		h.logger.Infof("error sendind to auth-service: %v", err)
		response.InternalServerError(w)
		return
	}
	if loginResponse.Err != "" {
		h.logger.Infof("error from auth-service: %v", loginResponse.Err)
		response.InternalServerError(w)
		return
	}

	expiredTime, _ := time.Parse(time.DateTime, loginResponse.Lifetime)

	cookie := &http.Cookie{
		Name:    loginResponse.Msg,
		Value:   loginResponse.Token,
		Expires: expiredTime,
	}

	http.SetCookie(w, cookie)

	response.Common(w, 200, "logged in")
}

func (h *handler) Logout(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var token string

	err := json.NewDecoder(r.Body).Decode(&token)
	if err != nil {
		h.logger.Fatalf("Erorr decoding: %v", err)
		response.BadRequest(w)
		return
	}

	conn := authPbService.NewConn(&h.logger)
	authClient := authPbService.NewAuthPbServiceClient(conn)
	loginResponse, err := authClient.Logout(context.Background(), &authPbService.LogoutRequest{
		Token: token,
	})
	if err != nil {
		h.logger.Infof("error sendind to auth-service: %v", err)
		response.InternalServerError(w)
		return
	}
	if loginResponse.Err != "" {
		h.logger.Infof("error from auth-service: %v", loginResponse.Err)
		response.InternalServerError(w)
		return
	}

	cookie := &http.Cookie{
		Name:    "",
		Value:   "",
		Expires: time.Now(),
	}

	http.SetCookie(w, cookie)

	response.Common(w, 200, "logged out")
}
