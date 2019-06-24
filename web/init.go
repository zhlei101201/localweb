package web

import (
	"localweb/db"
	"localweb/log"
)

func Init() error{
	paths, err := db.GetAllPath()
	if err != nil {
		log.Log.Errorf("db.GetAllPath() failed: %v\n", err)
		return err
	}
	log.Log.Infof("db.GetAllPath().length = %d\n", len(paths))

	allPaths = make(map[string]bool)
	for i, v := range paths {
		log.Log.Infof("  %d : %v\n", i, v.Path)
		if webConfig.mode == MODE_RELOAD {
			allPaths[v.Path] = false
		} else {
			allPaths[v.Path] = true
		}
	}
	allPathsNum = len(paths)
	log.Log.Infof("allPathsNum = %d\n", allPathsNum)

	if webConfig.mode == MODE_NONE { //NONE模式仅创建本地web
		return nil
	}

	waitsQueue = make(chan *urlInfo, webConfig.wlen)
	go WaitTasks() //等待队列任务处理

	tasksQueue = make([]chan *urlInfo, webConfig.tasks)
	for i := 0; i < webConfig.tasks; i++ {
		tasksQueue[i] = make(chan *urlInfo, TaskQueueLen)
		go Consume(tasksQueue[i], i)
	}

	if webConfig.mode != MODE_AUTO { //AUTO模式，在本地无请求页面时进行自动同步
		NewTasks(webConfig.path, webConfig.path, webConfig.depth, true)
	}

	return nil
}
