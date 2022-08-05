package gateways

import (
	"github.com/emicklei/go-restful/v3"
)

func (api *API) RegisterRoutes(ws *restful.WebService) {
	ws.Path("/database")
	//ws.Route(ws.POST(productPath).To(api.productPOSTHandler).Reads(restful.MIME_JSON).Writes(restful.MIME_JSON).Doc("Adds products"))
	//ws.Route(ws.GET(productPath).To(api.productGETHandler).Writes(restful.MIME_JSON).Doc("Gets products"))
}
