package config

import (
	"pureheroky.com/server/utils"
)

var Client = utils.GetEnvValue("CLIENTCOLL")
var Database = utils.GetEnvValue("DATABASE")
var Projects = utils.GetEnvValue("PROJECTCOLL")
var MongoURI = utils.GetEnvValue("MONGOURI")
