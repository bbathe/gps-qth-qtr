package main

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	// don't output normal log messages
	log.SetOutput(ioutil.Discard)

	os.Exit(m.Run())
}
