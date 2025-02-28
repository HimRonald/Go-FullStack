package main

// Importing the necessary packages
import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

var collection *mongo.Collection

// Structure the Todo collection for MongoDB in ATLAS
type Todo struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"` 
	Completed bool               `json:"completed"`
	Body      string             `json:"body"`
}

func main() {
	fmt.Println("Hello, GO! Ronald") // Testing the server with a my own message

	err := godotenv.Load(".env") // load the .env file that contains the MongoDB URL and PORT
	if err != nil {
		fmt.Println("Error loading .env file", err)
		return
	}

	MONGODB_URL := os.Getenv("MONGODB_URL") // Get the MongoDB URL from the .env file
	clientOptions := options.Client().ApplyURI(MONGODB_URL) 
	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(context.Background()) // disconnect connection after functions execution

	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MongoDB ATLAS") // Testing the connection to MongoDB ATLAS with given message

	collection = client.Database("goland_db").Collection("todos") // Create a collection in DB

	app := fiber.New()

	app.Get("/api/todos", getTodos)
	app.Post("/api/todos", createTodos)
	app.Patch("/api/todos/:id", updateTodos)
	app.Delete("/api/todos/:id", deleteTodos)

	port := os.Getenv("PORT") // false error handling for port 
	if port == "" {
		port = "5555"
	}

	log.Fatal(app.Listen("0.0.0.0:" + port))
}

// create Functions for get all the todos 
func getTodos(c *fiber.Ctx) error {
	var todos []Todo // declare todos variable in array

	cursor, err := collection.Find(context.Background(), bson.M{})

	if err != nil {
		return err
	}

	defer cursor.Close(context.Background()) // postponed the execution of the function until the surrounding function returns

	for cursor.Next(context.Background()) {
		var todo Todo
		if err := cursor.Decode(&todo); err != nil {
			return err
		}
		todos = append(todos, todo)
	}

	return c.JSON(todos)
}

// create function to create todos 
func createTodos(c *fiber.Ctx) error {
	todo := new(Todo)

	if err := c.BodyParser(todo); err != nil {
		return err
	}

	if todo.Body == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Todo body cannot be empty"})
	}

	insertResult, err := collection.InsertOne(context.Background(), todo)
	if err != nil {
		return err
	}

	todo.ID = insertResult.InsertedID.(primitive.ObjectID)

	return c.Status(201).JSON(todo) // 201 Created
}

// create function to update todos
func updateTodos(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id) // convert to primitive.ObjectID

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid todo ID"})
	}

	filter := bson.M{"_id": objectID} // filter the object ID to update
	update := bson.M{"$set": bson.M{"completed": true}}

	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{"message": true})
}

// create function to delete todos
func deleteTodos(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid todo ID"})
	}

	filter := bson.M{"_id": objectID}

	_, err = collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(fiber.Map{"message": true})
}
