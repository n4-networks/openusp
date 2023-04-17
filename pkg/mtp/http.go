package mtp

import (
	"errors"
	"fmt"
	"log"
	"net/http"
)

func (m *Mtp) HttpServerStart(exit chan int32) (err error) {

	if m.Cfg == nil {
		log.Panic("HTTP Configuration are not set")
		return errors.New("HTTP Config has not been set yet")
	}
	//http.HandleFunc("/cwmp", CwmpHandler)

	switch m.Cfg.Http.Mode {
	case "nontls":
		log.Println("Starting HTTP server in NonTLS mode")
		go m.HttpServer(exit)
	case "tls":
		log.Println("Starting HTTP server in TLS mode")
		go m.HttpServerTLS(exit)
	case "both":
		log.Println("Starting HTTP server in both nonTLS & TLS mode")
		go m.HttpServer(exit)
		go m.HttpServerTLS(exit)
	}
	return nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func (m *Mtp) HttpServer(exit chan int32) {

	//http.HandleFunc("/", handler)
	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))
	addr := ":" + m.Cfg.Http.Port
	log.Println("Starting HTTP Server Instance at: ", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		log.Panic("Could not invoke http server instance", err.Error())
		exit <- HTTP_SERVER
	}
}

func (m *Mtp) HttpServerTLS(exit chan int32) {
	addr := ":" + m.Cfg.Http.TLSPort
	log.Println("Starting HTTP TLS Server Instance at: ", addr)
	err := http.ListenAndServeTLS(addr, m.Cfg.Http.CertFile, m.Cfg.Http.KeyFile, nil)
	if err != nil {
		log.Panic("Could not invoke https server instance", err.Error())
		exit <- HTTP_SERVER_TLS
	}
}
