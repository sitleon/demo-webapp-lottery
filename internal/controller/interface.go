package controller

import "github.com/gin-gonic/gin"

type RestController interface {
	RegisterRoute(root *gin.RouterGroup)
}
