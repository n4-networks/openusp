package apiserver

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/handlers"
)

func (as *ApiServer) Server() error {

	// CORS handlers
	headers := handlers.AllowedHeaders([]string{"content-type", "authorization"})
	origins := handlers.AllowedOrigins([]string{"*"})
	methods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	srv := &http.Server{
		Handler:      handlers.CORS(headers, origins, methods)(as.router),
		Addr:         ":" + as.cfg.httpPort,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Println("Starting HTTP server at:", as.cfg.httpPort)
	if as.cfg.isTlsOn {
		log.Println("Running server with TLS ...")
		return srv.ListenAndServeTLS("ssl/server.csr", "ssl/server.key")
	} else {
		return srv.ListenAndServe()
	}
	//return http.ListenAndServe(":8081", handlers.CORS(headers, origins, methods)(as.router))
}
