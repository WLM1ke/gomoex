package gomoex

import (
	"fmt"
	"go.uber.org/goleak"
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	runWithLeakDetector(m, http.DefaultClient.CloseIdleConnections)
}

func runWithLeakDetector(m *testing.M, teardownFunc func()) {
	exitCode := m.Run()
	teardownFunc()

	if exitCode == 0 {
		if err := goleak.Find(); err != nil {
			fmt.Printf("%v\n\nleaks on successful test run\n", err)
			exitCode = 1
		}
	}

	os.Exit(exitCode)
}
