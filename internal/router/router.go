package router

import (
	"github.com/gin-gonic/gin"

	"github.com/smolse/fluffy-pancake/internal/handler"
	"github.com/smolse/fluffy-pancake/internal/service"
)

// NewRouter returns a new Gin router with configured routes.
func NewRouter(svc *service.RiskService) *gin.Engine {
	handler := handler.NewHandler(svc)

	router := gin.Default()
	v1 := router.Group("/v1")
	{
		risks := v1.Group("/risks")
		{
			risks.GET(":id", handler.GetRisk)
			risks.POST("", handler.CreateRisk)
			risks.GET("", handler.ListRisks)
		}
	}

	return router
}
