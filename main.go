package main

import (
	"context"
	"fmt"
	"github.com/go-openapi/runtime/middleware"
	goHandlers "github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/nicholasjackson/env"
	"log"
	"ms/data"
	"ms/handlers"
	"net/http"
	"os"
	"os/signal"
	"time"
)

var bindAddress = env.Int("BIND_ADDRESS", false, 9090, "Bind address for the server")
var logLevel = env.String("LOG_LEVEL", false, "debug", "Log output level for the server [debug, info, trace]")
var basePath = env.String("BASE_PATH", false, "/tmp/images", "Base path to save images")

func main() {
	_ = env.Parse()

	l := log.New(os.Stdout, "products-api ", log.LstdFlags)

	//logImages := hclog.New(
	//	&hclog.LoggerOptions{
	//		Name:  "product-images",
	//		Level: hclog.LevelFromString(*logLevel),
	//	},
	//)
	//
	//sl := logImages.StandardLogger(&hclog.StandardLoggerOptions{InferLevels: true})
	//stor, err := files.NewLocal(*basePath, 1024*1000*5)
	//if err != nil {
	//	logImages.Error("Unable to create storage", "error", err)
	//	os.Exit(1)
	//}

	// create the handlers
	//fh := handlers.NewFiles(stor, logImages)

	v := data.NewValidation()

	// create the handlers
	ph := handlers.NewProducts(l, v)

	// create a new serve mux and register the handlers
	sm := mux.NewRouter()

	// handlers for API
	getR := sm.Methods(http.MethodGet).Subrouter()
	getR.HandleFunc("/products", ph.ListAll)
	getR.HandleFunc("/products/{id:[0-9]+}", ph.ListSingle)

	//getR.Handle(
	//	"/images/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-z]{3}}",
	//	http.StripPrefix("/images/", http.FileServer(http.Dir(*basePath))),
	//)

	putR := sm.Methods(http.MethodPut).Subrouter()
	putR.HandleFunc("/products", ph.Update)
	putR.Use(ph.MiddlewareValidateProduct)

	postR := sm.Methods(http.MethodPost).Subrouter()
	postR.HandleFunc("/products", ph.Create)
	postR.Use(ph.MiddlewareValidateProduct)
	//
	//postR.HandleFunc("/images/{id:[0-9]+}/{filename:[a-zA-Z]+\\.[a-z]{3}}", fh.ServeHTTP)

	deleteR := sm.Methods(http.MethodDelete).Subrouter()
	deleteR.HandleFunc("/products/{id:[0-9]+}", ph.Delete)

	// handler for documentation
	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)

	getR.Handle("/docs", sh)
	getR.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	gh := goHandlers.CORS(goHandlers.AllowedOrigins([]string{"*"}))

	// create a new server
	s := http.Server{
		Addr:         fmt.Sprintf(":%v", *bindAddress), // configure the bind address
		Handler:      gh(sm),                           // set the default handler
		ErrorLog:     l,                                // set the logger for the server
		ReadTimeout:  5 * time.Second,                  // max time to read request from the sdk
		WriteTimeout: 10 * time.Second,                 // max time to write response to the sdk
		IdleTimeout:  120 * time.Second,                // max time for connections using TCP Keep-Alive
	}

	// start the server
	go func() {
		l.Println(hclog.Debug,
			fmt.Sprint("Starting server on port ", *bindAddress))
		l.Println(hclog.Debug,
			fmt.Sprint("Trace level server ", *logLevel))

		err := s.ListenAndServe()
		if err != nil {
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	// trap sigterm or interupt and gracefully shutdown the server
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	// Block until a signal is received.
	sig := <-c
	log.Println("Got signal:", sig)

	// gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	s.Shutdown(ctx)
}
