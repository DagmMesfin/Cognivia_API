package routers

import (
	"cognivia-api/Delivery/controllers"
	"cognivia-api/infrastructure"

	"github.com/gin-gonic/gin"
)

func SetupRouter(
	authHandler *controllers.UserHandler,
	notebookHandler *controllers.NotebookHandler,
) *gin.Engine {
	router := gin.Default()

	userRoutes := router.Group("/api/v1/users")
	{
		userRoutes.POST("/register", authHandler.Register)
		userRoutes.POST("/login", authHandler.Login)
		userRoutes.GET("/:id", authHandler.GetUser)
		userRoutes.PUT("/:id", authHandler.UpdateUser)
		userRoutes.DELETE("/:id", authHandler.DeleteUser)
	}

	notebookRoutes := router.Group("/api/v1/notebooks")
	{
		// Protected routes - require JWT authentication
		notebookRoutes.Use(infrastructure.JWTAuth())
		notebookRoutes.POST("/", notebookHandler.CreateNotebook)
		notebookRoutes.GET("/:id", notebookHandler.GetNotebook)
		notebookRoutes.GET("/user", notebookHandler.GetNotebooksByUserID)
		notebookRoutes.PUT("/:id", notebookHandler.UpdateNotebook)
		notebookRoutes.DELETE("/:id", notebookHandler.DeleteNotebook)
		notebookRoutes.GET("/:id/snapnotes", notebookHandler.GetSnapnotes)
		notebookRoutes.GET("/:id/prep-pilot", notebookHandler.GetPrepPilot)
	}

	return router
}
