package apiserver

import (
	"log"
	"net/http"
)

func (as *ApiServer) setMiddlewares() error {
	log.Println("Registering middleware logging")
	as.router.Use(middlewareLogging)
	log.Println("Registering middleware access control")
	as.router.Use(middlewareUserAuth)
	return nil
}

func middlewareLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Middleware Logging: %v, %v\n", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
