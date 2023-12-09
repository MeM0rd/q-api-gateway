package quote

import (
	"context"
	"encoding/json"
	"github.com/MeM0rd/q-api-gateway/internal/handlers"
	"github.com/MeM0rd/q-api-gateway/pkg/logger"
	quotePbService "github.com/MeM0rd/q-api-gateway/pkg/pb/quote"
	"github.com/MeM0rd/q-api-gateway/pkg/sessions"
	"github.com/MeM0rd/q-api-gateway/pkg/utils/response"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"strconv"
)

type handler struct {
	logger logger.Logger
}

func NewHandler(logger *logger.Logger) handlers.Handler {
	return &handler{
		logger: *logger,
	}
}

func (h *handler) Route(r *httprouter.Router) {
	r.GET("/quotes", h.GetList)
	r.POST("/quotes", h.Create)
	r.DELETE("/quotes", h.Delete)
}

func (h *handler) GetList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	conn := quotePbService.NewConn(&h.logger)
	quotePbClient := quotePbService.NewQuotePbServiceClient(conn)
	getListResponse, err := quotePbClient.GetList(context.Background(), &quotePbService.GetListRequest{})
	if err != nil {
		h.logger.Infof("error in quote service: %v", err)
		response.InternalServerError(w)
		return
	}

	response.Success(w, getListResponse.Quotes)
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	var quote Quote

	err := json.NewDecoder(r.Body).Decode(&quote)
	if err != nil {
		h.logger.Infof("error decoding request: %v", err)
		response.BadRequest(w)
		return
	}

	conn := quotePbService.NewConn(&h.logger)
	quotePbClient := quotePbService.NewQuotePbServiceClient(conn)
	createResponse, err := quotePbClient.Create(context.Background(), &quotePbService.CreateRequest{
		Title:  quote.Title,
		Text:   quote.Text,
		UserId: int64(quote.UserId),
	})
	if err != nil {
		h.logger.Infof("error in quote service: %v", err)
		response.InternalServerError(w)
		return
	}

	if createResponse == nil {
		response.InternalServerError(w)
		return
	}

	response.Success(w, createResponse.Quote)
}

func (h *handler) Delete(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	session, err := sessions.CheckSession(r)
	if err != nil {
		log.Printf("error checking sesion in delete quote func: %v", err)
		response.BadRequest(w)
		return
	}
	h.logger.Infof("quote id from request: %v", params.ByName("id"))
	quoteId, err := strconv.Atoi(params.ByName("id"))
	if err != nil {
		h.logger.Infof("error getting quote id: %v", err)
		response.BadRequest(w)
		return
	}

	conn := quotePbService.NewConn(&h.logger)
	quotePbClient := quotePbService.NewQuotePbServiceClient(conn)
	deleteResponse, err := quotePbClient.Delete(context.Background(), &quotePbService.DeleteRequest{
		QuoteId: int64(quoteId),
		UserId:  int64(session.UserId),
	})
	if err != nil {
		response.Common(w, 200, deleteResponse.Err)
		return
	}

	if deleteResponse.Status != true {
		response.InternalServerError(w)
		return
	}

	response.Success(w, []byte(deleteResponse.Msg))
}
