# localweb
golang开发的网站镜像工具，可以将网站复制到本地，并以web的形式在本地操作

go build -o localweb main.go

localweb version : 1.0.0
Usage: localweb [-w webName] [-tn taskNumber] [-dl depthLevel] [-wl waitLength] [-wm workMode] [-ln localName] [-ld logDir] [-ll logLevel] [-lp localPort] [-nt] [-st sleepTime]

Options:    
  &emsp;-dl int    
        &emsp;&emsp;&emsp;Web网页链接的深度,[1-100] (default 2)    
  &emsp;-ld string    
        &emsp;&emsp;&emsp;日志文件存储目录    
  &emsp;-ll string    
        &emsp;&emsp;&emsp;日志级别,[DEBUG|INFO|WARNING|ERROR] (default "WARNING")    
  &emsp;-ln string    
        &emsp;&emsp;&emsp;本地镜像网站的存储文件名称,默认同域名    
  &emsp;-lp int    
        &emsp;&emsp;&emsp;本地镜像网站的服务端口,[80-65535] (default 80)    
  &emsp;-nt    
        &emsp;&emsp;&emsp;部分网站声明的编码格式有误，允许不做编码格式转换(默认转换)    
  &emsp;-st int    
        &emsp;&emsp;&emsp;遇503错误，暂停运行的时长（毫秒）,[1-1000000] (default 1000)    
  &emsp;-tn int    
        &emsp;&emsp;&emsp;网页下载时的并行任务数,[1-100] (default 1)    
  &emsp;-wl int    
        &emsp;&emsp;&emsp;等待下载的任务队列长度,[10000-10000000] (default 100000)    
  &emsp;-wm string    
        &emsp;&emsp;&emsp;工作模式：[NONE|AUTO|APPEND|RELOAD] (default "NONE")    
  &emsp;-wn string    
        &emsp;&emsp;&emsp;目标网站名称,[如：https://www.xxx.com]    
