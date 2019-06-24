# localweb
golang开发的网站镜像工具，可以将网站复制到本地，并以web的形式在本地操作

go build -o localweb main.go

localweb version : 1.0.0
Usage: localweb [-w webName] [-tn taskNumber] [-dl depthLevel] [-wl waitLength] [-wm workMode] [-ln localName] [-ld logDir] [-ll logLevel] [-lp localPort] [-nt] [-st sleepTime]

Options:    
  -dl int    
          Web网页链接的深度,[1-100] (default 2)    
  -ld string    
        日志文件存储目录    
  -ll string    
        日志级别,[DEBUG|INFO|WARNING|ERROR] (default "WARNING")    
  -ln string    
        本地镜像网站的存储文件名称,默认同域名    
  -lp int    
        本地镜像网站的服务端口,[80-65535] (default 80)    
  -nt    
        部分网站声明的编码格式有误，允许不做编码格式转换(默认转换)    
  -st int    
        遇503错误，暂停运行的时长（毫秒）,[1-1000000] (default 1000)    
  -tn int    
        网页下载时的并行任务数,[1-100] (default 1)    
  -wl int    
        等待下载的任务队列长度,[10000-10000000] (default 100000)    
  -wm string    
        工作模式：[NONE|AUTO|APPEND|RELOAD] (default "NONE")    
  -wn string    
        目标网站名称,[如：https://www.xxx.com]    
