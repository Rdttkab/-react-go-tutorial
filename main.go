package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type Todo struct {
	ID        int    `json:"id"`
	Completed bool   `json:"completed"`
	Body      string `json:"body"`
}

func main() {
	fmt.Println("Hello World")

	app := fiber.New()

	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	PORT := os.Getenv("PORT")

	if PORT == "" {
        PORT = "4000"
    }

	todos := []Todo{}

	// get all todo
	app.Get("/api/todo",func(c *fiber.Ctx) error {
		return c.Status(200).JSON(todos)
	})

	// get a todo
	app.Get("/api/todo/:id",func(c *fiber.Ctx) error {
		id:= c.Params("id")

		for i,todo := range(todos) {	
			if fmt.Sprint(todo.ID) == id {
				return c.Status(200).JSON(todos[i])
			}
		}

		return c.Status(404).JSON(fiber.Map{"error": "Todo is not found"})
	})

	//create a todo
	app.Post("/api/todo", func(c *fiber.Ctx) error {
		todo := &Todo{}

		err := c.BodyParser(todo)

		if err != nil {
			return c.Status(400).JSON(fiber.Map{"error":"Todo body is required"})
		}

		if todo.Body == "" {
			return c.Status(400).JSON(fiber.Map{"error":"Todo body is required"})
		}

		todo.ID = len(todos) + 1

		todos = append(todos, *todo)

		return c.Status(201).JSON(todo)
	})

	//update a todo
	app.Put("/api/todo/:id",func(c *fiber.Ctx) error {
		id:= c.Params("id")
		// todo := &Todo{}

		for i,todo := range(todos) {	
			if fmt.Sprint(todo.ID) == id {
				todos[i].Completed = true

				return c.Status(200).JSON(todos[i])
			}
		}

		return c.Status(404).JSON(fiber.Map{"error": "Todo is not found"})
	})

	app.Delete("/api/todo/:id",func(c *fiber.Ctx) error {
		id:= c.Params("id")

		for i,todo := range(todos) {	
			if fmt.Sprint(todo.ID) == id {
				todos = append(todos[:i], todos[i+1:]...)

				return c.Status(200).JSON(fiber.Map{"success":true})
			}
		}

		return c.Status(404).JSON(fiber.Map{"error": "Todo is not found"})
	})

	log.Fatal(app.Listen(":" + PORT))
}
