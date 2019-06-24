package web

import (
	"localweb/log"
	"time"
)

func NewTasks(download, path string, depth int, parse bool) {
	var full bool

	if depth < 0 {
		return
	}
	if len(path) <= 1 { //首页面
		path = Index
		download = Index
	}
	task := &urlInfo{
		download: download,
		path: path,
		depth: depth,
	}

	lock.Lock()
	v, ok := allPaths[path]
	if !ok || !v { //适用于AUTO、APPEND、RELOAD三种模式
		task.source = true
	} else if webConfig.mode == MODE_APPEND  && parse {
		task.source = false
	} else {
		lock.Unlock()
		return
	}

	select {
	case waitsQueue <- task:
		if !ok || !v {
			allPaths[path] = true
		}
		if !ok {
			allPathsNum += 1
		}
	default:
		full = true //避免引起内存过大消耗，可以通过APPEND模式触发重新下载
	}
	lock.Unlock()

	if full {
		log.Log.Warningf("WaitQueue is full<%d>: %+v\n", len(waitsQueue), task)
	} else {
		log.Log.Infof("SendWaitQueue: %+v\n", task)
	}
}

func WaitTasks() {
	var (
		task *urlInfo
		id = -1
	)
	for task = range waitsQueue {
		id = SendTaskQueue(id, task)
		log.Log.Infof("SendTaskQueue[%d]: %+v\n", id, task)
	}
}

func SendTaskQueue(id int, task *urlInfo) int {
	for {
		if id = FindIdelChannel(id); id < 0 {
			log.Log.Infof("TaskQueue is full, retry in 1s\n")
			time.Sleep(time.Duration(1) * time.Second)
		} else {
			tasksQueue[id] <- task
			break
		}
	}
	return id
}

func FindIdelChannel(id int) int {
	if id < -1 || id >= webConfig.tasks {
		id = -1
	}
	for i := 0; i < webConfig.tasks; i++ {
		if id += 1; id >= webConfig.tasks {
			id = 0
		}
		if len(tasksQueue[id]) < TaskQueueLen {
			return id
		}
	}

	return -1
}
