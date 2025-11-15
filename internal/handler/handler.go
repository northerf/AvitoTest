package handler

import (
	"github.com/gin-gonic/gin"

	"Avito/internal/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.Default()
	teams := router.Group("/team")
	{
		teams.POST("/add", h.CreateTeam)
		teams.GET("/get", h.GetTeam)
		teams.POST("/deactivate", h.DeactivateUsers)
	}
	users := router.Group("/users")
	{
		users.POST("/setIsActive", h.SetUserActive)
		users.GET("/getReview", h.GetUserReviews)
	}
	pr := router.Group("/pullRequest")
	{
		pr.POST("/create", h.CreatePR)
		pr.POST("/merge", h.MergePR)
		pr.POST("/reassign", h.ReassignReviewer)
	}
	stats := router.Group("/stats")
	{
		stats.GET("/allstats", h.GetStatistics)
		stats.GET("/user", h.GetUserStatistics)
	}

	router.GET("/health", h.Health)
	return router
}
