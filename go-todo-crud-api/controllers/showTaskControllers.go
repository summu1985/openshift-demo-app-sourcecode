package controllers

import (
	"go-todo-crud-api/initializers"
	"go-todo-crud-api/models"

	"github.com/gofiber/fiber/v2"
)

func ShowTask(c *fiber.Ctx) error {
	var alltasks []models.TodoTask
	result := initializers.DB.Find(&alltasks)

	if result.RowsAffected >= 1 && result.Error == nil {
		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "Task Fetched",
			"tasks":   alltasks,
		})
	} else {
		return c.Status(200).JSON(fiber.Map{
			"status":  "error",
			"message": "No Task Available",
		})
	}
}
