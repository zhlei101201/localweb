package db

import "fmt"

func SetName(v string) {
	dbConfig.Name = v
}

func GetName() string {
	return dbConfig.Name
}

func Display() string {
	return fmt.Sprintf("%+v", dbConfig)
}