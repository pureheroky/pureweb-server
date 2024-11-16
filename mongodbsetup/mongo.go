package mongodbsetup

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo/options"
	"os"
	"pureheroky.com/server/config"
	"pureheroky.com/server/models"
	"pureheroky.com/server/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

/*
SetupMongodbClient connecting to database and return
connection to it
*/
func SetupMongodbClient() (*mongo.Client, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(config.MongoURI).SetServerAPIOptions(serverAPI)
	return mongo.Connect(context.Background(), opts)
}

/*
SaveImageInDB read the image file from the specified path.
Encode the file data to base64 and update the MongoDB document with the encoded image.
Return any errors encountered during the process.
*/
func SaveImageInDB(coll *mongo.Collection, id, imagePath string) error {
	fileData, err := os.ReadFile(imagePath)
	if err != nil {
		return err
	}

	base64Data := base64.StdEncoding.EncodeToString(fileData)

	filter := bson.M{"id": id}
	update := bson.M{"$set": bson.M{"image": base64Data}}
	_, err = coll.UpdateOne(context.TODO(), filter, update)
	return err
}

/*
GetImage find a document by ID and retrieve the base64 encoded image data.
Decode the image data and return the image bytes and content type.
Return errors if the document is not found or decoding fails.
*/
func GetImage(coll *mongo.Collection, id string) ([]byte, string, error) {
	var result bson.M
	filter := bson.D{{Key: "id", Value: id}}

	err := coll.FindOne(context.TODO(), filter).Decode(&result)
	if err != nil {
		return nil, "", err
	}

	imageData, ok := result["image"].(string)
	if !ok {
		return nil, "", errors.New("image not found or invalid type")
	}

	decodedImage, err := base64.StdEncoding.DecodeString(imageData)
	if err != nil {
		return nil, "", err
	}

	contentType := "image/png"

	return decodedImage, contentType, nil
}

/*
GetUserValue find a user document by name and decode it into a User struct.
Retrieve the value of a specified field using reflection.
Return the field value as a string.
*/
func GetUserValue(coll *mongo.Collection, title string, key string) (string, error) {
	var result bson.M
	err := coll.FindOne(context.TODO(), bson.D{{Key: "name", Value: title}}).Decode(&result)
	if err != nil {
		return "", err
	}

	user := &result
	jsonData, err := json.MarshalIndent(user, "", "    ")
	if err != nil {
		return "", err
	}

	var resStruct models.User
	err = json.Unmarshal(jsonData, &resStruct)

	if err != nil {
		fmt.Printf("Error unmarshalling json: %v\n", err)
		return "", err
	}
	return utils.GetFieldValue(resStruct, key), nil
}

/*
GetProject find a project document by ID and return the result.
If the project is not found, print an error message.
Return the project document or nil if not found.
*/
func GetProject(coll *mongo.Collection, id string) (*bson.M, error) {
	var result bson.M
	err := coll.FindOne(context.TODO(), bson.D{{Key: "id", Value: id}}).Decode(&result)

	if errors.Is(err, mongo.ErrNoDocuments) {
		fmt.Printf("No project was found with the id %s\n", id)
		return nil, nil
	}

	if err != nil {
		panic(err)
	}
	return &result, nil
}

//func SetUserSkills(coll *mongo.Collection, skills []string) {
//	filter := bson.D{{Key: "name", Value: "pureheroky"}}
//	update := bson.M{
//		"$push": bson.M{
//			"skills": skills,
//		},
//	}
//
//	result, err := coll.UpdateOne(context.TODO(), filter, update)
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	fmt.Printf("Matched %v documents and updated %v documents.\n", result.MatchedCount, result.ModifiedCount)
//}
//
//func SetUpProjects(client *mongo.Client, id string, name string, git string, img string, status string, complete string, date string, desc string, weblink string) {
//	prjcoll := client.Database("pureweb").Collection("projects")
//
//	docs := []interface{}{
//		Project{ID: id, Name: name, PrjGitlink: git, PrjImage: img, PrjStatus: status, PrjComplete: complete, PrjDate: date, PrjDesc: desc, PrjWeblink: weblink},
//	}
//
//	_, err := prjcoll.InsertMany(context.TODO(), docs)
//
//	if err != nil {
//		panic(err)
//	}
//}
