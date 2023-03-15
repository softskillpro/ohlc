package data

import (
	"github.com/gin-gonic/gin"
)

func RouterRegister(route gin.RouterGroup) {
	routeGroup := route.Group(routeContext)

	routeGroup.GET("", Get)
	routeGroup.POST("", Post)
}
