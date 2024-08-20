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
var client string = getEnvValue("CLIENTCOLL")
var database string = getEnvValue("DATABASE")
var projects string = getEnvValue("PROJECTCOLL")

func main() {
	/*
		Initialize the MongoDB client and check for errors.
		Setup the Gin router with necessary middleware and routes.
		Start the Gin server on localhost:8080.
	*/
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
	/*
		Load environment variables from a .env file.
		Return the value of the specified environment variable.
	*/
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv(key)
}

func setupMongodbClient() (*mongo.Client, error) {
	/*
		Create a MongoDB client with the URI obtained from environment variables.
		Set the server API version and return the client instance.
	*/
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	mongouri := getEnvValue("MONGOURI")
	opts := options.Client().ApplyURI(mongouri).SetServerAPIOptions(serverAPI)
	return mongo.Connect(context.Background(), opts)
}

func getUserValueHandler(c *gin.Context) {
	/*
		Extract the value type from the request parameters.
		Retrieve user value from the MongoDB collection and send it in the response.
	*/
	value := c.Param("valuetype")
	coll := mongoClient.Database(database).Collection(client)
	data := mongodbsetup.GetUserValue(coll, "pureheroky", value)
	c.JSON(http.StatusOK, gin.H{"data": data, "status": http.StatusOK})
}

func getImageHandler(c *gin.Context) {
	/*
		Extract collection name and query name from request parameters.
		Retrieve image data and content type from the MongoDB collection.
		Send image data in the response with appropriate content type.
	*/
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
	/*
		Extract project ID from request parameters.
		Retrieve project details from the MongoDB collection and send them in the response.
	*/
	prjId := c.Param("projectid")
	coll := mongoClient.Database(database).Collection(projects)
	data := mongodbsetup.GetProject(coll, prjId)
	c.JSON(http.StatusOK, gin.H{"data": data, "status": http.StatusOK})
}

func uploadImageHandler(c *gin.Context) {
	/*
		Extract parameters for image upload from the request.
		Save the image to the database using the specified parameters.
		Log fatal errors if any occur during the upload process.
	*/
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
	/*
		Bind the incoming JSON data to the Project struct.
		Insert the project data into the MongoDB collection.
		Send appropriate responses based on the success or failure of the insertion.
	*/
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
	/*
		Set CORS headers to allow requests from "https://pureheroky.com".
		Handle preflight OPTIONS requests with a 204 status code.
	*/
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
