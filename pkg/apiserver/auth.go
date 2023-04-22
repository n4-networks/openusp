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
)

var users = map[string]string{
	"n4admin": "n4defaultpass",
}

func isAuthorized(username, password string) bool {
	pass, ok := users[username]
	if !ok {
		return false
	}
	return pass == password

}

func middlewareUserAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "OPTIONS" {
			next.ServeHTTP(w, r)
			return
		}
		log.Println(r.RequestURI)
		username, password, ok := r.BasicAuth()
		if !ok {
			w.Header().Add("WWW-Authenticate", `Basic realm="Give username and password"`)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"message" : "No basic auth present"}`))
			//w.Header().Set("Access-Control-Allow-Origin", "*")  // require for UI to avoid CORS Policy
			//w.Header().Set("Access-Control-Allow-Headers", "*") // require for UI to avoid CORS Policy
			log.Println("No basic auth present")
			return
		}
		if !isAuthorized(username, password) {
			w.Header().Add("WWW-Authenticate", `Basic realm="Give username and password"`)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(`{"message" : "Invalid username and password"}`))
			//w.Header().Set("Access-Control-Allow-Origin", "*")  // require for UI to avoid CORS Policy
			//w.Header().Set("Access-Control-Allow-Headers", "*") // require for UI to avoid CORS Policy
			log.Println("Invalid username and password")
			return
		}
		log.Println("Passed Authorization test")
		next.ServeHTTP(w, r)
	})
}
