package main

import (
	"log"
	"sync"

	"pendekin/utils"

	"github.com/gin-gonic/gin"
	_ "github.com/joho/godotenv/autoload"
)

func startServer(wg *sync.WaitGroup) {
	defer wg.Done()

	server := gin.New()

	db, err := utils.InitDB()
	if err != nil {
		log.Fatal("Failed to init database!")
	}

	defer utils.DisconnectDB(db)

	rdb, err := utils.InitRedis()
	if err != nil {
		log.Fatal("Failed to init redis!")
	}

	server = utils.SetUpRoutes(server, db, rdb)

	if err := server.Run(":8080"); err != nil {
		log.Fatal("Failed to run server!")
	}
}

func main() {
	wg := sync.WaitGroup{}

	wg.Add(1)
	go startServer(&wg)

	wg.Wait()
}
