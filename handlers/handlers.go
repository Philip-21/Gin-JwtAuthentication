package handlers

import (
	"context"
	"go-jwt/database"
	helper "go-jwt/helpers"
	"go-jwt/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

func HashPassword()

func Verifypassword()

func Signup() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user *models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		//Db operation
		validateErr := validate.Struct(&user)
		if validateErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validateErr.Error()})
			return
		}
		//countemail validates for email
		//checks the number of time a particular  email appears
		countemail, err := userCollection.CountDocuments(ctx, bson.H{"email": user.Email})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while checking for email"})
		}
		countphone, err := userCollection.CountDocuments(ctx, bson.H{"phone": user.Phone})

		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while checking for phone number"})
		}
		if countemail > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email already exists"})
		}
		if countphone > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "PHone number exists alredy"})
		}
	}
}

func Login()

func GetUsers()

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.Param("user_id") //( refers to the id /users/:user_id)

		//checks if the user is an admin or not based on the id
		if err := helper.MatchUserTypeToUid(c, userId); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		var user models.User
		//gettig a particular user from db
		//using decode to unmarshall the data for golang to understand
		err := userCollection.FindOne(ctx, bson.H{"user_id": userId}).Decode(&user)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}
