package main

import (
	"github.com/op/go-logging"
	"os"
)

var logger = logging.MustGetLogger("consul-backup")
var format = logging.MustStringFormatter(`%{time:2006-01-02 15:04:05.000} %{color}%{level}%{color:reset} (consul-backup):  %{message}`)

func SetupLogging() {
	logging.Reset()
	handler := logging.NewLogBackend(os.Stdout, "", 0)
	formatter := logging.NewBackendFormatter(handler, format)
	logging.SetBackend(formatter)
}
