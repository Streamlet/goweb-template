package main

import (
	"github.com/Streamlet/gohttp"
	"goweb-template/common/webframe"
	"goweb-template/handlers"
)

func registerHandlers(application gohttp.Application[webframe.HttpContext]) {
	application.Handle("/", handlers.FallbackHandler)
	application.Handle("/status", handlers.StatusHandler)
}
