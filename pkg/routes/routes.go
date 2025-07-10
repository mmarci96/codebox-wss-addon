package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/mmarci96/codebox-wss-addon/pkg/config"
	"github.com/mmarci96/codebox-wss-addon/pkg/handlers"
)

func SetupRouter(cfg *config.Config) *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Logger())

	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set(
			"Access-Control-Allow-Methods",
			"POST, GET, OPTIONS, PUT, DELETE",
		)
		c.Writer.Header().Set(
			"Access-Control-Allow-Headers",
			"Content-Type, Authorization",
		)
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	})
	cm := handlers.NewClientManager()

	config := config.Load()

	r.GET(cfg.WsEndpoint, handlers.WsHandler(cm, config))
	r.GET("/ping", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "pong")
	})

	return r
}
