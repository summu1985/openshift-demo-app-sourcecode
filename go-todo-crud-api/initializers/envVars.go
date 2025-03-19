package initializers

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvVariables() {
	config_path := os.Getenv("CONFIG_PATH")
	config_file := config_path + "/todo-api.properties"
	if _, err := os.Stat(config_file); err != nil {
		log.Fatal(config_file + " File does not exist\n")
	}
	fmt.Println(config_file)
	err := godotenv.Load(config_file)
	if err != nil {
		log.Fatal("Error Loading Environment File")
	}
}
