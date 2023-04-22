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

package main

import (
	"log"

	"github.com/n4-networks/openusp/pkg/apiserver"
)

func main() {
	log.SetFlags(log.Lshortfile)

	as := &apiserver.ApiServer{}
	log.Println("Initializing API Server...")
	if err := as.Init(); err != nil {
		log.Println("Error:", err)
	}
	log.Println("Starting API Server...")
	if err := as.Server(); err != nil {
		log.Println("Error: API Server is exiting...Err:", err)
	}
}
