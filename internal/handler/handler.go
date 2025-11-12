package handler

import (
	"Avito/internal/service"
	"github.com/gin-gonic/gin"
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
	return router
}
