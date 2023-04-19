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
