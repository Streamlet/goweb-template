package system

import "goweb/common/webframe"

func StatusHandler(c webframe.HttpContext) {
	c.Success(nil)
}
