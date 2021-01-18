package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/luschnat-ziegler/musicbox/app"
	"github.com/luschnat-ziegler/musicbox/logger"
	"os"
)

func init() {
	// loads values from .env into the system
	if err := godotenv.Load(); err != nil {
		logger.Error("No .env file found")
	}
}

func main() {
	fmt.Println(os.LookupEnv("MUSIC_PATH"))
	app.Start()
}
