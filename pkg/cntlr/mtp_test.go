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

package cntlr

import (
	"fmt"
	"log"
	"testing"
	"time"
)

func TestMtpStart(t *testing.T) {
	var cfg MtpConfig
	cfg.Http.Port = "22"
	cfg.Http.Mode = "nontls"

	fmt.Printf("Starting TestMtpStart\n")
	log.Printf("Printing from log\n")
	exit, err := MtpStart(&cfg)
	if err != nil {
		t.Fatal("Error in http server: ")
	}
	timer := time.NewTimer(2 * time.Second)
	select {
	case <-exit:
		t.Fatal("Error in http server: ")
	case <-timer.C:
		t.Logf("Timer expired")
	}
}
