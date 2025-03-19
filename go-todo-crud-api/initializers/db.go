package initializers

import (
	"go-todo-crud-api/models"
	"log"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDatabase() {
	var err error
	mhost := os.Getenv("MYSQL_HOST")
	mport := os.Getenv("MYSQL_PORT")
	muser := os.Getenv("MYSQL_USER")
	mpass := os.Getenv("MYSQL_PASS")
	mdb := os.Getenv("MYSQL_DB")

	dsn := muser + ":" + mpass + "@tcp(" + mhost + ":" + mport + ")/" + mdb + "?charset=utf8mb4&parseTime=True&loc=Local"
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Could Not Connect to Database " + mdb + " at " + mhost + " Port:" + mport)
	}
}

func SyncDB() {
	DB.AutoMigrate(&models.TodoTask{})
	DB.Migrator().CreateIndex(&models.TodoTask{}, "id")
}
