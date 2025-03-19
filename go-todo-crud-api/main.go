package main

import (
	"go-todo-crud-api/controllers"
	"go-todo-crud-api/initializers"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectToDatabase()
	initializers.SyncDB()
}
func main() {
	app := fiber.New()
	app.Use(logger.New())

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "Welcome to Golang, Fiber, and GORM",
		})
	})

	app.Post("/create", controllers.CreateTask)
	app.Get("/show", controllers.ShowTask)
	app.Post("/update/:id", controllers.UpdateTask)
	app.Delete("/delete/:id", controllers.DeleteTask)

	app.Listen(":" + os.Getenv("APP_PORT"))
}
