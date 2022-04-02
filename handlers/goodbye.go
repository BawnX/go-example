package handlers

import (
	"log"
	"net/http"
)

type Goodbye struct {
	l *log.Logger
}

func (g Goodbye) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	writer.Write([]byte("Byeee"))
}

func NewGoodbye(l *log.Logger) *Goodbye {
	return &Goodbye{l}
}
