package utils

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"reflect"
	"strings"
	"unicode"
)

/*
GetEnvValue load environment variables from a .env file.
Return the value of the specified environment variable.
*/
func GetEnvValue(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv(key)
}

/*
Capitalize if the string is empty, return it as is.
Capitalize the first letter of the string and return the result.
*/
func Capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	runes := []rune(s)
	runes[0] = unicode.ToUpper(runes[0])
	return string(runes)
}

/*
GetFieldValue Use reflection to get the value of a nested field in a struct.
Capitalize field names and return the field value as a string.
*/
func GetFieldValue(v interface{}, key string) string {
	rv := reflect.ValueOf(v)
	keys := strings.Split(key, ".")

	for _, k := range keys {
		k = Capitalize(k)
		rv = reflect.Indirect(rv).FieldByName(k)
		if !rv.IsValid() {
			fmt.Printf("No such field: %s in obj\n", k)
			return ""
		}
	}

	return fmt.Sprintf("%v", rv.Interface())
}
