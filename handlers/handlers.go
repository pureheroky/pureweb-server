package handlers

import (
	"context"
	"errors"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"pureheroky.com/server/config"
	"pureheroky.com/server/models"
	"pureheroky.com/server/mongodbsetup"
)

/*
GetUserValueHandler extract the value type from the request parameters.
Retrieve user value from the MongoDB collection and send it in the response.
*/
func (h *Handler) GetUserValueHandler(c *gin.Context) {
	value := c.Param("valuetype")
	coll := h.MongoClient.Database(config.Database).Collection(config.Client)
	data, err := mongodbsetup.GetUserValue(coll, "pureheroky", value)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": data, "status": http.StatusOK})
}

/*
GetImageHandler extract collection name and query name from request parameters.
Retrieve image data and content type from the MongoDB collection.
Send image data in the response with appropriate content type.
*/
func (h *Handler) GetImageHandler(c *gin.Context) {
	coll := h.MongoClient.Database(config.Database).Collection(c.Param("collname"))
	imageData, contentType, err := mongodbsetup.GetImage(coll, c.Param("queryname"))
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve image"})
		}
		return
	}

	c.Header("Content-Type", contentType)
	_, err = c.Writer.Write(imageData)
	if err != nil {
		log.Printf("Failed to write image: %v", err)
	}
}

/*
GetProjectHandler extract project ID from request parameters.
Retrieve project details from the MongoDB collection and send them in the response.
*/
func (h *Handler) GetProjectHandler(c *gin.Context) {
	prjId := c.Param("projectid")
	coll := h.MongoClient.Database(config.Database).Collection(config.Projects)
	data, err := mongodbsetup.GetProject(coll, prjId)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, gin.H{"data": data, "status": http.StatusOK})
}

/*
UploadImageHandler extract parameters for image upload from the request.
Save the image to the database using the specified parameters.
Log fatal errors if any occur during the upload process.
*/
func (h *Handler) UploadImageHandler(c *gin.Context) {
	var req models.UploadImage
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	coll := h.MongoClient.Database(config.Database).Collection(req.Query)
	err := mongodbsetup.SaveImageInDB(coll, req.ID, "../images/"+req.ImageName)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(200, gin.H{"message": "Image uploaded successfully"})
}

/*
CreateProjectHandler bind the incoming JSON data to the Project struct.
Insert the project data into the MongoDB collection.
Send appropriate responses based on the success or failure of the insertion.
*/
func (h *Handler) CreateProjectHandler(c *gin.Context) {
	var project models.Project
	if err := c.BindJSON(&project); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	coll := h.MongoClient.Database(config.Database).Collection(config.Projects)
	_, err := coll.InsertOne(context.Background(), project)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, project)
}
