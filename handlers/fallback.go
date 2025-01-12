package handlers

import (
	"goweb-template/common/webframe"
	"net/http"
)

func FallbackHandler(c webframe.HttpContext) {
	c.HttpError(http.StatusNotFound)
}
