// Copyright 2023 N4-Networks.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
