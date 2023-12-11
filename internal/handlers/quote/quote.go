package quote

import (
	"context"
	"encoding/json"
	"github.com/MeM0rd/q-api-gateway/internal/handlers"
	mw "github.com/MeM0rd/q-api-gateway/internal/middleware"
	"github.com/MeM0rd/q-api-gateway/pkg/logger"

	quotePbService "github.com/MeM0rd/q-api-gateway/pkg/pb/quote"
	"github.com/MeM0rd/q-api-gateway/pkg/sessions"
	"github.com/MeM0rd/q-api-gateway/pkg/utils/response"
	"github.com/julienschmidt/httprouter"
	"google.golang.org/grpc"
	"log"
	"net/http"
	"strconv"
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
	r.GET("/quotes", h.GetList)
	r.POST("/quotes", mw.Auth(h.Create))
	r.DELETE("/quotes/:id", mw.Auth(h.Delete))
}

var conn *grpc.ClientConn
var quotePbClient quotePbService.QuotePbServiceClient

func Init(l *logger.Logger) {
	var err error

	conn, err = quotePbService.NewConnection()
	if err != nil {
		l.Infof("Cannot create quotePbService conn: %v", err)
	}

	quotePbClient = quotePbService.NewQuotePbServiceClient(conn)
}

func (h *handler) GetList(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
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
