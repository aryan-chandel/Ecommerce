package controllers

import (
	"context"
	"ecommerce-cart/database"
	"ecommerce-cart/models"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)
type Application struct{
	prodCollection *mongo.Collection
	userCollection *mongo.Collection
}
func NewApplication(prodcollection,usercollection *mongo.Collection) *Application{
	return &Application{
		prodCollection: prodcollection,
		userCollection: usercollection,
	}
}
//cart specific functions
func(app *Application) AddToCart() echo.HandlerFunc{
	return func(c echo.Context) error {
		ProductQueryid:= c.QueryParam("id")
		if ProductQueryid==""{
			return c.JSON(http.StatusBadRequest,map[string]string{"error":"productid is missing"})
		}
		userid:=c.QueryParam("userid")
		if userid==""{
			return c.JSON(http.StatusBadRequest,map[string]string{"error":"userid is missing"})
		}
		//check if ids are genuine
		productId,err := primitive.ObjectIDFromHex("ProductQueryid")
		if err!=nil{
			return c.JSON(http.StatusInternalServerError,errors.New("invalid id"))
		}
		var ctx , cancel= context.WithTimeout(context.Background(),5*time.Second)
		defer cancel()

		er:= database.AddItemToCart(ctx, app.prodCollection,app.userCollection,productId, userid)
		if er!=nil{
			return c.JSON(http.StatusInternalServerError,map[string]string{"error":er.Error()})
		}
		return c.JSONPretty(http.StatusOK, map[string]string{"message":"successfully added"}," ")
	}

}

func(app *Application) RemoveItem() echo.HandlerFunc{
	return func(c echo.Context) error {
		ProductQueryid:= c.QueryParam("id")
		if ProductQueryid==""{
			return c.JSON(http.StatusBadRequest,map[string]string{"error":"productid is missing"})
		}
		userid:=c.QueryParam("userid")
		if userid==""{
			return c.JSON(http.StatusBadRequest,map[string]string{"error":"userid is missing"})
		}
		//check if ids are genuine
		productId,err := primitive.ObjectIDFromHex("ProductQueryid")
		if err!=nil{
			return c.JSON(http.StatusInternalServerError,errors.New("invalid id"))
		}
		var ctx , cancel= context.WithTimeout(context.Background(),5*time.Second)
		defer cancel()
		er:= database.RemoveItemFromCart(ctx, productId,userid,app.prodCollection,app.userCollection)
		if er!=nil{
			return c.JSON(http.StatusInternalServerError,map[string]string{"error":er.Error()})
		}
		return c.JSON(http.StatusOK,map[string]string{"message":"successfully removed"})
	}	

}

func(app *Application) InstantBuy() echo.HandlerFunc{
	return func(c echo.Context) error {
		ProductQueryid:= c.QueryParam("id")
		if ProductQueryid==""{
			return c.JSON(http.StatusBadRequest,map[string]string{"error":"productid is missing"})
		}
		userid:=c.QueryParam("userid")
		if userid==""{
			return c.JSON(http.StatusBadRequest,map[string]string{"error":"userid is missing"})
		}
		//check if ids are genuine
		productId,err := primitive.ObjectIDFromHex("ProductQueryid")
		if err!=nil{
			return c.JSON(http.StatusInternalServerError,errors.New("invalid id"))
		}
		var ctx , cancel= context.WithTimeout(context.Background(),5*time.Second)
		defer cancel()
		er:= database.InstantBuyer(ctx, userid,productId,app.prodCollection,app.userCollection)
		if er!=nil{
			return c.JSON(http.StatusInternalServerError,map[string]string{"error":er.Error()})
		}
		return c.JSON(http.StatusOK,map[string]string{"message":"order successfully placed"})
		
	}	
}

func(app*Application) GetItemFromCart() echo.HandlerFunc{
	return func(c echo.Context) error {
		userid:= c.QueryParam("id")
		if userid==""{
			return c.JSON(http.StatusNotFound,map[string]string{"error":"id not found"})
		}
		userst_id,err:= primitive.ObjectIDFromHex(userid)
		if err!=nil{
			log.Println(err)
			return c.String(http.StatusInternalServerError,"invalid")
		}
		var ctx,cancel=context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()
		var cartholder models.User
		err = UserCollection.FindOne(ctx, bson.D{primitive.E{Key:"_id",Value: userst_id}}).Decode(&cartholder)
		if err!=nil{
			log.Panic(err)
			return c.JSONPretty(500,map[string]string{"error":"something went wrong"}," ")
		}
		//aggregation
		filter_match:= bson.D{{Key:"$match",Value: bson.D{{Key: "_id",Value: userst_id}}}}
		unwind:=bson.D{{Key: "$unwind",Value: bson.D{{Key: "path",Value:"$usercart" }}}}
		grouping:= bson.D{{Key: "$group",Value: bson.D{{Key: "_id",Value: "$_id"},{Key: "total",Value: bson.D{{Key: "$sum",Value: "$usercart.price"}}}}}}
		cursor,er:=UserCollection.Aggregate(ctx,  mongo.Pipeline{filter_match,unwind, grouping})
		if er!=nil{
			log.Panic(er)
			return c.String(500,"something went wrong")
		}
		var list []bson.M
		if err= cursor.All(ctx,&list);err!=nil{
			log.Panic(err)
			return c.JSON(500,map[string]string{"error":"something went wrong"})
		}
		var totalPrice interface{} = nil
		if len(list) > 0 {
		totalPrice = list[0]["total"] // Extract the total price if available
		}

		return c.JSON(200, map[string]interface{}{
		"usercart": cartholder.UserCart,
		"total":    totalPrice,
		})

	}
}


func(app *Application) BuyFromCart() echo.HandlerFunc{
	return func(c echo.Context) error {
		userid := c.QueryParam("userid")
		if userid==""{
			return c.JSON(http.StatusBadRequest,map[string]string{"error":"empty user id"})
		}
		var ctx, cancel =context.WithTimeout(context.Background(),5*time.Second)
		defer cancel()
		er:= database.BuyItemFromCart(ctx, userid, app.userCollection,app.prodCollection)
		if er!=nil{
			return c.JSON(http.StatusInternalServerError,map[string]string{"error":er.Error()})
		}
		return c.JSONPretty(200,map[string]string{"..":"order placed"}," ")
	}
} 