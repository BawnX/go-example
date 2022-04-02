package handlers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type Hello struct {
	l *log.Logger
}

func (h *Hello) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	h.l.Println("Hola Mundo")

	d, err := ioutil.ReadAll(request.Body)
	if err != nil {
		http.Error(writer, "Oooops", http.StatusBadRequest)
		return
	}

	fmt.Fprintf(writer, "Hola %s\n", d)
}

func NewHello(l *log.Logger) *Hello {
	return &Hello{l}
}
