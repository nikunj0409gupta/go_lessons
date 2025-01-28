package handlers

import (
	"context"
	"fmt"
	"log"

	"net/http"
	"time"

	"crud-app/models"

	"github.com/labstack/echo/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var userCollection *mongo.Collection

func InitMongoDB(uri, dbName, collectionName string) {
	clientOptions := options.Client().ApplyURI(uri)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	userCollection = client.Database(dbName).Collection(collectionName)
	fmt.Println("Connected to MongoDB!")
}

// This bind method is doing the same work as that of req.body in node
// It is binding the request body to the user struct
// If we want to modify the data that can be done after binding
// Creating timeout context is important to avoid memory leaks - if request takes too long to respond
// Need to talk about things like to do 1000 requests at a time
func CreateUser(c echo.Context) error {
	user := new(models.User)
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	res, err := userCollection.InsertOne(ctx, user)
	//fmt.Println(c.Request().Body, user)

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
	}
	fmt.Println("User created successfully", res)
	return c.JSON(http.StatusCreated, map[string]interface{}{
		"message": "User created successfully",
		"res":     res,
	})
}

// This is simple
// create the context for preventing timeout
// find all the users
// We ave created cursor and can't return it directly
// We will have to store it in a data structure and then return it
func GetAllUsers(c echo.Context) error {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := userCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch users"})
	}
	defer cursor.Close(ctx)

	var users []models.User
	if err := cursor.All(ctx, &users); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to parse users"})
	}

	return c.JSON(http.StatusOK, users)
}

// This is also simple
// .We get params from contecxt through echo we dont have req, res here
// It's contect which has data and helps in returing as well
// Converting id to object Id is important else it will not read
// Methods related to cursor are passed with pointer so that they can update the values of varibles
func GetUserByID(c echo.Context) error {
	id := c.Param("id")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var user models.User
	err1 := userCollection.FindOne(ctx, bson.M{"_id": objectID}).Decode(&user)
	if err1 != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
	}

	return c.JSON(http.StatusOK, user)
}

// It is quite similar to create and get user request
// Same process to pass req body or body which is through bind.
// If you want updated user getting retuned then decode it using pointer
// convert id to object id
func UpdateUserByID(c echo.Context) error {
	id := c.Param("id")

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}
	user := new(models.User)
	if err := c.Bind(user); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err1 := userCollection.UpdateOne(ctx, bson.M{"_id": objectID}, bson.M{"$set": user})
	if err1 != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update user"})
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "User updated successfully"})
}

// Fair enough and quite clear as well to me
func DeleteUserByID(c echo.Context) error {
	id := c.Param("id")
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid ID format"})
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_, err1 := userCollection.DeleteOne(ctx, bson.M{"_id": objectID})
	if err1 != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete user"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "User deleted successfully"})
}
