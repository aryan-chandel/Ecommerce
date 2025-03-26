package middleware

import (
	
	"ecommerce-cart/tokens"
	"net/http"

	"github.com/labstack/echo/v4"
)
func Authentication() echo.MiddlewareFunc {
    return func(next echo.HandlerFunc) echo.HandlerFunc {
        return func(c echo.Context) error {
            clientToken := c.Request().Header.Get("token")
            if clientToken == "" {
                return c.JSON(http.StatusNotFound, map[string]string{"error": "token not found"})
            }

            claims, msg := tokens.ValidateToken(clientToken)
            if msg != "" {
                return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
            }

            c.Set("email", claims.Email)
            c.Set("uid", claims.Uid)

            return next(c) 
        }
    }
}
