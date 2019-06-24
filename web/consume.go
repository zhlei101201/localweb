package web

import (
	"localweb/db"
	"localweb/log"
)

func Consume(channel chan *urlInfo, id int) {
	var (
		task *urlInfo
		body []byte
		code int
		mode string
		err error
	)

	log.Log.Infof("Consume[%d] is starting...\n", id)
	for task = range channel {
		log.Log.Infof("Consume[%d] Rcv: %+v\n", id, task)
		if task.source {
			body, code, mode, err = HttpGet(task)
			if err != nil || code != 200 {
				continue
			}
		} else {
			body, mode, err = db.GetPath(task.path) //db中包含完整的路径
			if err != nil {
				log.Log.Warningf("db.GetPath(%v) failed: %v\n", task, err)
				continue
			}
		}
		task.mode = mode

		log.Log.Infof("Download %s: body=%d, code=%d\n", task.path, len(body), code)
		c := PageParse(task, body)

		//数据存储
		if task.source {
			if err := db.AddPath(c.task.path, c.title, c.keywords, c.description, c.task.mode, true, []byte(c.body)); err != nil {
				log.Log.Errorf("db.AdPath(%v) failed: %v\n", c, err)
			}
		}
	}
}
