package route

import (
	"github.com/HyeonHo-Park/simple-mutating-webhook/internal/handler"
	"github.com/gin-gonic/gin"
)

func InitDeploymentRoute(r *gin.Engine) {
	groupRoute := r.Group("/api/v1/deployment")

	dHandler := handler.NewDeploymentHandler()

	groupRoute.POST("/mutate", dHandler.Mutate)
}
