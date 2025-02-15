package handler

import (
	"goweb/common/utility"
	"goweb/common/webframe"
	"goweb/handler/api"
	"goweb/handler/system"
	"goweb/web"
	"io/fs"
	"net/http"
	"os"

	"github.com/Streamlet/gohttp"
)

func registerWebHandlers(router gohttp.Application[webframe.HttpContext], webroot string) {
	var rootFs, includeFs fs.FS
	if webroot != "" {
		rootFs = os.DirFS(webroot)
		includeFs = rootFs
	} else {
		rootFs = web.RootFs
		includeFs = web.IncludeFs
	}
	ssiFs := utility.NewSsiFS(rootFs, includeFs, []string{".html"}, "<!-- #include=\"(.*)\" -->")
	router.RawHandle("/", http.FileServer(http.FS(ssiFs)))
}

func registerApiHandlers(router gohttp.Application[webframe.HttpContext]) {
	// BEGIN AUTO GENERATED API HANDLERS
	router.Handle("/api/version", api.VersionHandler)
	// END AUTO GENERATED API HANDLERS
}

func registerSystemHandlers(router gohttp.Application[webframe.HttpContext]) {
	router.Handle("/status", system.StatusHandler)
	// router.Handle("/", system.FallbackHandler)
}

func Registers(router gohttp.Application[webframe.HttpContext], webroot string) {
	registerApiHandlers(router)
	registerSystemHandlers(router)
	registerWebHandlers(router, webroot)
}
