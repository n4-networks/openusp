package rest

import (
	"log"
	"net/http"
)

func (re *Rest) setMiddlewares() error {
	log.Println("Registering middleware logging")
	re.router.Use(middlewareLogging)
	log.Println("Registering middleware access control")
	re.router.Use(middlewareUserAuth)
	return nil
}

func middlewareLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Middleware Logging: %v, %v\n", r.Method, r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
