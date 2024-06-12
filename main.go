package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Todo struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Completed bool               `json:"completed"`
	Body      string             `json:"body"`
}

var collection *mongo.Collection

func main() {
	if os.Getenv("ENV") != "production" {
		err := godotenv.Load(".env")

		if err != nil {
			log.Fatal("Error loading .env file")
		}
	}

	MONGODB_URI := os.Getenv("MONGODB_URI")

	clientOptions := options.Client().ApplyURI(MONGODB_URI)

	client, err := mongo.Connect(context.Background(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(context.Background())

	err = client.Ping(context.Background(), nil)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Connected to MONGODB ATLAS")

	collection = client.Database("golang_db").Collection("todos")

	app := fiber.New()

	// app.Use(cors.New(cors.Config{
	// 	AllowOrigins: "http://localhost:5173",
	// 	AllowHeaders: "Origin,Content-Type,Accept",
	// }))

	app.Get("/api/todo", getAllTodo)
	app.Get("/api/todo/:id", getTodo)
	app.Post("/api/todo", createTodo)
	app.Put("/api/todo/:id", updateTodo)
	app.Delete("/api/todo/:id", deleteTodo)

	PORT := os.Getenv("PORT")

	if PORT == "" {
		PORT = "4000"
	}

	if os.Getenv("ENV") == "production" {
		app.Static("/", "./client/dist")
	}

	log.Fatal(app.Listen(":" + PORT))
}

// get all todo
func getAllTodo(c *fiber.Ctx) error {
	var todos []Todo

	cursor, err := collection.Find(context.Background(), bson.M{})

	if err != nil {
		return err
	}

	defer cursor.Close(context.Background())

	for cursor.Next(context.Background()) {
		var todo Todo

		err := cursor.Decode(&todo)

		if err != nil {
			return err
		}

		todos = append(todos, todo)
	}

	return c.Status(200).JSON(todos)
}

// // get a todo
func getTodo(c *fiber.Ctx) error {
	id := c.Params("id")

	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}

	var todo Todo

	err = collection.FindOne(context.Background(), bson.M{"_id": objectID}).Decode(&todo)

	if err != nil {
		// return err
		return c.Status(404).JSON(fiber.Map{"error": "Todo is not found"})
	}

	return c.Status(200).JSON(todo)
}

// create a todo
func createTodo(c *fiber.Ctx) error {
	todo := new(Todo)

	err := c.BodyParser(todo)

	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Todo body is required"})
	}

	if todo.Body == "" {
		return c.Status(400).JSON(fiber.Map{"error": "Todo body is required"})
	}

	result, err := collection.InsertOne(context.Background(), todo)

	if err != nil {
		return err
	}

	todo.ID = result.InsertedID.(primitive.ObjectID)

	fmt.Println("Todo is: ", *todo)

	return c.Status(201).JSON(todo)
}

// update a todo
func updateTodo(c *fiber.Ctx) error {
	id := c.Params("id")

	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}

	update := bson.M{"$set": bson.M{"completed": true}}

	var todo Todo

	err = collection.FindOneAndUpdate(context.Background(), bson.M{"_id": objectID}, update).Decode(&todo)

	if err != nil {
		return err
	}

	return c.Status(404).JSON(todo)
}

func deleteTodo(c *fiber.Ctx) error {
	id := c.Params("id")

	objectID, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		return err
	}

	var todo Todo
	err = collection.FindOneAndDelete(context.Background(), bson.M{"_id": objectID}).Decode(&todo)

	if err != nil {
		return err
	}

	return c.Status(200).JSON(todo)
}
