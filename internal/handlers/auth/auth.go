package auth

import (
	"context"
	"encoding/json"
	"github.com/MeM0rd/q-api-gateway/internal/handlers"
	"github.com/MeM0rd/q-api-gateway/pkg/logger"
	authPbService "github.com/MeM0rd/q-api-gateway/pkg/pb/auth"
	"github.com/MeM0rd/q-api-gateway/pkg/sessions"
	"github.com/MeM0rd/q-api-gateway/pkg/utils/response"
	"github.com/julienschmidt/httprouter"
	"google.golang.org/grpc"
	"net/http"
	"os"
	"time"
)

type handler struct {
	logger logger.Logger
}

func NewHandler(l *logger.Logger) handlers.Handler {
	Init(l)
	return &handler{
		logger: *l,
	}
}

func (h *handler) Route(r *httprouter.Router) {
	r.POST("/auth/register", h.Register)
	r.POST("/auth/login", h.Login)
	r.GET("/auth/logout", h.Logout)
}

var conn *grpc.ClientConn
var authPbClient authPbService.AuthPbServiceClient

func Init(l *logger.Logger) {
	var err error

	conn, err = authPbService.NewConnection()
	if err != nil {
		l.Infof("Cannot create authPbService conn: %v", err)
	}

	authPbClient = authPbService.NewAuthPbServiceClient(conn)
}

func (h *handler) Register(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var rm RegisterModel

	err := json.NewDecoder(r.Body).Decode(&rm)
	if err != nil {
		h.logger.Fatalf("Erorr decoding: %v", err)
		response.BadRequest(w)
		return
	}

	registerResponse, err := authPbClient.Register(context.Background(), &authPbService.RegisterRequest{
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

	loginResponse, err := authPbClient.Login(context.Background(), &authPbService.LoginRequest{
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

	h.logger.Infof("%v, %v, %v", loginResponse.Msg, loginResponse.Token, expiredTime)

	cookie := &http.Cookie{
		Name:    loginResponse.Msg,
		Value:   loginResponse.Token,
		Expires: expiredTime,
		Path:    "/",
	}

	http.SetCookie(w, cookie)

	response.Common(w, 200, "logged in")
}

func (h *handler) Logout(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	token, err := sessions.GetToken(r)
	if err != nil {
		h.logger.Infof("error getting token during logout: %v", err)
		response.BadRequest(w)
		return
	}

	loginResponse, err := authPbClient.Logout(context.Background(), &authPbService.LogoutRequest{
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
		Name:    os.Getenv("COOKIE_NAME"),
		Value:   "",
		Expires: time.Now(),
		Path:    "/",
	}

	http.SetCookie(w, cookie)

	response.Common(w, 200, "logged out")
}
