package main

import (
	"ecommerce-cart/controllers"
	"ecommerce-cart/database"
	"ecommerce-cart/middleware"
	"ecommerce-cart/routes"
	"fmt"
	"os"

	"github.com/labstack/echo/v4"
)

func main() {
	var app *controllers.Application = controllers.NewApplication(database.ProductData(database.Client, "ProductCollection"), database.UserData(database.Client, "UserCollection"))
	e := echo.New()

	fmt.Println("secret key:", os.Getenv("SECRET_KEY"))
	routes.RegisterRoutes(e)
	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(200, map[string]string{"message": "status working"})
	})
	usergroup := e.Group("/users")
	usergroup.Use(middleware.Authentication())
	e.POST("/admin/addproduct", controllers.ProductViewAdmin())
	usergroup.GET("/addtocart", app.AddToCart())
	usergroup.GET("/removeitem", app.RemoveItem())
	usergroup.GET("/cartcheckout", app.BuyFromCart())
	usergroup.GET("/instantbuy", app.InstantBuy())

	e.Logger.Fatal(e.Start(":8000"))
}
