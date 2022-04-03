package handlers

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"ms/data"
	"net/http"
	"strconv"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

func (p *Products) GetProducts(writer http.ResponseWriter, request *http.Request) {
	p.l.Println("Handle GET Products")
	lp := data.GetProducts()
	err := lp.ToJson(writer)
	if err != nil {
		http.Error(writer, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) AddProduct(writer http.ResponseWriter, request *http.Request) {
	p.l.Println("Handle POST Product")
	prod := request.Context().Value(KeyProduct).(data.Product)
	data.AddProduct(&prod)
}

func (p *Products) UpdateProduct(writer http.ResponseWriter, request *http.Request) {
	p.l.Println("Handle PUT Product")

	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(writer, "Unable to convert id", http.StatusBadRequest)
		return
	}

	prod := request.Context().Value(KeyProduct).(data.Product)

	err = data.UpdateProduct(id, &prod)
	if err == data.ErrorProductNotFound {
		http.Error(writer, "Product Not Found", http.StatusNotFound)
		return
	}

	if err != nil {
		http.Error(writer, "Product Not Found", http.StatusInternalServerError)
		return
	}
}

var KeyProduct = fmt.Sprintf("product")

func (p Products) MiddlewareProductValidation(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		prod := data.Product{}
		err := prod.FromJson(request.Body)
		if err != nil {
			p.l.Println("Failed Convert From Json Product ", err)
			http.Error(writer, "Unable to unmarshal json", http.StatusBadRequest)
			return
		}

		ctx := context.WithValue(request.Context(), KeyProduct, prod)
		req := request.WithContext(ctx)
		next.ServeHTTP(writer, req)
	})
}
