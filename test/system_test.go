//go:build component
// +build component

//nolint:testpackage
package test

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/gavv/httpexpect/v2"
)

func TestMain(m *testing.M) {
	if isRunningInDockerContainer() {
		if err := StartPokerContainer(); err != nil {
			fmt.Printf("start container failed: %v\n", err)
			os.Exit(1)
		}
	} else {
		pokerHost = "localhost"
		pokerPort = "8080"
	}

	os.Exit(m.Run())
}

func TestAlive(t *testing.T) {
	e := httpexpect.New(t, "http://"+pokerHost+":"+pokerPort)
	e.GET("/alive").Expect().Status(http.StatusOK)
}
