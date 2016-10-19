package main

import (
	"github.com/op/go-logging"
	"os"
)

var logger = logging.MustGetLogger("consul-bak")
var format = logging.MustStringFormatter(`%{time:2006-01-02 15:04:05.000} %{color}%{level}%{color:reset} (consul-bak.%{shortfile}):  %{message}`)

// SetupLogging setup up logging infra, should be called once on start
func SetupLogging() {
	logging.Reset()
	handler := logging.NewLogBackend(os.Stdout, "", 0)
	formatter := logging.NewBackendFormatter(handler, format)
	logging.SetBackend(formatter)
}
