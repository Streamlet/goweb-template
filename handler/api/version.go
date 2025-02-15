package api

import (
	"goweb/common/webframe"
)

type versionResonse struct {
	Version string `json:"version"`
}

func VersionHandler(c webframe.HttpContext) {
	c.Success(versionResonse{
		Version: "1.0.0",
	})
}
