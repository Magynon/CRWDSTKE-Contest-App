package api

import (
	"exam-api/domain"
	"fmt"
	"net/http"

	"github.com/emicklei/go-restful/v3"
	log "github.com/sirupsen/logrus"
)

func (api *API) createProductMemorySingle(req *restful.Request, resp *restful.Response) {
	product := &domain.Product{}
	err := req.ReadEntity(product)

	if err != nil {
		log.Errorf("Failed to read product, err=%v", err)
		_ = resp.WriteError(http.StatusBadRequest, err)
		return
	}
	id, alreadyExists, err := api.storage.Save(*product)
	if err != nil {
		log.Errorf("Failed to save product in storage, err=%v", err)
		_ = resp.WriteError(http.StatusInternalServerError, fmt.Errorf("failed to save product"))
		return
	}

	if alreadyExists {
		log.Infof("Product %s already in store", id)
		_ = resp.WriteError(http.StatusConflict, fmt.Errorf("product already exists"))
		return
	}

	log.Infof("Product %s saved in store", id)

	_ = resp.WriteAsJson(map[string]string{
		"id": id,
	})

}

func (api *API) createProductHTTPSingle(req *restful.Request, resp *restful.Response) {
	product := &domain.Product{}
	err := req.ReadEntity(product)

	if err != nil {
		log.Errorf("Failed to read product, err=%v", err)
		_ = resp.WriteError(http.StatusBadRequest, err)
		return
	}
	id, alreadyExists, err := api.client.Save(*product)
	if err != nil {
		log.Errorf("Failed to save product in storage, err=%v", err)
		_ = resp.WriteError(http.StatusInternalServerError, fmt.Errorf("failed to save product"))
		return
	}

	if alreadyExists {
		log.Infof("Product %s already in store", id)
		_ = resp.WriteError(http.StatusConflict, fmt.Errorf("product already exists"))
		return
	}

	log.Infof("Product %s saved in store", id)

	_ = resp.WriteAsJson(map[string]string{
		"id": id,
	})

}

func (api *API) getProductMemorySingle(req *restful.Request, resp *restful.Response) {
	id := req.QueryParameter("id")
	if id == "" {
		log.Infof("No id provided in request")
		_ = resp.WriteError(http.StatusBadRequest, fmt.Errorf("id must be provided"))
		return
	}
	product, exists, err := api.storage.Get(id)
	if err != nil {
		log.Errorf("Failed to get product from storage, err=%v", err)
		_ = resp.WriteError(http.StatusInternalServerError, fmt.Errorf("failed to get product from store"))
		return
	}
	if !exists {
		log.Infof("Product %s not in store", id)
		_ = resp.WriteError(http.StatusNotFound, fmt.Errorf("product not found"))
		return
	}
	_ = resp.WriteAsJson(product)
}

func (api *API) getProductHTTPSingle(req *restful.Request, resp *restful.Response) {
	id := req.QueryParameter("id")
	if id == "" {
		log.Infof("No id provided in request")
		_ = resp.WriteError(http.StatusBadRequest, fmt.Errorf("id must be provided"))
		return
	}
	product, exists, err := api.client.Get(id)
	if err != nil {
		log.Errorf("Failed to get product from storage, err=%v", err)
		_ = resp.WriteError(http.StatusInternalServerError, fmt.Errorf("failed to get product from store"))
		return
	}
	if !exists {
		log.Infof("Product %s not in store", id)
		_ = resp.WriteError(http.StatusNotFound, fmt.Errorf("product not found"))
		return
	}
	_ = resp.WriteAsJson(product)
}

func (api *API) updateProductMemorySingle(req *restful.Request, resp *restful.Response) {
	productDiff := &domain.ProductDiff{}
	err := req.ReadEntity(productDiff)

	if err != nil {
		log.Errorf("Failed to read product, err=%v", err)
		_ = resp.WriteError(http.StatusBadRequest, err)
		return
	}

	// check if id exists in storage
	_, exists, err := api.storage.Get(productDiff.ID)
	if err != nil {
		log.Errorf("Failed to get product from storage, err=%v", err)
		_ = resp.WriteError(http.StatusInternalServerError, fmt.Errorf("failed to get product from store"))
		return
	}

	if !exists {
		log.Infof("Product %s not in store", productDiff.ID)
		_ = resp.WriteError(http.StatusNotFound, fmt.Errorf("product not found"))
		return
	}

	// update product in storage
	updated, err := api.storage.Update(productDiff.ID, *productDiff)

	if err != nil {
		log.Errorf("Failed to update product in storage, err=%v", err)
		_ = resp.WriteError(http.StatusInternalServerError, fmt.Errorf("failed to update product"))
		return
	}

	if !updated {
		log.Infof("Product %s not in store", productDiff.ID)
		_ = resp.WriteError(http.StatusNotFound, fmt.Errorf("product not found"))
		return
	} else {
		log.Infof("Product %s updated in store", productDiff.ID)
		// http response ok
		_ = resp.WriteAsJson("Product " + productDiff.ID + " updated in store")
	}
}

func (api *API) updateProductHTTPSingle(req *restful.Request, resp *restful.Response) {
	productDiff := &domain.ProductDiff{}
	err := req.ReadEntity(productDiff)

	if err != nil {
		log.Errorf("Failed to read product, err=%v", err)
		_ = resp.WriteError(http.StatusBadRequest, err)
		return
	}

	// check if id exists in storage
	_, exists, err := api.storage.Get(productDiff.ID)
	if err != nil {
		log.Errorf("Failed to get product from storage, err=%v", err)
		_ = resp.WriteError(http.StatusInternalServerError, fmt.Errorf("failed to get product from store"))
		return
	}

	if !exists {
		log.Infof("Product %s not in store", productDiff.ID)
		_ = resp.WriteError(http.StatusNotFound, fmt.Errorf("product not found"))
		return
	}

	// update product in storage
	updated, err := api.client.Update(productDiff.ID, *productDiff)

	if err != nil {
		log.Errorf("Failed to update product in storage, err=%v", err)
		_ = resp.WriteError(http.StatusInternalServerError, fmt.Errorf("failed to update product"))
		return
	}

	if !updated {
		log.Infof("Product %s not in store", productDiff.ID)
		_ = resp.WriteError(http.StatusNotFound, fmt.Errorf("product not found"))
		return
	} else {
		log.Infof("Product %s updated in store", productDiff.ID)
		// http response ok
		_ = resp.WriteAsJson("Product " + productDiff.ID + " updated in store")
	}
}

func (api *API) deleteProductMemorySingle(req *restful.Request, resp *restful.Response) {
	id := req.QueryParameter("id")
	if id == "" {
		log.Infof("No id provided in request")
		_ = resp.WriteError(http.StatusBadRequest, fmt.Errorf("id must be provided"))
		return
	}

	// delete product from api storage
	deleted, err := api.storage.Delete(id)
	if err != nil {
		log.Errorf("Failed to delete product from storage, err=%v", err)
		_ = resp.WriteError(http.StatusInternalServerError, fmt.Errorf("failed to delete product from store"))
		return
	}

	if !deleted {
		log.Infof("Product %s not in store", id)
		_ = resp.WriteError(http.StatusNotFound, fmt.Errorf("product not found"))
		return
	} else {
		log.Infof("Product %s deleted from store", id)
		_ = resp.WriteAsJson("Product " + id + " deleted from store")
	}

}

func (api *API) deleteProductHTTPSingle(req *restful.Request, resp *restful.Response) {
	id := req.QueryParameter("id")
	if id == "" {
		log.Infof("No id provided in request")
		_ = resp.WriteError(http.StatusBadRequest, fmt.Errorf("id must be provided"))
		return
	}

	// delete product from api storage
	deleted, err := api.client.Delete(id)
	if err != nil {
		log.Errorf("Failed to delete product from storage, err=%v", err)
		_ = resp.WriteError(http.StatusInternalServerError, fmt.Errorf("failed to delete product from store"))
		return
	}

	if !deleted {
		log.Infof("Product %s not in store", id)
		_ = resp.WriteError(http.StatusNotFound, fmt.Errorf("product not found"))
		return
	} else {
		log.Infof("Product %s deleted from store", id)
		_ = resp.WriteAsJson("Product " + id + " deleted from store")
	}

}
