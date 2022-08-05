package api

import (
	"encoding/json"
	"exam-api/domain"
	"fmt"
	"github.com/emicklei/go-restful/v3"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
)

func (api *API) createProductMemoryBatch(req *restful.Request, resp *restful.Response) {
	body := req.Request.Body
	if body == nil {
		log.Printf("[ERROR] Couldn't read request body")
		resp.WriteError(http.StatusInternalServerError, restful.NewError(http.StatusInternalServerError, "nil body\n"))
		return

	}
	defer body.Close()
	var err error
	if err != nil {
		log.Printf("[ERROR] Couldn't read request body")
		resp.WriteError(http.StatusInternalServerError, restful.NewError(http.StatusInternalServerError, err.Error()+"\n"))
		return
	}
	data, err := ioutil.ReadAll(body)
	if err != nil {
		log.Printf("[ERROR] Couldn't read request body")
		resp.WriteError(http.StatusInternalServerError, restful.NewError(http.StatusInternalServerError, err.Error()+"\n"))
		return
	}

	// parse list of products from data
	var products []*domain.Product
	if err := json.Unmarshal(data, &products); err != nil {
		log.Printf("[ERROR] Couldn't parse request body")
		resp.WriteError(http.StatusInternalServerError, restful.NewError(http.StatusInternalServerError, err.Error()+"\n"))
		return
	}

	wg := &sync.WaitGroup{}

	// iterate over products
	for _, product := range products {
		wg.Add(1)
		go func(product *domain.Product) {
			defer wg.Done()
			id, alreadyExists, err := api.storage.Save(*product)
			if err != nil {
				log.Printf("[ERROR] Couldn't save product in storage")
				resp.WriteError(http.StatusInternalServerError, restful.NewError(http.StatusInternalServerError, err.Error()+"\n"))
				return
			}

			if alreadyExists {
				log.Printf("[ERROR] Product %s already in store", id)
				resp.WriteError(http.StatusConflict, restful.NewError(http.StatusConflict, "product already exists\n"))
				return
			}

			log.Printf("[INFO] Product %s saved in store", id)
		}(product)
	}

	wg.Wait()
	resp.WriteAsJson(products)
}

func (api *API) createProductHTTPBatch(req *restful.Request, resp *restful.Response) {
	body := req.Request.Body
	if body == nil {
		log.Printf("[ERROR] Couldn't read request body")
		resp.WriteError(http.StatusInternalServerError, restful.NewError(http.StatusInternalServerError, "nil body\n"))
		return

	}
	defer body.Close()
	var err error
	if err != nil {
		log.Printf("[ERROR] Couldn't read request body")
		resp.WriteError(http.StatusInternalServerError, restful.NewError(http.StatusInternalServerError, err.Error()+"\n"))
		return
	}
	data, err := ioutil.ReadAll(body)
	if err != nil {
		log.Printf("[ERROR] Couldn't read request body")
		resp.WriteError(http.StatusInternalServerError, restful.NewError(http.StatusInternalServerError, err.Error()+"\n"))
		return
	}

	// parse list of products from data
	var products []*domain.Product
	if err := json.Unmarshal(data, &products); err != nil {
		log.Printf("[ERROR] Couldn't parse request body")
		resp.WriteError(http.StatusInternalServerError, restful.NewError(http.StatusInternalServerError, err.Error()+"\n"))
		return
	}

	wg := &sync.WaitGroup{}

	// iterate over products
	for _, product := range products {
		wg.Add(1)
		go func(product *domain.Product) {
			defer wg.Done()
			id, alreadyExists, err := api.client.Save(*product)
			if err != nil {
				log.Printf("[ERROR] Couldn't save product in storage")
				resp.WriteError(http.StatusInternalServerError, restful.NewError(http.StatusInternalServerError, err.Error()+"\n"))
				return
			}

			if alreadyExists {
				log.Printf("[ERROR] Product %s already in store", id)
				resp.WriteError(http.StatusConflict, restful.NewError(http.StatusConflict, "product already exists\n"))
				return
			}

			log.Printf("[INFO] Product %s saved in store", id)
		}(product)
	}

	wg.Wait()
	resp.WriteAsJson(products)
}

func (api *API) getProductMemoryBatch(req *restful.Request, resp *restful.Response) {
	ids := req.QueryParameters("id")
	if ids == nil {
		log.Printf("[ERROR] Failed to read id")
		resp.WriteError(http.StatusBadRequest, fmt.Errorf("element id must be provided\n"))
		return
	}

	wg := &sync.WaitGroup{}
	var mu sync.Mutex
	var finalResponse string = ""

	for _, id := range ids {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			product, ok, err := api.storage.Get(id)
			if err != nil {
				log.Printf("[ERROR] Couldn't get product from storage")
				resp.WriteError(http.StatusInternalServerError, restful.NewError(http.StatusInternalServerError, err.Error()+"\n"))
				return
			}

			if !ok {
				log.Printf("[ERROR] Product %s not found", id)
				resp.WriteError(http.StatusNotFound, restful.NewError(http.StatusNotFound, "product not found\n"))
				return
			}

			mu.Lock()

			finalResponse += fmt.Sprintf("%v\n", product)
			mu.Unlock()
		}(id)
	}

	wg.Wait()
	resp.Write([]byte(finalResponse))
}

func (api *API) getProductHTTPBatch(req *restful.Request, resp *restful.Response) {
	ids := req.QueryParameters("id")
	if ids == nil {
		log.Printf("[ERROR] Failed to read id")
		resp.WriteError(http.StatusBadRequest, fmt.Errorf("element id must be provided\n"))
		return
	}

	wg := &sync.WaitGroup{}
	var mu sync.Mutex
	var finalResponse string = ""

	for _, id := range ids {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			product, ok, err := api.client.Get(id)

			if err != nil {
				log.Printf("[ERROR] Couldn't get product from storage")
				resp.WriteError(http.StatusInternalServerError, restful.NewError(http.StatusInternalServerError, err.Error()+"\n"))
				return
			}

			if !ok {
				log.Printf("[ERROR] Product %s not found", id)
				resp.WriteError(http.StatusNotFound, restful.NewError(http.StatusNotFound, "product not found\n"))
				return
			}

			mu.Lock()

			finalResponse += fmt.Sprintf("%v\n", product)
			mu.Unlock()
		}(id)
	}

	wg.Wait()
	resp.Write([]byte(finalResponse))
}

func (api *API) updateProductMemoryBatch(req *restful.Request, resp *restful.Response) {
	body := req.Request.Body
	if body == nil {
		log.Printf("[ERROR] Couldn't read request body")
		resp.WriteError(http.StatusInternalServerError, restful.NewError(http.StatusInternalServerError, "nil body\n"))
		return

	}
	defer body.Close()
	var err error
	if err != nil {
		log.Printf("[ERROR] Couldn't read request body")
		resp.WriteError(http.StatusInternalServerError, restful.NewError(http.StatusInternalServerError, err.Error()+"\n"))
		return
	}
	data, err := ioutil.ReadAll(body)
	if err != nil {
		log.Printf("[ERROR] Couldn't read request body")
		resp.WriteError(http.StatusInternalServerError, restful.NewError(http.StatusInternalServerError, err.Error()+"\n"))
		return
	}

	// parse list of productDiff from data
	var productDiffs []*domain.ProductDiff
	if err := json.Unmarshal(data, &productDiffs); err != nil {
		log.Printf("[ERROR] Couldn't parse request body")
		resp.WriteError(http.StatusInternalServerError, restful.NewError(http.StatusInternalServerError, err.Error()+"\n"))
		return
	}
	var finalResponse string = ""
	wg := &sync.WaitGroup{}

	// iterate over productDiffs
	for _, productDiff := range productDiffs {
		wg.Add(1)
		go func(productDiff *domain.ProductDiff) {
			defer wg.Done()
			ok, err := api.storage.Update(productDiff.ID, *productDiff)
			if err != nil {
				log.Printf("[ERROR] Couldn't update product in storage")
				resp.WriteError(http.StatusInternalServerError, restful.NewError(http.StatusInternalServerError, err.Error()+"\n"))
				return
			}

			if !ok {
				log.Printf("[ERROR] Product %s not found", productDiff.ID)
				resp.WriteError(http.StatusNotFound, restful.NewError(http.StatusNotFound, "product not found\n"))
				return
			}

			finalResponse += fmt.Sprintf("%v\n", productDiff)
			log.Printf("[INFO] Product %s updated in store", productDiff.ID)
		}(productDiff)
	}

	wg.Wait()
	resp.Write([]byte(finalResponse))
}

func (api *API) updateProductHTTPBatch(req *restful.Request, resp *restful.Response) {
	body := req.Request.Body
	if body == nil {
		log.Printf("[ERROR] Couldn't read request body")
		resp.WriteError(http.StatusInternalServerError, restful.NewError(http.StatusInternalServerError, "nil body\n"))
		return

	}
	defer body.Close()
	var err error
	if err != nil {
		log.Printf("[ERROR] Couldn't read request body")
		resp.WriteError(http.StatusInternalServerError, restful.NewError(http.StatusInternalServerError, err.Error()+"\n"))
		return
	}
	data, err := ioutil.ReadAll(body)
	if err != nil {
		log.Printf("[ERROR] Couldn't read request body")
		resp.WriteError(http.StatusInternalServerError, restful.NewError(http.StatusInternalServerError, err.Error()+"\n"))
		return
	}

	// parse list of productDiff from data
	var productDiffs []*domain.ProductDiff
	if err := json.Unmarshal(data, &productDiffs); err != nil {
		log.Printf("[ERROR] Couldn't parse request body")
		resp.WriteError(http.StatusInternalServerError, restful.NewError(http.StatusInternalServerError, err.Error()+"\n"))
		return
	}
	var finalResponse string = ""
	wg := &sync.WaitGroup{}

	// iterate over productDiffs
	for _, productDiff := range productDiffs {
		wg.Add(1)
		go func(productDiff *domain.ProductDiff) {
			defer wg.Done()
			ok, err := api.client.Update(productDiff.ID, *productDiff)
			if err != nil {
				log.Printf("[ERROR] Couldn't update product in storage")
				resp.WriteError(http.StatusInternalServerError, restful.NewError(http.StatusInternalServerError, err.Error()+"\n"))
				return
			}

			if !ok {
				log.Printf("[ERROR] Product %s not found", productDiff.ID)
				resp.WriteError(http.StatusNotFound, restful.NewError(http.StatusNotFound, "product not found\n"))
				return
			}

			finalResponse += fmt.Sprintf("%v\n", productDiff)
			log.Printf("[INFO] Product %s updated in store", productDiff.ID)
		}(productDiff)
	}

	wg.Wait()
	resp.Write([]byte(finalResponse))
}

func (api *API) deleteProductMemoryBatch(req *restful.Request, resp *restful.Response) {
	ids := req.QueryParameters("id")
	if ids == nil {
		log.Printf("[ERROR] Failed to read id")
		resp.WriteError(http.StatusBadRequest, fmt.Errorf("element id must be provided\n"))
		return
	}

	wg := &sync.WaitGroup{}
	var mu sync.Mutex
	var finalResponse string = ""

	for _, id := range ids {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			ok, err := api.storage.Delete(id)
			if err != nil {
				log.Printf("[ERROR] Couldn't delete product from storage")
				resp.WriteError(http.StatusInternalServerError, restful.NewError(http.StatusInternalServerError, err.Error()+"\n"))
				return
			}

			if !ok {
				log.Printf("[ERROR] Product %s not found", id)
				resp.WriteError(http.StatusNotFound, restful.NewError(http.StatusNotFound, "product not found\n"))
				return
			}

			mu.Lock()

			finalResponse += fmt.Sprintf("%v\n", id)
			mu.Unlock()
		}(id)
	}

	wg.Wait()
	resp.Write([]byte(finalResponse))
}

func (api *API) deleteProductHTTPBatch(req *restful.Request, resp *restful.Response) {
	ids := req.QueryParameters("id")
	if ids == nil {
		log.Printf("[ERROR] Failed to read id")
		resp.WriteError(http.StatusBadRequest, fmt.Errorf("element id must be provided\n"))
		return
	}

	wg := &sync.WaitGroup{}
	var mu sync.Mutex
	var finalResponse string = ""

	for _, id := range ids {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			ok, err := api.client.Delete(id)
			if err != nil {
				log.Printf("[ERROR] Couldn't delete product from storage")
				resp.WriteError(http.StatusInternalServerError, restful.NewError(http.StatusInternalServerError, err.Error()+"\n"))
				return
			}

			if !ok {
				log.Printf("[ERROR] Product %s not found", id)
				resp.WriteError(http.StatusNotFound, restful.NewError(http.StatusNotFound, "product not found\n"))
				return
			}

			mu.Lock()

			finalResponse += fmt.Sprintf("%v\n", id)
			mu.Unlock()
		}(id)
	}

	wg.Wait()
	resp.Write([]byte(finalResponse))
}
