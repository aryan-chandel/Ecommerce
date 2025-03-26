package database

import (
	"context"
	"ecommerce-cart/models"
	"errors"
	"log"
	"time"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)
var (
	ErrCantFindProduct= errors.New("cant find the product")
	ErrcantDecodeProduct=errors.New("can't find the product")
	ErrUserIdIsNotValid=errors.New("this user is not valid")
	ErrCantupdateUser=errors.New("user cant be updated")
	ErrCantRemoveItemCart=errors.New("this item cannot be removed from cart")
	ErrCantGetitem=errors.New("unable to get item from the cart")
	ErrCantbuyCartItem=errors.New("cant update the purchase")
	ErrCantBuyItem=errors.New("cant place the order")
)
func AddItemToCart(c context.Context, prodcollection *mongo.Collection, usercollection *mongo.Collection,prodid primitive.ObjectID,userid string)error  {
	
	
	var item models.ProductUser
	
	err:= prodcollection.FindOne(c, bson.M{"_id":prodid}).Decode(&item)
	if err!=nil{
		log.Println(err)
		return ErrcantDecodeProduct
	}
	uid,msg:=primitive.ObjectIDFromHex(userid)
	if msg!=nil{
		log.Println(msg)
		return ErrUserIdIsNotValid
	}
	filter:= bson.D{{Key: "_id",Value: uid}}
	update:= bson.D{{Key: "$push",Value: bson.D{{Key: "usercart",Value: item}}}}
	
	result,er:= usercollection.UpdateOne(c,filter,update)
	if er!=nil{
		log.Println(er)
		return ErrCantupdateUser
	}
	if result.ModifiedCount==0{
		log.Print("cant update user")
		return ErrCantupdateUser
	}
	return nil 
}
func RemoveItemFromCart(c context.Context, prodid primitive.ObjectID,userid string, prodcollection *mongo.Collection, usercollection *mongo.Collection)error{
	id,er:= primitive.ObjectIDFromHex(userid)
	if er!=nil{
		log.Println(er)
		return ErrUserIdIsNotValid
	}
	filter:= bson.M{"_id":id}
	update:= bson.M{"$pull":bson.M{"usercart":bson.M{"_id":prodid}}}
	_,err:=usercollection.UpdateOne(c,filter,update)
	if err!=nil{
		log.Println(err)
		return ErrCantRemoveItemCart
	}

	return nil

}
func BuyItemFromCart(c context.Context,userid string, usercoll *mongo.Collection, prodcoll *mongo.Collection)error{
	id,er:= primitive.ObjectIDFromHex(userid)
	if er!=nil{
		log.Println(er)
		return ErrCantbuyCartItem
	}
	var cartholder models.User
	err:=usercoll.FindOne(c,bson.M{"_id":id}).Decode(&cartholder)
	if err!=nil{
		log.Println(err)
		return ErrUserIdIsNotValid
	}
	var GetCartItem []bson.M
	var new_order models.Order
	new_order.Order_Id= primitive.NewObjectID()
	new_order.Ordered_At=time.Now()
	new_order.Payment_Method.Cod=true
	new_order.Order_Cart=make([]models.ProductUser, 0)
	match:= bson.D{{Key: "$match",Value: bson.D{{Key: "_id",Value: id}}}}
	unwind:= bson.D{{Key: "$unwind",Value: bson.D{{Key: "path",Value: "$usercart"}}}}
	grouping:= bson.D{{Key: "$group",Value: bson.D{{Key:"_id",Value: "$_id" },{Key: "total",Value: bson.D{{Key: "$sum",Value: "$usercart.price"}}}}}}
	result,err:=usercoll.Aggregate(c,mongo.Pipeline{match,unwind,grouping})
	if err!=nil{
		log.Println(err)
		return ErrCantGetitem
	}
	err= result.All(c,&GetCartItem)
	if err!=nil{
		log.Println(err)
		return ErrcantDecodeProduct
	}
	var total_price int32
	for _,items:=range GetCartItem{
		price:= items["total"]
		total_price= price.(int32) 
	}
	new_order.Price=int(total_price)
	new_order.Order_Cart=cartholder.UserCart
	emptycart :=make([]models.ProductUser, 0)
	filter:= bson.D{{Key: "_id",Value: id}}
	update:= bson.D{{Key: "$push",Value: bson.D{{Key:"orders",Value: new_order}}}}
	results,err:=usercoll.UpdateOne(c,filter,update)
	if err!=nil{
		log.Println(err)
		return ErrCantupdateUser
	}
	if results.ModifiedCount==0{
		return ErrCantupdateUser
	}
	update1:=bson.D{{Key: "$set",Value:bson.D{{Key: "usercart",Value: emptycart}}}}
	results,err=usercoll.UpdateOne(c,filter,update1)
	if err!=nil{
		log.Panicln(err)
		return ErrCantbuyCartItem
	}
	if results.ModifiedCount==0{
		return ErrCantbuyCartItem
	}

	return nil

}
func InstantBuyer(c context.Context,userid string,prodid primitive.ObjectID,prodcoll *mongo.Collection,usercoll *mongo.Collection)error{
	//fetch product,,.. create order then update user orderstatus
	id,er:= primitive.ObjectIDFromHex(userid)
	if er!=nil{
		log.Println(er)
		return ErrUserIdIsNotValid
	}
	var product_detail models.ProductUser
	var order_detail models.Order
	order_detail.Order_Id=primitive.NewObjectID()
	order_detail.Ordered_At=time.Now()
	order_detail.Order_Cart= make([]models.ProductUser, 0)
	order_detail.Payment_Method.Cod=true

	filter:= bson.M{"_id":prodid}
	
	er=prodcoll.FindOne(c,filter).Decode(&product_detail)
	if er!=nil{
		log.Println(er)
		return ErrCantFindProduct
	}
	order_detail.Price=product_detail.Price
	order_detail.Order_Cart=append(order_detail.Order_Cart, product_detail)
	filter1 := bson.M{"_id":id}
	update1 := bson.M{"$push":bson.M{"orders":order_detail}}
	_,err:= usercoll.UpdateOne(c,filter1,update1)
	if err!=nil{
		log.Println(err)
		return ErrCantBuyItem
	}
	return nil

}