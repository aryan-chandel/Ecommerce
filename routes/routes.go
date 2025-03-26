package routes

import (
	"ecommerce-cart/controllers"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo) {
	e.POST("/users/signup", controllers.SignUp())
	e.POST("/users/login", controllers.Login())
	e.GET("/users/productview", controllers.SearchProduct())
	e.GET("/users/search", controllers.SearchProductByQuery())
}
