package api

import (
	"exam-store/domain"
	"github.com/emicklei/go-restful/v3"
)

const (
	productPath = "/product"
)

type API struct {
	storage domain.Storage
}

func NewAPI(store domain.Storage) *API {
	return &API{
		storage: store,
	}
}

func (api *API) RegisterRoutes(ws *restful.WebService) {
	ws.Path("/store")
	ws.Route(ws.POST(productPath).To(api.createProductSingle))
	ws.Route(ws.GET(productPath).To(api.getProductSingle))
	ws.Route(ws.PATCH(productPath).To(api.updateProductSingle))
	ws.Route(ws.DELETE(productPath).To(api.deleteProductSingle))
}
