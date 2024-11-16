package routes

import (
	"github.com/gin-gonic/gin"
	"pureheroky.com/server/handlers"
	"pureheroky.com/server/middlewares"
	"pureheroky.com/server/models"
)

/*
SetupRoutes initializes and sets up API endpoints for the Gin router.
It applies the CorsMiddleware and defines various routes handled by the application.
*/
func SetupRoutes(r *gin.Engine, app *models.App) {
	r.Use(middlewares.CorsMiddleware())
	r.Use(middlewares.LogIPMiddleware())

	handler := handlers.NewHandler(app.MongoClient)

	r.GET("/getuservalue/:valuetype", handler.GetUserValueHandler)
	r.GET("/getimage/:collname/:queryname", handler.GetImageHandler)
	r.GET("/getproject/:projectid", handler.GetProjectHandler)
	r.GET("/uploadimage/:query/:id/:imagename", handler.UploadImageHandler)
	r.POST("/projects", handler.CreateProjectHandler)
}
