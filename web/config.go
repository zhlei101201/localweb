package web

import (
	"fmt"
	"strings"
)

func SetSchema(v string) {
	webConfig.schema = v
}

func GetSchema() string {
	return webConfig.schema
}

func SetDomain(v string) {
	webConfig.domain = v
}

func GetDomain() string {
	return webConfig.domain
}

func SetPath(v string) {
	if len(v) <= 1 {
		v = Index
	}

	webConfig.path = v
	if idx := strings.LastIndexByte(v, '/'); idx >= 0 {
		webConfig.dir = v[:idx+1]
	} else {
		webConfig.dir = "/"
	}
}

func GetPath() string {
	return webConfig.path
}

func SetQuery(v string) {
	webConfig.query = v
}

func GetQuery() string {
	return webConfig.query
}

func SetTasks(v int) {
	webConfig.tasks = v
}

func GetTasks() int {
	return webConfig.tasks
}

func SetDepth(v int) {
	webConfig.depth = v
}

func GetDepth() int {
	return webConfig.depth
}

func SetWlen(v int) {
	webConfig.wlen = v
}

func GetWlen() int {
	return webConfig.wlen
}

func SetMode(v int) {
	webConfig.mode = v
}

func GetMode() int {
	return webConfig.mode
}

func SetTrans(v bool) {
	webConfig.notrans = v
}

func GetTrans() bool {
	return webConfig.notrans
}

func SetSTime(v int64) {
	webConfig.stime = v
}

func GetSTime() int64 {
	return webConfig.stime
}

func Display() string {
	return fmt.Sprintf("%+v", webConfig)
}
