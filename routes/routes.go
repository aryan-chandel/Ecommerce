package routes

import (
	"ecommerce-cart/controllers"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo) {
	e.POST("/signup", controllers.SignUp())
	e.POST("/login", controllers.Login())
	e.GET("/productview", controllers.SearchProduct())
	e.GET("/search", controllers.SearchProductByQuery())
}
