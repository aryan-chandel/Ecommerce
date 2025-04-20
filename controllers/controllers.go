package controllers

import (
	"context"
	"ecommerce-cart/database"
	"ecommerce-cart/models"
	"ecommerce-cart/tokens"
	"log"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var validate = validator.New()
var UserCollection *mongo.Collection = database.UserData(database.Client, "UserCollection")
var ProductCollection *mongo.Collection = database.ProductData(database.Client, "ProductCollection")

func HashPassword(password string) string {
	pass, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(pass)
}

func Verifypassword(userpassword string, givenpass string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(givenpass), []byte(userpassword))
	valid := true
	msg := ""
	if err != nil {
		valid = false
		msg = "login or password is incorrect"
	}
	return valid, msg
}

func Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		var founduser models.User
		err := c.Bind(&user)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})

		}
		err = UserCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&founduser)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "user cant be found"})
		}
		passIsValid, msg := Verifypassword(*user.Password, *founduser.Password)
		if !passIsValid {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": msg})
		}
		token, refreshtoken, err := tokens.TokenGenerator(*founduser.Email, *founduser.First_Name, *founduser.Last_Name, founduser.User_Id)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "something went wrong try again later"})
		}
		tokens.UpdateAllTokens(token, refreshtoken, founduser.User_Id)
		return c.JSON(http.StatusFound, map[string]string{"message": "login successfully"})

	}

}

func SignUp() echo.HandlerFunc {
	return func(c echo.Context) error {

		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second) // for db queries
		defer cancel()

		var user models.User
		err := c.Bind(&user)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{"error": validationErr})
		}

		count, err := UserCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		if count > 0 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "user already exists"})
		}

		count, err = UserCollection.CountDocuments(ctx, bson.M{"phone": user.Phone})
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
		}
		if count > 0 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "user already exists"})
		}

		password := HashPassword(*user.Password)
		user.Password = &password
		user.Created_At = time.Now()
		user.Updated_At = time.Now()
		user.ID = primitive.NewObjectID()
		user.User_Id = user.ID.Hex()

		token, refreshToken, _ := tokens.TokenGenerator(*user.Email, *user.First_Name, *user.Last_Name, user.User_Id) //jwt
		user.Token = &token
		user.Refresh_Token = &refreshToken //check
		user.UserCart = []models.ProductUser{}
		user.Adress_Details = []models.Address{}
		user.Order_Status = []models.Order{}

		_, insertErr := UserCollection.InsertOne(ctx, user)
		if insertErr != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": insertErr.Error()})
		}

		return c.JSONPretty(http.StatusCreated, map[string]string{"message": "Successfully signed up"}, "  ")
	}
}

func SearchProduct() echo.HandlerFunc {
	return func(c echo.Context) error {
		var Productlist []models.Product
		var ctx, cancel = context.WithTimeout(context.Background(), 50*time.Second)
		defer cancel()

		cursor, err := ProductCollection.Find(ctx, bson.D{{}})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "something went wrong"})
		}
		err = cursor.All(ctx, &Productlist)
		if err != nil {
			log.Println(err)
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "something went wrong"})
		}
		defer cursor.Close(ctx)
		if err = cursor.Err(); err != nil {
			log.Println(err)
			return c.String(http.StatusInternalServerError, "invalid")

		}
		return c.JSON(http.StatusOK, Productlist)

	}
}

func SearchProductByQuery() echo.HandlerFunc {
	return func(c echo.Context) error {
		var SearchProduct []models.Product
		query := c.QueryParam("name")
		if query == "" {
			log.Println("query is empty")
			return c.JSON(http.StatusNotFound, map[string]string{">>": "invalid search index"})

		}
		var ctx, cancel = context.WithTimeout(context.Background(), 50*time.Second)
		defer cancel()
		cursor, err := ProductCollection.Find(ctx, bson.M{"product_name": bson.M{"$regex": query, "$options": "i"}})
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{"message": "something went wrong"})

		}
		err = cursor.All(ctx, SearchProduct)
		if err != nil {
			log.Print("err")
			return c.String(400, "invalid")
		}
		defer cursor.Close(ctx)
		if err = cursor.Err(); err != nil {
			log.Print(err)
			return c.String(400, "invalid request")
		}
		return c.JSON(200, SearchProduct)
	}

}

func ProductViewAdmin() echo.HandlerFunc {
	return func(c echo.Context) error {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var product models.Product
		err := c.Bind(&product)
		if err != nil {
			log.Panic(err)
			return c.JSON(http.StatusBadRequest, "invalid")
		}
		product.Product_ID = primitive.NewObjectID()
		_, anyerr := ProductCollection.InsertOne(ctx, product)
		if anyerr != nil {
			log.Panic(anyerr)
			return c.JSON(http.StatusInternalServerError, "product cant be inserted")
		}
		return c.JSON(http.StatusOK, "successfullly added")

	}

}
