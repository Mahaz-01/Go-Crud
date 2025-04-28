package routes

import (
	"gin-crud/internal/handlers"
	"gin-crud/internal/middleware"
    "github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {
    // Unprotected routes for authentication
    router.POST("/register", handlers.RegisterUser)
    router.POST("/login", handlers.LoginUser)

    // Protected routes for items (require JWT authentication)
    items := router.Group("/items")
    items.Use(middleware.JWTMiddleware())
    {
        items.GET("", handlers.GetItems)
        items.GET(":id", handlers.GetItem)
        items.POST("", handlers.CreateItem)
        items.PUT(":id", handlers.UpdateItem)
        items.DELETE(":id", handlers.DeleteItem)
    }
}