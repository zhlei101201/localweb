package log

import "fmt"

func SetDir(v string) {
	logConfig.dir = v
}

func GetDir() string {
	return logConfig.dir
}

func SetLevel(v string) {
	logConfig.level = v
}

func GetLevel() string {
	return logConfig.level
}

func SetSize(v int) {
	logConfig.size = v
}

func GetSize() int {
	return logConfig.size
}

func Display() string {
	return fmt.Sprintf("%+v", logConfig)
}