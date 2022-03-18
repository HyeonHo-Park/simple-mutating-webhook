package main

import (
	"github.com/HyeonHo-Park/simple-mutating-webhook/internal/route"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main()  {
	router := setupRouter()
	router.Run(":8080")
}

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"*"},
		AllowHeaders:  []string{"*"},
		AllowWildcard: true,
	}))

	route.InitDeploymentRoute(r)

	return r
}