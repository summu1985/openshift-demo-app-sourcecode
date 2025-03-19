package models

type TodoTask struct {
	ID         uint `gorm:"primaryKey"`
	TaskHeader string
	TaskBody   string
}
