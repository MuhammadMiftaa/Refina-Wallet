package service

import (
	"io"
	"os"
	"testing"

	"refina-wallet/config/log"

	"github.com/sirupsen/logrus"
)

// TestMain initializes shared test dependencies before running any tests.
func TestMain(m *testing.M) {
	// Initialize logger so that log.Info / log.Error / etc. don't panic
	log.Log = logrus.New()
	log.Log.SetOutput(io.Discard)
	log.Log.SetLevel(logrus.PanicLevel)

	os.Exit(m.Run())
}
