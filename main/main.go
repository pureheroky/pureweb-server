package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"pureheroky.com/mongodbsetup"
)

var mongoClient *mongo.Client

func main() {
	var err error
	mongoClient, err = setupMongodbClient()
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}
	defer mongoClient.Disconnect(context.Background())

	r := gin.Default()
	r.Use(corsMiddleware())

	r.GET("/getuservalue/:valuetype", getUserValueHandler)
	r.GET("/getimage/:collname/:queryname", getImageHandler)
	r.GET("/getproject/:projectid", getProjectHandler)
	r.GET("/uploadimage/:query/:id/:imagename", uploadImageHandler)
	r.POST("/projects", createProjectHandler)

	r.Run("localhost:8080")
}

func getEnvValue(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv(key)
}

var client string = getEnvValue("CLIENTCOLL")
var database string = getEnvValue("DATABASE")
var projects string = getEnvValue("PROJECTSCOLL")

func setupMongodbClient() (*mongo.Client, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	mongouri := getEnvValue("MONGOURI")
	opts := options.Client().ApplyURI(mongouri).SetServerAPIOptions(serverAPI)
	return mongo.Connect(context.Background(), opts)
}

func getUserValueHandler(c *gin.Context) {
	value := c.Param("valuetype")
	coll := mongoClient.Database(database).Collection(client)
	data := mongodbsetup.GetUserValue(coll, "pureheroky", value)
	c.JSON(http.StatusOK, gin.H{"data": data, "status": http.StatusOK})
}

func getImageHandler(c *gin.Context) {
	coll := mongoClient.Database(database).Collection(c.Param("collname"))
	imageData, contentType, err := mongodbsetup.GetImage(coll, c.Param("queryname"))
	if err != nil {
		if err == mongo.ErrNoDocuments {
			c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve image"})
		}
		return
	}

	c.Header("Content-Type", contentType)
	c.Writer.Write(imageData)
}

func getProjectHandler(c *gin.Context) {
	prjId := c.Param("projectid")
	coll := mongoClient.Database(database).Collection(projects)
	data := mongodbsetup.GetProject(coll, prjId)
	c.JSON(http.StatusOK, gin.H{"data": data, "status": http.StatusOK})
}

func uploadImageHandler(c *gin.Context) {
	id := c.Param("id")
	query := c.Param("query")
	image := c.Param("imagename")
	coll := mongoClient.Database(database).Collection(query)
	err := mongodbsetup.SaveImageInDB(coll, id, "temp/"+image)
	if err != nil {
		log.Fatal(err)
	}
}

func createProjectHandler(c *gin.Context) {
	var project mongodbsetup.Project
	if err := c.BindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	coll := mongoClient.Database(database).Collection(projects)
	_, err := coll.InsertOne(context.Background(), project)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, project)
}

func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "https://pureheroky.com")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}
