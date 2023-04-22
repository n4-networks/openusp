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

package clitest

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
)

func runAndReport(cmd string, cmpStr string) (*bytes.Buffer, *bytes.Buffer, error) {

	logStream := &bytes.Buffer{}
	shellStream := &bytes.Buffer{}
	var err error = nil

	log.SetOutput(logStream)
	defer log.SetOutput(os.Stdout)

	cliH.SetOut(shellStream)
	defer cliH.SetOut(os.Stdout)

	if err = cliH.ProcessCmd(cmd); err == nil {
		if !strings.Contains(shellStream.String(), cmpStr) {
			err = fmt.Errorf("Error: Output does not have: %v\n", cmpStr)
		} else if strings.Contains(logStream.String(), "Error") {
			err = fmt.Errorf("Error: From log\n")
		}
	}
	return shellStream, logStream, err
}

func runAndCheckErr(t *testing.T, cmd string) {
	cliH.ClearLastCmdErr()
	if err := cliH.ProcessCmd(cmd); err != nil {
		t.Fatal(err)
	}
	if err := cliH.GetLastCmdErr(); err != nil {
		t.Error(err)
	}
}

func getInstancePathByAlias(name string) (string, error) {
	return cliH.GetInstancePathByAlias(name)
}

func runAndReturn(cmd string) error {
	cliH.ClearLastCmdErr()
	if err := cliH.ProcessCmd(cmd); err != nil {
		return err
	}
	return cliH.GetLastCmdErr()
}
