package controllers

import (
	"context"
	"ecommerce-cart/models"
	"log"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func AddAddress() echo.HandlerFunc {
	return func(c echo.Context) error {
		user_id:=c.QueryParam("id")
		if user_id==""{
			return c.String(http.StatusNotFound,"id not found")
		}
		usert_id,err:= primitive.ObjectIDFromHex(user_id)
		if err!=nil{
			return c.String(http.StatusInternalServerError,"invalid ")
		}
		var addresses models.Address
		addresses.Address_id=primitive.NewObjectID()
		err=c.Bind(&addresses)
		if err!=nil{
			return c.JSON(http.StatusNotAcceptable,err)
		}
		var ctx,cancel=context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()
		//aggregate{filter,unwind,group and func}
		filter:= bson.D{{Key: "$match",Value: bson.D{{Key: "_id",Value: usert_id}}}}
		unwind:= bson.D{{Key: "$unwind",Value: bson.D{{Key: "path",Value: "$address"}}}}
		group:= bson.D{{Key: "$group",Value: bson.D{{Key: "_id",Value: "$address_id"},{Key: "count",Value: bson.D{{Key: "$sum",Value:1 }}}}}}
		cursor,er:= UserCollection.Aggregate(ctx, mongo.Pipeline{filter,unwind,group})
		if er!=nil{
			return c.JSONPretty(500,"internal server error"," ")

		}
		var addressinfo []bson.M
		err=cursor.All(ctx,&addressinfo)
		if err!=nil{
			log.Panic(err)
		}
		var size int32
		for _,address_no :=range addressinfo{
			count :=address_no["count"]
			size =count.(int32)

		}
		if size<2{
			//update user doc
			filter:=bson.D{{Key: "_id",Value: usert_id}}
			update:=bson.D{{Key: "$push",Value: bson.D{{Key: "address",Value: addresses}}}}
			_,err= UserCollection.UpdateOne(ctx,filter,update)
			if err!=nil{
				return c.JSON(500, "internal server error")
			}
			return c.JSONPretty(200,"successfully added"," ")

		}else{
			return c.JSON(400,"not allowed")
		}
		
	}
}

func EditHomeAddress() echo.HandlerFunc {
	return func(c echo.Context) error {
		user_id:=c.QueryParam("id")
		if user_id==""{
			return c.String(http.StatusNotFound,"id not found")
		}
		usert_id,err:= primitive.ObjectIDFromHex(user_id)
		if err!=nil{
			return c.String(http.StatusInternalServerError,"invalid ")
		}
		var editaddress models.Address
		err=c.Bind(&editaddress)
		if err!=nil{
			return c.JSON(http.StatusNotAcceptable,err)
		}
		var ctx,cancel=context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()

		filter:=bson.D{{Key: "$match",Value: bson.D{{Key: "_id",Value: usert_id}}}}
		update :=bson.D{{Key: "$set",Value: bson.D{{Key: "address.0.house",Value: editaddress.House},{Key: "address.0.street",Value: editaddress.Street},{Key: "address.0.city",Value: editaddress.City},{Key: "address.0.pincode",Value: editaddress.Pincode}}}}
		_,err=UserCollection.UpdateOne(ctx,filter,update)
		if err!=nil{
			return c.JSONPretty(http.StatusInternalServerError,"something went wrong","")

		}
		return c.JSONPretty(http.StatusOK,"successfully updated"," ")
	
	}


}

func EditWorkAddress() echo.HandlerFunc {
	return func(c echo.Context) error {
		user_id:=c.QueryParam("id")
		if user_id==""{
			return c.String(http.StatusNotFound,"id not found")
		}
		usert_id,err:= primitive.ObjectIDFromHex(user_id)
		if err!=nil{
			return c.String(http.StatusInternalServerError,"invalid ")
		}
		var editaddress models.Address
		err=c.Bind(&editaddress)
		if err!=nil{
			return c.JSON(http.StatusNotAcceptable,err)
		}
		var ctx,cancel=context.WithTimeout(context.Background(),100*time.Second)
		defer cancel()

		filter:=bson.D{{Key: "$match",Value: bson.D{{Key: "_id",Value: usert_id}}}}
		update :=bson.D{{Key: "$set",Value: bson.D{{Key: "address.1.house",Value: editaddress.House},{Key: "address.1.street",Value: editaddress.Street},{Key: "address.1.city",Value: editaddress.City},{Key: "address.1.pincode",Value: editaddress.Pincode}}}}
		_,err=UserCollection.UpdateOne(ctx,filter,update)
		if err!=nil{
			return c.JSONPretty(http.StatusInternalServerError,"something went wrong","")

		}
		return c.JSON(200,"successfully updated")
	}
	
}

func DeleteAddress() echo.HandlerFunc {
	return func(c echo.Context) error {
		user_id:= c.QueryParam("id")
		if user_id==""{
			return c.JSON(http.StatusNotFound,map[string]string{"error":"id not found"})
		}
		addresses:= make([]models.Address,0)
		userstored_id,err:= primitive.ObjectIDFromHex(user_id)
		if err!=nil{
			return c.JSONPretty(http.StatusInternalServerError,map[string]string{">>":"cant find user"}," ")
		}
		var ctx ,cancel=context.WithTimeout(context.Background(),50*time.Second)
		defer cancel()
		filter:= bson.D{primitive.E{Key:"_id",Value:userstored_id}}
		update:= bson.M{"$set": bson.D{primitive.E{Key:"address",Value: addresses}}}
		_,err=UserCollection.UpdateOne(ctx,filter,update)
		if err!=nil{
			return c.JSONPretty(http.StatusInternalServerError,map[string]string{"error":"something went wrong"}," ")

		}
		return c.String(200,"successfully deleted")


	}
	
}
