package controllers

import (
	"go-todo-crud-api/initializers"
	"go-todo-crud-api/models"

	"github.com/gofiber/fiber/v2"
)

func DeleteTask(c *fiber.Ctx) error {
	deleteId := c.Params("id")
	result := initializers.DB.Delete(&models.TodoTask{}, deleteId)
	if result.RowsAffected < 1 {
		return c.Status(200).JSON(fiber.Map{
			"status":  "error",
			"message": "Task Not Available",
			"task_id": c.Params("id"),
		})
	} else {
		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "Task Deleted",
			"task_id": c.Params("id"),
		})
	}

}
