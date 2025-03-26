package main

import (
	"ecommerce-cart/controllers"
	"ecommerce-cart/database"
	"ecommerce-cart/middleware"
	"ecommerce-cart/routes"

	"github.com/labstack/echo/v4"
)

func main() {
	var app *controllers.Application= controllers.NewApplication(database.ProductData(database.Client, "ProductCollection"), database.UserData(database.Client, "UserCollection"))

	e := echo.New()
	routes.RegisterRoutes(e)
	e.POST("/admin/addproduct", controllers.ProductViewAdmin())
	e.Use(middleware.Authentication())
	e.GET("/addtocart", app.AddToCart())
	e.GET("/removeitem", app.RemoveItem())
	e.GET("/cartcheckout", app.BuyFromCart())
	e.GET("/instantbuy", app.InstantBuy())

	e.Logger.Fatal(e.Start(":8000"))
}
