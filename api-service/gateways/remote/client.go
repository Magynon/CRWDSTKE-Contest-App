package remote

import (
	"encoding/json"
	"exam-api/domain"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// Implement the following client to connect to the remote storage server

type Client struct {
	client http.Client
}

func NewClient(client http.Client) *Client {
	return &Client{client: client}
}

func (c *Client) Save(product domain.Product) (string, bool, error) {
	url := "http://localhost:8081/store/product"
	method := "POST"

	client := &http.Client{}

	marshalledProduct, err := json.Marshal(product)

	if err != nil {
		return "", false, err
	}
	req, err := http.NewRequest(method, url, strings.NewReader(string(marshalledProduct)))

	if err != nil {
		fmt.Println(err)
		return "", false, err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", false, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "", false, err
	}

	fmt.Println(string(body))
	return product.GetHash(), false, nil
}

func (c *Client) Get(id string) (domain.Product, bool, error) {
	url := "http://localhost:8081/store/product/"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url+"?id="+id, nil)

	if err != nil {
		fmt.Println(err)
		return domain.Product{}, false, err
	}
	req.Header.Add("Content-Type", "application/json")

	res, _ := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return domain.Product{}, false, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return domain.Product{}, false, err
	}

	// unmarshal the body into a product
	var product domain.Product
	err = json.Unmarshal(body, &product)
	if err != nil {
		fmt.Println(err)
		return domain.Product{}, false, err
	}
	return product, true, nil
}

func (c *Client) Update(id string, diff domain.ProductDiff) (bool, error) {
	url := "http://localhost:8081/store/product"
	method := "PATCH"

	client := &http.Client{}

	// initialize new product
	newProduct := domain.Product{
		Name:         id,
		Manufacturer: "",
		Price:        diff.Diff.Price,
		Stock:        diff.Diff.Stock,
		Tags:         diff.Diff.Tags,
	}

	marshalledProduct, err := json.Marshal(newProduct)

	req, err := http.NewRequest(method, url, strings.NewReader(string(marshalledProduct)))

	if err != nil {
		fmt.Println(err)
		return false, err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	defer res.Body.Close()

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return false, err
	}

	return true, nil
}

func (c *Client) Delete(id string) (bool, error) {
	url := "http://localhost:8081/store/product/"
	method := "DELETE"

	client := &http.Client{}
	req, err := http.NewRequest(method, url+"?id="+id, nil)

	if err != nil {
		fmt.Println(err)
		return false, err
	}
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return false, err
	}
	fmt.Println(string(body))

	return true, nil
}
