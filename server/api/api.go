package api

import (
	"os"

	"github.com/2432001677/chat-gopt/server/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func StartServer() {

	r := gin.New()
	r.Use(cors.Default())

	g := r.Group("/api/v1")
	{
		g.POST("ask", service.Ask)
		g.GET("history", service.History)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8060"
	}
	r.Run(":" + port)
}
