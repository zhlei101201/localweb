package web

import (
	"fmt"
	"io/ioutil"
	"localweb/db"
	"localweb/log"
	"net/http"
	"sort"
	"strings"
	"time"
)

func HttpGet(task *urlInfo) (body []byte, code int, mode string, err error) {
	var response *http.Response

	if task.download == Index{
		response, err = http.Get(webConfig.schema + webConfig.domain)
	} else if task.download[0] == '/' {
		response, err = http.Get(webConfig.schema + webConfig.domain + task.download)
	} else { //含域名
		response, err = http.Get(task.download)
	}
	if err != nil || response.StatusCode == 503 { //503临时错误，为防限流暂停1s
		if err != nil {
			log.Log.Warningf("http.Get(%s) failed: %v\n", task.download, err)
		} else {
			log.Log.Warningf("http.Get(%s) failed: StatusCode=%v\n", task.download, response.StatusCode)
			time.Sleep(time.Duration(webConfig.stime) * time.Millisecond)
		}
		lock.Lock()
		delete(allPaths, task.path)
		allPathsNum -= 1
		lock.Unlock()
		return
	}
	defer response.Body.Close()


	if code = response.StatusCode; code != 200 {
		log.Log.Warningf("http.Get(%s) failed: StatusCode=%d\n", task.download, code)
		//不再重新下载
		return
	}

	body, err = ioutil.ReadAll(response.Body)
	if err != nil {
		log.Log.Errorf("ioutil.ReadAll(%s) failed: %v\n", task.download, err)
		lock.Lock()
		delete(allPaths, task.path)
		allPathsNum -= 1
		lock.Unlock()
		return
	}

	mode = response.Header.Get("Content-Type")

	if task.path == Index { //首页http返回301，自动跳转至https
		if webConfig.schema != response.Request.URL.Scheme + "://" {
			webConfig.schema = response.Request.URL.Scheme + "://"
		}
	}

	return
}

func GetLocalHtml(path string) (context []byte, mode string) {
	log.Log.Infof("  GetLocalHtml %s\n", path)
	if len(path) <= 1 {
		path = Index
	}

	var err error
	if _, ok := allPaths[path]; ok { // 检查是否已经下载，已经下载该页面
		context, mode, err = db.GetPath(path) //db中包含完整的路径
		if err != nil {
			context = []byte("排队等待下载中，请稍后刷新页面")
			mode = "text/html; charset=utf-8"
			log.Log.Warningf("    %s is in queue\n", path)
			return
		}
		//同步核对该网页的所有链接（a.href）是否已经缓存
		if webConfig.mode != MODE_NONE && strings.Contains(mode, "text/html") {
			go LocalPageParse(path, mode, context)
		}
		return
	}

	if path == Index {
		context, mode = GetCommands()
		return
	}

	context = []byte("本地缓存无此页面信息")
	mode = "text/html; charset=utf-8"
	return

	/* 没有必要，LocalPageParse满足需求
	if webConfig.mode == MODE_NONE {
	}


	NewTasks(path, path, 2, true)

	context = []byte("缓存中无此网页，已经发起自动同步，请稍后刷新页面")
	mode = "text/html; charset=utf-8"

	return
	*/
}

func GetAllPaths() (body []byte, mode string) {
	var tnum, fnum int
	var all []string
	var isPrint bool

	lock.Lock()
	str := fmt.Sprintf("Count = %d<br />\n", allPathsNum)
	if allPathsNum < 60000 {
		for k, v := range allPaths {
			all = append(all, fmt.Sprintf("%s : %v", k, v))
			if v {
				tnum += 1
			} else {
				fnum += 1
			}
		}
		isPrint = true
	}
	lock.Unlock()

	body = append(body, []byte(fmt.Sprintf("waitQueue.length=%d<br />\n", len(waitsQueue)))...)
	for i := 0; i < webConfig.tasks; i++ {
		body = append(body, []byte(fmt.Sprintf("tasksQueue[%d].length=%d<br />\n", i, len(tasksQueue[i])))...)
	}
	body = append(body, []byte(str)...)
	if isPrint {
		body = append(body, []byte(fmt.Sprintf("isTure = %d, isFalse = %d<br />\n", tnum, fnum))...)
		sort.Strings(all)
		for i, v := range all {
			body = append(body, []byte(fmt.Sprintf("%d : %s<br />\n", i + 1, v))...)
		}
	}
	mode = "text/html; charset=utf-8"

	return
}

func GetCommands() (context []byte, mode string) {
	cmds, err := db.GetCommands()
	if err != nil {
		context = []byte("获取Commands失败")
		mode = "text/html; charset=utf-8"
		log.Log.Warningf("db.GetCommands() failed: %v\n", err)
		return
	}

	context = append(context, []byte(fmt.Sprintf("Total = %d<br />\n", len(cmds)))...)
	for _, v := range cmds {
		context = append(context, []byte(fmt.Sprintf("<br />\n%v<br />\n", v))...)
	}

	mode = "text/html; charset=utf-8"
	return
}
