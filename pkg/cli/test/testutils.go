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
