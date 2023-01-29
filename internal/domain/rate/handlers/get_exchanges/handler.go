package get_exchanges

import (
	service2 "Currency/internal/domain/rate/service"
	"Currency/internal/infrastructure/dto"
	service3 "Currency/internal/infrastructure/service"
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"net/http"
)

const ExchangeRateHandlerTag = "GetExchangeHandler"

type GetExchangeRateHandler struct {
	service        *service2.ExchangeRateService
	queryExtractor *service3.QueryExtractor
}

func NewHandler(srv *service2.ExchangeRateService) *GetExchangeRateHandler {
	return &GetExchangeRateHandler{
		service:        srv,
		queryExtractor: service3.NewQueryExtractor(),
	}
}

func (h *GetExchangeRateHandler) ExchangeRate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	reqParam, err := h.queryExtractor.ExtractRateParams(r.URL.Query())
	exchangeRate := h.service.GetExchangeRate(reqParam.From, reqParam.To, reqParam.OnDate.OnDate)

	if err != nil {
		createResponse(dto.NewError(err), w)
		return
	}

	createResponse(exchangeRate, w)
	return
}

func (h *GetExchangeRateHandler) Convert(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	reqParam, err := h.queryExtractor.ExtractConvertParams(r.URL.Query())

	if err != nil {
		createResponse(dto.NewError(err), w)
		return
	}

	exchangeRate := h.service.ConvertCurrency(reqParam.From, reqParam.To, reqParam.OnDate.OnDate, reqParam.Value)
	createResponse(exchangeRate, w)
	return
}

func createResponse(v any, w http.ResponseWriter) {
	marshalled, err := json.Marshal(v)

	if err != nil {
		fallbackResponse, _ := json.Marshal(`{"errorMessage": "inconvertible response"}`)
		_, _ = w.Write(fallbackResponse)
	}

	_, _ = w.Write(marshalled)
}
