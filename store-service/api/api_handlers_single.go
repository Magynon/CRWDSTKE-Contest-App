package api

import (
	"exam-store/domain"
	exam_api_domain "exam-store/domain"
	"fmt"
	"github.com/emicklei/go-restful/v3"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func (api *API) createProductSingle(req *restful.Request, resp *restful.Response) {
	product := domain.Product{}
	err := req.ReadEntity(&product)
	if err != nil {
		log.Errorf("Failed to read product, err=%v", err)
		resp.WriteError(http.StatusConflict, fmt.Errorf("read error: %v", err))
		return
	}

	id, alreadyInDatabase, err := api.storage.Save(product)

	if err != nil {
		log.Errorf("Failed to save product, err=%v", err)
		resp.WriteError(http.StatusConflict, fmt.Errorf("save error: %v", err))
		return
	}

	if alreadyInDatabase {
		log.Errorf("Failed to save product, err=%v", err)
		resp.WriteError(http.StatusConflict, fmt.Errorf("already in database"))
		return
	}

	resp.WriteAsJson(id)
	log.Infof("Product %v created", id)

}

func (api *API) getProductSingle(req *restful.Request, resp *restful.Response) {
	id := req.QueryParameter("id")
	product := domain.Product{}
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
	log.Infof("Product %v got", id)
}

func (api *API) updateProductSingle(req *restful.Request, resp *restful.Response) {
	product := exam_api_domain.Product{}
	err := req.ReadEntity(&product)
	if err != nil {
		log.Errorf("Failed to read product, err=%v", err)
		resp.WriteError(http.StatusConflict, fmt.Errorf("read error: %v", err))
		return
	}

	id := product.Name
	if id == "" {
		log.Errorf("Failed to read id")
		_ = resp.WriteError(http.StatusBadRequest, fmt.Errorf("id not ok"))
		return
	}

	alreadyInDatabase, err := api.storage.Update(id, product)

	if err != nil {
		log.Errorf("Failed to insert in database: %v", err)
		resp.WriteError(http.StatusConflict, fmt.Errorf("read error: %v", err))
		return
	}

	if !alreadyInDatabase {
		log.Errorf("Product not found in database")
		_ = resp.WriteError(http.StatusNotFound, fmt.Errorf("product not found"))
		return
	}

	resp.WriteAsJson(id)
	log.Infof("Product %v updated", id)
}

func (api *API) deleteProductSingle(req *restful.Request, resp *restful.Response) {

	id := req.QueryParameter("id")
	if id == "" {
		log.Errorf("Failed to read id")
		_ = resp.WriteError(http.StatusBadRequest, fmt.Errorf("id must be provided"))
		return
	}

	productFound, err := api.storage.Delete(id)

	if err != nil {
		log.Errorf("Failed to find product in database: %v", err)
		resp.WriteError(http.StatusConflict, fmt.Errorf("read error: %v", err))
		return
	}

	if !productFound {
		log.Errorf("Failed to find product in database: %v", err)
		resp.WriteError(http.StatusConflict, fmt.Errorf("read error: %v", err))
		return
	}

	resp.WriteAsJson(id)
	log.Infof("Product %v deleted", id)

}
