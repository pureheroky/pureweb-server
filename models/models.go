package models

import "go.mongodb.org/mongo-driver/mongo"

type App struct {
	MongoClient *mongo.Client
}

// UploadImage represents body of request to update profile image
type UploadImage struct {
	ID        string `json:"id"`
	Query     string `json:"query"`
	ImageName string `json:"imagename"`
}

// User represents a user document in the MongoDB collection.
type User struct {
	ID       string   `json:"_id"`
	Age      int64    `json:"age"`
	Gitlink  string   `json:"gitlink"`
	Image    string   `json:"image"`
	Username string   `json:"name"`
	Tglink   string   `json:"tglink"`
	Status   bool     `json:"status"`
	Skills   []string `json:"skills"`
}

// Skills represents a list of skills.
type Skills struct {
	Skills []string `json:"skills"`
}

// Project represents a project document in the MongoDB collection.
type Project struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	PrjGitlink  string `json:"prjgitlink"`
	PrjImage    string `json:"image"`
	PrjStatus   string `json:"prjstatus"`
	PrjComplete string `json:"prjcomplete"`
	PrjDate     string `json:"prjdate"`
	PrjDesc     string `json:"prjDesc"`
	PrjWeblink  string `json:"prjweblink"`
}
