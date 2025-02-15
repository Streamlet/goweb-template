package system

import (
	"goweb/common/webframe"
	"net/http"
)

func FallbackHandler(c webframe.HttpContext) {
	c.HttpError(http.StatusNotFound)
}
