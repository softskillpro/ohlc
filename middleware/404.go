package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"ohcl/controllers/outputForms"
)

const (
	notFoundPath   = "not found"
	notFoundMethod = "not found method"
)

func NotFoundPath(c *gin.Context) {
	c.JSON(http.StatusNotFound, outputForms.NewState().
		SetStatus(false).SetCode(http.StatusNotFound).SetMessage(notFoundPath))
}

func NotFoundMethod(c *gin.Context) {
	c.JSON(http.StatusMethodNotAllowed, outputForms.NewState().
		SetCode(http.StatusMethodNotAllowed).SetStatus(false).SetMessage(notFoundMethod))
}
