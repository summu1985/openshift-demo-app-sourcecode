package controllers

import (
	"go-todo-crud-api/initializers"
	"go-todo-crud-api/models"

	"github.com/gofiber/fiber/v2"
)

func CreateTask(c *fiber.Ctx) error {

	payload := struct {
		TaskHeader string `json:"TaskHeader"`
		TaskBody   string `json:"TaskBody"`
	}{}

	if err := c.BodyParser(&payload); err != nil {
		return err
	}
	addtask := models.TodoTask{TaskHeader: payload.TaskHeader, TaskBody: payload.TaskBody}
	initializers.DB.Create(&addtask)

	return c.Status(200).JSON(fiber.Map{
		"header":  payload.TaskHeader,
		"body":    payload.TaskBody,
		"status":  "success",
		"message": "Task Created",
		"task_id": addtask.ID,
	})

}
