package database

import (
	"context"
	"fmt"
	"log"
	"time"
	"os"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var uri= os.Getenv("MONGO_URI")
func DBset() *mongo.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Println("unable to connect to mongodb")
		return nil
	}
	fmt.Println("successfully connected")
	return client
}

var Client = DBset()

func UserData(client *mongo.Client, collectioname string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("Ecommmerce").Collection(collectioname)
	return collection

}
func ProductData(client *mongo.Client, collectioname string) *mongo.Collection {
	var productcollection = client.Database("Ecommerce").Collection(collectioname)
	return productcollection

}
