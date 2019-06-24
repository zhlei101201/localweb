package log

import (
	"github.com/cxr29/log"
)

func Open() (err error) {
	logfile, err = log.Open(logConfig.base, int64(logConfig.second), int64(logConfig.size), int64(logConfig.max))
	if err != nil {
		return
	}

	Log = log.New(logfile, "", log.Ltime, log.NameLevel(logConfig.level))
	return
}

func Close() {
	logfile.Close()
}