package handlers

import (
	"log"
	"ms/data"
	"net/http"
	"regexp"
	"strconv"
)

type Products struct {
	l *log.Logger
}

func (p *Products) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		p.getProducts(writer, request)
		return
	}

	if request.Method == http.MethodPost {
		p.addProduct(writer, request)
		return
	}

	if request.Method == http.MethodPut {
		regex := regexp.MustCompile(`/([0-9]+)`)
		group := regex.FindAllStringSubmatch(request.URL.Path, -1)
		if len(group) != 1 {
			http.Error(writer, "Invalid URI", http.StatusBadRequest)
			return
		}

		if len(group[0]) != 2 {
			http.Error(writer, "Invalid URI", http.StatusBadRequest)
			return
		}

		idString := group[0][1]
		id, err := strconv.Atoi(idString)

		if err != nil {
			http.Error(writer, "Invalid URI", http.StatusBadRequest)
			return
		}

		p.updateProduct(id, writer, request)
		return
	}

	writer.WriteHeader(http.StatusMethodNotAllowed)
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) getProducts(writer http.ResponseWriter, request *http.Request) {
	p.l.Println("Handle GET Products")
	lp := data.GetProducts()
	err := lp.ToJson(writer)
	if err != nil {
		http.Error(writer, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) addProduct(writer http.ResponseWriter, request *http.Request) {
	p.l.Println("Handle POST Product")
	prod := &data.Product{}
	err := prod.FromJson(request.Body)
	if err != nil {
		http.Error(writer, "Unable to unmarshal json", http.StatusBadRequest)
	}
	data.AddProduct(prod)
}

func (p *Products) updateProduct(id int, writer http.ResponseWriter, request *http.Request) {
	p.l.Println("Handle PUT Product")

	prod := &data.Product{}
	err := prod.FromJson(request.Body)
	if err != nil {
		http.Error(writer, "Unable to unmarshal json", http.StatusBadRequest)
	}
	err = data.UpdateProduct(id, prod)
	if err == data.ErrorProductNotFound {
		http.Error(writer, "Product Not Found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(writer, "Product Not Found", http.StatusInternalServerError)
		return
	}
}
