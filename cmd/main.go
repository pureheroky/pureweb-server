package main

import (
	"context"
	"log"
	"pureheroky.com/mongodbsetup/models"
	"pureheroky.com/mongodbsetup/mongodbsetup"

	"github.com/gin-gonic/gin"
	"pureheroky.com/mongodbsetup/routes"
)

func main() {
	mongoClient, err := mongodbsetup.SetupMongodbClient()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer func() {
		if err := mongoClient.Disconnect(context.Background()); err != nil {
			log.Printf("Failed to disconnect from MongoDB: %v", err)
		}
	}()

	app := &models.App{
		MongoClient: mongoClient,
	}

	r := gin.Default()
	routes.SetupRoutes(r, app)

	if err := r.Run("localhost:8080"); err != nil {
		log.Fatalf("Error while starting server: %v", err)
	}
}
