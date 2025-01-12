package handlers

import "goweb-template/common/webframe"

func StatusHandler(c webframe.HttpContext) {
	c.Success(nil)
}
