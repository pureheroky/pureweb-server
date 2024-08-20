package mongodbsetup

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strings"
	"unicode"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

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
	ID          string             `json:"id"`
	Name        string             `json:"name"`
	PrjGitlink  string             `json:"prjgitlink"`
	PrjImage    string             `json:"image"`
	PrjStatus   string             `json:"prjstatus"`
	PrjComplete string             `json:"prjcomplete"`
	PrjDate     string             `json:"prjdate"`
	PrjDesc     string             `json:"prjDesc"`
	PrjWeblink  string             `json:"prjweblink"`
}

func capitalize(s string) string {
	/*
		If the string is empty, return it as is.
		Capitalize the first letter of the string and return the result.
	*/
	if len(s) == 0 {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

func SaveImageInDB(coll *mongo.Collection, id, imagePath string) error {
	/*
		Read the image file from the specified path.
		Encode the file data to base64 and update the MongoDB document with the encoded image.
		Return any errors encountered during the process.
	*/
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

func getFieldValue(v interface{}, key string) string {
	/*
		Use reflection to get the value of a nested field in a struct.
		Capitalize field names and return the field value as a string.
	*/
	rv := reflect.ValueOf(v)
	keys := strings.Split(key, ".")

	for _, k := range keys {
		k = capitalize(k)
		rv = reflect.Indirect(rv).FieldByName(k)
		if !rv.IsValid() {
			fmt.Printf("No such field: %s in obj\n", k)
			return ""
		}
	}

	return fmt.Sprintf("%v", rv.Interface())
}

func GetImage(coll *mongo.Collection, id string) ([]byte, string, error) {
	/*
		Find a document by ID and retrieve the base64 encoded image data.
		Decode the image data and return the image bytes and content type.
		Return errors if the document is not found or decoding fails.
	*/
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

func GetUserValue(coll *mongo.Collection, title string, key string) string {
	/*
		Find a user document by name and decode it into a User struct.
		Retrieve the value of a specified field using reflection.
		Return the field value as a string.
	*/
	var result bson.M
	err := coll.FindOne(context.TODO(), bson.D{{Key: "name", Value: title}}).Decode(&result)

	if err != nil {
		panic(err)
	}

	user := &result
	jsonData, err := json.MarshalIndent(user, "", "    ")
	if err != nil {
		panic(err)
	}

	var resStruct User

	err = json.Unmarshal(jsonData, &resStruct)
	if err != nil {
		fmt.Printf("Error unmarshalling json: %v\n", err)
		panic(err)
	}
	return getFieldValue(resStruct, key)
}

func GetProject(coll *mongo.Collection, id string) *bson.M {
	/*
		Find a project document by ID and return the result.
		If the project is not found, print an error message.
		Return the project document or nil if not found.
	*/
	var result bson.M
	err := coll.FindOne(context.TODO(), bson.D{{Key: "id", Value: id}}).Decode(&result)

	if err == mongo.ErrNoDocuments {
		fmt.Printf("No project was found with the id %s\n", id)
		return nil
	}

	if err != nil {
		panic(err)
	}
	return &result
}

