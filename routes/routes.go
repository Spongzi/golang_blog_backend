package routes

import (
	"github.com/gin-gonic/gin"
	"golang_blog_backend/logger"
)

func SetupRouter() *gin.Engine {
	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true))
	return r
}
