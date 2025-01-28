package main

import (
	"crud-app/handlers"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	mongoURI := "mongodb://localhost:27017"
	dbName := "testdb"
	collectionName := "users"

	handlers.InitMongoDB(mongoURI, dbName, collectionName)

	// Creating server
	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Decalre all routes
	e.GET("/users", handlers.GetAllUsers)
	e.POST("/users", handlers.CreateUser)
	e.GET("/users/:id", handlers.GetUserByID)
	e.PUT("/users/:id", handlers.UpdateUserByID)
	e.DELETE("/users/:id", handlers.DeleteUserByID)

	e.Logger.Fatal(e.Start(":8080"))
}
