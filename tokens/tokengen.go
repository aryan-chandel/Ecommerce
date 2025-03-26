package tokens

import (
	"context"
	"ecommerce-cart/database"
	"log"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SignedDetails struct { //fields used for token generator
	Email      string
	First_Name string
	Last_Name  string
	Uid        string
	jwt.StandardClaims
}

var SECRET_KEY = os.Getenv("SECRET_KEY")
var user = database.UserData(database.Client, "UserCollection")

func TokenGenerator(email string, firstname string, lastname string, uid string) (signedtoken string, refreshtoken string, err error) {
	claims := &SignedDetails{Email: email, First_Name: firstname, Last_Name: lastname, Uid: uid,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(24)).Unix(),
		},
	}
	refreshclaims := &SignedDetails{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(168)).Unix(),
		},
	}
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", "", err
	}
	refreshtoken, er := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshclaims).SignedString([]byte(SECRET_KEY))
	if er != nil {
		log.Panic(er)
		return
	}
	return token, refreshtoken, nil
}
func ValidateToken(signed string) (claims *SignedDetails, msg string) {
	token, err := jwt.ParseWithClaims(signed, &SignedDetails{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		log.Panic(err)
		msg = err.Error()
		return
	}
	claims, ok := token.Claims.(*SignedDetails)
	if !ok {
		msg = "the token is invalid"
		return
	}
	//check expiry of token
	if claims.ExpiresAt < time.Now().Local().Unix() {
		msg = "token is already expired"
		return
	}
	return claims, msg

}
func UpdateAllTokens(signedtoken string, signedrefresh string, uid string) {
	var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
	defer cancel()

	var updateobj primitive.D
	updateobj = append(updateobj, bson.E{Key: "token", Value: signedtoken})
	updateobj = append(updateobj, bson.E{Key: "refresh_token", Value: signedrefresh})
	updatedat, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	updateobj = append(updateobj, bson.E{Key: "updated_at", Value: updatedat})
	upsert := true

	filter := bson.D{{Key: "user_id", Value: uid}}
	update := bson.D{{Key: "$set", Value: updateobj}}
	opt := options.UpdateOptions{
		Upsert: &upsert,
	}
	_, err := user.UpdateOne(ctx, filter, update, &opt)
	if err != nil {
		log.Panic(err)
		return
	}

}
