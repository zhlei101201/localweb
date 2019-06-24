package web

import "sync"

const (
	MODE_NONE	= iota
	MODE_AUTO
	MODE_APPEND
	MODE_RELOAD

	CHARSET_GBK
	CHARSET_GB18030
	CHARSET_UTF8
)

const (
	Index			= "/index.html"
	Suffix			= "/x.suffix"
	TaskQueueLen	= 1000
)

type config struct {
	schema		string
	domain		string
	path		string
	dir			string
	query		string
	tasks		int
	depth		int
	wlen		int
	mode		int
	notrans		bool
	stime		int64
}

type urlInfo struct {
	download string
	path	string
	mode	string
	depth	int
	charset int
	source	bool
}

type Context struct {
	task *urlInfo
	dir string
	body string
	title string
	keywords string
	description string
	err error
}

var (
	webConfig = &config{
	schema: "http://",
	tasks: 1,
	depth: 3,
	mode: MODE_NONE,
	}

	lock sync.Mutex
	allPaths map[string]bool //url冲突检测表
	allPathsNum int
	tasksQueue []chan *urlInfo //下载任务队列
	waitsQueue chan *urlInfo //等候下载队列，避免下载任务队列过长
)
