package routes

import (
	"go-api/controllers"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterRoutes(e *echo.Echo, client *mongo.Client) {
	controllers.SetClient(client)

	e.POST("/add", controllers.AddPoints)
	e.POST("/spend", controllers.SpendPoints)
	e.GET("/balance", controllers.GetBalance)
}
