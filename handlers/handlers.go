package handlers

import (
	"context"
	"fmt"
	"go-jwt/database"
	helper "go-jwt/helpers"
	"go-jwt/models"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}
	return string(bytes)
}

func Verifypassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	var msg string
	if err != nil {
		msg = fmt.Sprintf("email or password is incorrect ")
		check = false
	}
	return check, msg
}

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

		password := HashPassword(user.Password)
		user.Password = password

		countphone, err := userCollection.CountDocuments(ctx, bson.H{"phone": user.Phone})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error while checking for phone number"})
			return
		}
		if countemail > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email already exists"})
			return
		}
		if countphone > 0 {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "PHone number exists alredy"})
			return
		}
		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.ID = primitive.NewObjectID()
		//the string to access the user
		user.User_id = user.ID.Hex() //returns the hex encoding of the ObjectID as a string

		//generating token sent to the user
		token, refreshToken, _ := helper.GenerateAllTokens(*&user.Email, *&user.First_name, *&user.Last_name, *&user.User_type, *&user.User_id)
		//set the user token
		user.Token = token
		user.Refresh_token = refreshToken

		//insert the user items into the database
		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user) //insert a single document into the collection one at a time()
		if insertErr != nil {
			msg := fmt.Sprintf("User item was not created ")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()
		c.JSON(http.StatusOK, resultInsertionNumber)
	}
}

func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		//a varialble declared when data exists in db
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		//db login action
		err := userCollection.FindOne(ctx, bson.H{"email": user.Email}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "email or password incorrect"})
			return
		}
		passwordIsValid, msg := Verifypassword(user.Password, foundUser.Password)
		defer cancel()
		if passwordIsValid != true {
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		//additional validations
		if &foundUser.Email == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "user not found"})
			return
		}
		token, refreshToken, err := helper.GenerateAllTokens(*&foundUser.Email, *&foundUser.First_name, *&foundUser.Last_name, *&foundUser.User_type, *&foundUser.User_id)
		if err != nil {
			return
		}
		helper.UpdateAllTokens(token, refreshToken, foundUser.User_id)
		userCollection.FindOne(ctx, bson.H{"user_id": foundUser.User_id}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, foundUser)
	}
}

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
