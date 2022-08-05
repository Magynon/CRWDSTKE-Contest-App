package api

import (
	exam_api_domain "exam-api/domain"
	"fmt"
	"github.com/emicklei/go-restful/v3"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func (api *API) createProductSingle(req *restful.Request, resp *restful.Response) {
	product := exam_api_domain.Product{}
	err := req.ReadEntity(&product)
	if err != nil {
		log.Errorf("Failed to read product, err=%v", err)
		resp.WriteError(http.StatusConflict, fmt.Errorf("read error: %v", err))
		return
	}

	id, alreadyInDatabase, err := api.storage.Save(product)

	if err != nil {
		log.Errorf("Failed to save product, err=%v", err)
		resp.WriteError(http.StatusInternalServerError, fmt.Errorf("save error: %v", err))
		return
	}

	if alreadyInDatabase {
		log.Errorf("Failed to save product, err=%v", err)
		resp.WriteError(http.StatusFound, fmt.Errorf("already in database"))
		return
	}

	resp.WriteAsJson(id)

}

func (api *API) getProductSingle(req *restful.Request, resp *restful.Response) {
	id := req.QueryParameter("id")
	product := exam_api_domain.Product{}
	if id == "" {
		log.Errorf("Failed to read id")
		_ = resp.WriteError(http.StatusBadRequest, fmt.Errorf("id must be provided"))
		return
	}

	product, isProductThere, err := api.storage.Get(id)
	if err != nil {
		log.Errorf("Failed to get product, err=%v", err)
		_ = resp.WriteError(http.StatusNotFound, fmt.Errorf("failed to get product %v", err))
		return
	}

	if !isProductThere {
		log.Errorf("Failed to get product, err=%v", err)
		_ = resp.WriteError(http.StatusNotFound, fmt.Errorf("product is not available"))
		return
	}

	resp.WriteAsJson(product)
}

func (api *API) updateProductSingle(req *restful.Request, resp *restful.Response) {
	product := exam_api_domain.Product{}
	err := req.ReadEntity(&product)
	if err != nil {
		log.Errorf("Failed to read product, err=%v", err)
		resp.WriteError(http.StatusConflict, fmt.Errorf("read error: %v", err))
		return
	}

	id := product.GetHash()
	if id == "" {
		log.Errorf("Failed to read id")
		_ = resp.WriteError(http.StatusBadRequest, fmt.Errorf("id not ok"))
		return
	}

	alreadyInDatabase, err := api.storage.Update(id, product)

	if !alreadyInDatabase {
		log.Errorf("Product not found in database")
		_ = resp.WriteError(http.StatusNotFound, fmt.Errorf("product not found"))
		return
	}
	
	resp.WriteAsJson(id)
}

func (api *API) deleteProductSingle(req *restful.Request, resp *restful.Response) {
	panic("TODO")
}
