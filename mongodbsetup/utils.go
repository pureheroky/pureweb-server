package mongodbsetup

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func SetUserSkills(coll *mongo.Collection, skills []string) {
	filter := bson.D{{Key: "name", Value: "pureheroky"}}
	update := bson.M{
		"$push": bson.M{
			"skills": skills,
		},
	}

	result, err := coll.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Matched %v documents and updated %v documents.\n", result.MatchedCount, result.ModifiedCount)
}

func SetUpProjects(client *mongo.Client, id string, name string, git string, img string, status string, complete string, date string, desc string, weblink string) {
	prjcoll := client.Database("pureweb").Collection("projects")

	docs := []interface{}{
		Project{ID: id, Name: name, PrjGitlink: git, PrjImage: img, PrjStatus: status, PrjComplete: complete, PrjDate: date, PrjDesc: desc, PrjWeblink: weblink},
	}

	_, err := prjcoll.InsertMany(context.TODO(), docs)

	if err != nil {
		panic(err)
	}
}
