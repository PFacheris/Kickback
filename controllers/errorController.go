package controllers

import (
	"github.com/martini-contrib/render"
)

func HandleError(encoding string, status int, e error, r render.Render) {
	reportError(e)

	switch encoding {
	case "json":
		r.JSON(status, map[string]string{
			"message": e.Error(),
		})
	default:
		r.HTML(status, "error", map[string]string{
			"message": e.Error(),
		})
	}
}

// @TODO: Error reporting service and/or logging.
func reportError(e error) {

}
