package log

import (
	"github.com/cxr29/log"
	"io"
)

type config struct {
	dir			string
	level		string
	base		string
	second		int
	size		int
	max			int
}

var (
	logConfig = &config{
	dir: "",
	level: "DEBUG",
	base: "log",
	second: 86400,
	size: 16*1024*1024,
	max: 200,
	}

	logfile io.WriteCloser
	Log *log.Logger
)

