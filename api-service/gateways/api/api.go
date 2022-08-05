package api

import (
	"exam-api/domain"
	"exam-api/gateways/remote"

	"github.com/emicklei/go-restful/v3"
)

const (
	memoryRootPath = "/memory"
	httpRootPath   = "/http"
	redisRootPath  = "/redis"

	productPath = "/product"

	versionSingle = "/single"
	versionBatch  = "/batch"
)

type API struct {
	storage domain.Storage
	client  remote.Client
}

func NewAPI(store domain.Storage) *API {
	return &API{
		storage: store,
	}
}

func (api *API) RegisterRoutes(ws *restful.WebService) {
	ws.Path("/store")
	ws.Route(ws.POST(memoryRootPath + productPath + versionSingle).To(api.createProductMemorySingle))
	ws.Route(ws.GET(memoryRootPath + productPath + versionSingle).To(api.getProductMemorySingle))
	ws.Route(ws.PATCH(memoryRootPath + productPath + versionSingle).To(api.updateProductMemorySingle))
	ws.Route(ws.DELETE(memoryRootPath + productPath + versionSingle).To(api.deleteProductMemorySingle))

	ws.Route(ws.POST(memoryRootPath + productPath + versionBatch).To(api.createProductMemoryBatch))
	ws.Route(ws.GET(memoryRootPath + productPath + versionBatch).To(api.getProductMemoryBatch))
	ws.Route(ws.PATCH(memoryRootPath + productPath + versionBatch).To(api.updateProductMemoryBatch))
	ws.Route(ws.DELETE(memoryRootPath + productPath + versionBatch).To(api.deleteProductMemoryBatch))

	ws.Route(ws.POST(httpRootPath + productPath + versionSingle).To(api.createProductHTTPSingle))
	ws.Route(ws.GET(httpRootPath + productPath + versionSingle).To(api.getProductHTTPSingle))
	ws.Route(ws.PATCH(httpRootPath + productPath + versionSingle).To(api.updateProductHTTPSingle))
	ws.Route(ws.DELETE(httpRootPath + productPath + versionSingle).To(api.deleteProductHTTPSingle))

	ws.Route(ws.POST(httpRootPath + productPath + versionBatch).To(api.createProductHTTPBatch))
	ws.Route(ws.GET(httpRootPath + productPath + versionBatch).To(api.getProductHTTPBatch))
	ws.Route(ws.PATCH(httpRootPath + productPath + versionBatch).To(api.updateProductHTTPBatch))
	ws.Route(ws.DELETE(httpRootPath + productPath + versionBatch).To(api.deleteProductHTTPBatch))

}
