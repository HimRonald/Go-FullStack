package main

import (
	"context"
	"fmt"
	"os"

	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var collection *mongo.Collection

type Todo struct {
	ID        int    `json:"id" bson:"id"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

func main() {
	fmt.Println("Hello, GO! Ronald")

	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Error loading .env file", err)
		return
	}

	MONGODB_URL := os.Getenv("MONGODB_URL")
	clientOptions := options.Client().ApplyURI(MONGODB_URL)
	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB ATLAS")

	collection = client.Database("goland_db").Collection("todos")

	app := fiber.New()

	app.Get("/api/todos", getTodos)
	app.Post("/api/todos", createTodos)
	app.Patch("/api/todos/:id", updateTodos)
	app.Delete("/api/todos/:id", deleteTodos)

	port := os.Getenv("PORT")
	if port == "" {
		port = "5555"
	}

	log.Fatal(app.Listen("0.0.0.0:" + port))
}

// func getTodos(c *fiber.Ctx) error {
// 	var todos []Todo

// 	collection.Find(context.Background(),bson.M{})
// }
// func createTodos(c *fiber.Ctx) error {}
// func updateTodos(c *fiber.Ctx) error {}
// func deleteTodos(c *fiber.Ctx) error {}
