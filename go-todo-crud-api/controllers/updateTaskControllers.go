package controllers

import (
	"go-todo-crud-api/initializers"
	"go-todo-crud-api/models"

	"github.com/gofiber/fiber/v2"
)

func UpdateTask(c *fiber.Ctx) error {
	c.Params("id")
	payload := struct {
		TaskHeader string `json:"TaskHeader"`
		TaskBody   string `json:"TaskBody"`
	}{}
	if err := c.BodyParser(&payload); err != nil {
		return err
	}
	var taskToBeUpdated models.TodoTask
	result := initializers.DB.First(&taskToBeUpdated, c.Params("id"))
	if result.RowsAffected == 1 {
		updateTask := models.TodoTask{TaskHeader: payload.TaskHeader, TaskBody: payload.TaskBody}
		initializers.DB.Model(&taskToBeUpdated).Updates(updateTask)
		return c.Status(200).JSON(fiber.Map{
			"updatedHeader": payload.TaskHeader,
			"updatedBody":   payload.TaskBody,
			"status":        "success",
			"message":       "Task Updated",
			"task_id":       c.Params("id"),
		})
	} else {
		return c.Status(200).JSON(fiber.Map{
			"updatedHeader": payload.TaskHeader,
			"updatedBody":   payload.TaskBody,
			"status":        "error",
			"message":       "No Task with the provided ID is available",
			"task_id":       c.Params("id"),
		})
	}

}
