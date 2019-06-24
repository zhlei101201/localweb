package main

import (
	"flag"
	"fmt"
	"localweb/db"
	"localweb/log"
	"localweb/web"
	"net/http"
	"os"
	"regexp"
	"strings"
	_  "net/http/pprof"
)

var (
	program		= "localweb"
	version		= "1.0.0"

	webName		string
	localName	string
	localPort	int

	logLevel	string
	logDir		string

	taskNumber	int
	depthLevel	int
	waitLength	int
	workMode	string

	noTrans		bool
	sleepTime	int64
)

func usage() {
	fmt.Printf("%s version : %s\n", program, version)
	fmt.Printf("Usage: %s [-w webName] [-tn taskNumber] [-dl depthLevel] [-wl waitLength] [-wm workMode] [-ln localName] [-ld logDir] [-ll logLevel] [-lp localPort] [-nt] [-st sleepTime]\n", program)
	fmt.Printf("\nOptions:\n")
	flag.PrintDefaults()
}

func init() {
	flag.StringVar(&webName, "wn", "", "目标网站名称,[如：https://www.xxx.com]")
	flag.IntVar(&taskNumber, "tn", 1, "网页下载时的并行任务数,[1-100]")
	flag.IntVar(&depthLevel, "dl", 2, "Web网页链接的深度,[1-100]")
	flag.IntVar(&waitLength, "wl", 100000, "等待下载的任务队列长度,[10000-10000000]")
	flag.StringVar(&workMode, "wm", "NONE", "工作模式：[NONE|AUTO|APPEND|RELOAD]")
	flag.StringVar(&localName, "ln", "", "本地镜像网站的存储文件名称,默认同域名")
	flag.StringVar(&logDir, "ld", "", "日志文件存储目录")
	flag.StringVar(&logLevel, "ll", "WARNING", "日志级别,[DEBUG|INFO|WARNING|ERROR]")
	flag.IntVar(&localPort, "lp", 80, "本地镜像网站的服务端口,[80-65535]")
	flag.BoolVar(&noTrans, "nt", false, "部分网站声明的编码格式有误，允许不做编码格式转换(默认转换)")
	flag.Int64Var(&sleepTime, "st", 1000, "遇503错误，暂停运行的时长（毫秒）,[1-1000000]")
	flag.Usage = usage
}

func main() {
	flag.Parse()

	if len(webName) == 0 && len(localName) == 0 {
		flag.Usage()
		return
	}

	pattern := regexp.MustCompile(`^(https?://)?([\w.-]+\.[0-9a-zA-Z]+)(/[\w./-]*)?(.*?)$`)
	matcher := pattern.FindStringSubmatch(strings.ToLower(webName))
	if matcher == nil {
		fmt.Printf("Incorrect webName: %s\n", webName)
		return
	}
	if len(matcher[1]) > 0 {
		web.SetSchema(matcher[1])
	}
	if len(matcher[2]) > 0 {
		web.SetDomain(matcher[2])
	} else {
		fmt.Printf("Incorrect webName: %s\n", webName)
		return
	}
	web.SetPath(matcher[3])
	if len(matcher[4]) > 0 {
		web.SetQuery(matcher[4])
	}
	if taskNumber < 1 || taskNumber > 100 {
		fmt.Printf("Incorrect taskNumber: %d\n", taskNumber)
		return
	}
	web.SetTasks(taskNumber)
	if depthLevel < 1 || depthLevel > 100 {
		fmt.Printf("Incorrect depthLevel: %d\n", depthLevel)
		return
	}
	web.SetDepth(depthLevel)
	if waitLength < 10000 || waitLength > 10000000 {
		fmt.Printf("Incorrect waitLength: %d\n", waitLength)
		return
	}
	web.SetWlen(waitLength)

	switch strings.ToUpper(workMode) {
	case "NONE":
		web.SetMode(web.MODE_NONE)
	case "AUTO":
		web.SetMode(web.MODE_AUTO)
	case "RELOAD":
		web.SetMode(web.MODE_RELOAD)
	case "APPEND":
		web.SetMode(web.MODE_APPEND)
	default:
		fmt.Printf("Incorrect workMode: %s\n", workMode)
		return
	}

	if len(localName) > 0 {
		db.SetName(localName)
	} else {
		db.SetName(matcher[2]+".db")
	}

	if len(logDir) > 0 {
		log.SetDir(logDir)
	}

	switch strings.ToUpper(logLevel) {
	case "DEBUG":
	case "INFO":
	case "WARNING":
	case "ERROR":
	default:
		fmt.Printf("Incorrect logLevel: %s\n", logLevel)
		return
	}
	log.SetLevel(strings.ToUpper(logLevel))

	if localPort < 80 || localPort > 65535 {
		fmt.Printf("Incorrect localPort: %d\n", localPort)
		return
	}

	web.SetTrans(noTrans)
	if sleepTime < 1 || sleepTime > 1000000{
		fmt.Printf("Incorrect sleepTime: %d\n", sleepTime)
		return
	}
	web.SetSTime(sleepTime)

	log.Open()
	defer log.Close()

	str := fmt.Sprintf("Command:%+v<br />\n...WebConfig:%+v<br />\n...DbConfig:%+v<br />\n...LogConfig:%+v<br />\n", os.Args, web.Display(), db.Display(), log.Display())
	log.Log.Criticalf("%s\n", str)

	if err := db.Open(str); err != nil {
		log.Log.Errorf("db.Open() failed: %v\n", err)
		return
	}
	defer db.Close()

	if err := web.Init(); err != nil {
		log.Log.Errorf("web.Init() failed: %v\n", err)
		return
	}

	http.HandleFunc("/", app)
	http.HandleFunc("/system/GetAllPaths", getAllPaths)
	http.HandleFunc("/system/GetAllCmds", getAllCmds)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", localPort), nil); err != nil {
		log.Log.Errorf("ListenAndServer(%d) failed: %v\n", localPort, err)
		return
	}
}

func app(w http.ResponseWriter, r *http.Request) {
	context, mode := web.GetLocalHtml(r.RequestURI)
	w.Header().Set("Content-Type", mode)
	w.WriteHeader(200)
	w.Write(context)
}

func getAllPaths(w http.ResponseWriter, r *http.Request) {
	context, mode := web.GetAllPaths()
	w.Header().Set("Content-Type", mode)
	w.WriteHeader(200)
	w.Write(context)
}

func getAllCmds(w http.ResponseWriter, r *http.Request) {
	context, mode := web.GetCommands()
	w.Header().Set("Content-Type", mode)
	w.WriteHeader(200)
	w.Write(context)
}
