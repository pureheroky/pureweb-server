package handlers

import (
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
)

type HandlerInterface interface {
	GetUserValueHandler(c *gin.Context)
	GetImageHandler(c *gin.Context)
	GetProjectHandler(c *gin.Context)
	UploadImageHandler(c *gin.Context)
	CreateProjectHandler(c *gin.Context)
}

type Handler struct {
	MongoClient *mongo.Client
}

func NewHandler(mongoClient *mongo.Client) *Handler {
	return &Handler{
		MongoClient: mongoClient,
	}
}
