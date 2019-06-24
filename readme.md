<font size=4>

# localweb
&emsp;&emsp;&emsp;golang开发的网站镜像工具，可以将网站复制到本地，并以web的形式在本地操作（本地压缩和sqlite3存储），支持css&js，支持动态网页。阅读效果与直接登录原网站相同。   
&emsp;&emsp;&emsp;暂不支持登录：各个网站差别较大，缺少通用性。   
  
### 工作模式说明（参数-wm workMode)   
**&emsp;NONE：**仅作为web服务器，供浏览器访问本地镜像文件，即使发现缺少网页也不会与网站进行同步。   
**&emsp;AUTO：**作为web服务器，供浏览器访问，浏览器访问时，会主动检查访问页面上的所有链接是否已经全部缓存到本地，若发现缺少则自动向网站发起同步。   
**&emsp;APPEND：**启动后，从入口（-wn参数指定）开始检查所有网页是否全部缓存到本地，若有缺少则自动向网站发起同步，检查深度由参数-dl指定。   
**&emsp;RELOAD：**启动后，无论本地是否有缓存，均向网站发起同步，以适应网站的页面更新。深度有参数-ld指定。

   
## 编译
go build -o localweb main.go   
*提醒：windows下是localweb.exe*


## localweb version : 1.0.0
**Usage: localweb [-w webName] [-ln localName] [-dl depthLevel] [-tn taskNumber] [-wl waitLength] [-wm workMode] [-lp localPort] [-ld logDir] [-ll logLevel] [-nt] [-st sleepTime]**

## Options:    
  **&emsp;-wn string**    
        &emsp;&emsp;&emsp;目标网站名称,[如：https://www.xxx.com]    
  **&emsp;-ln string**    
        &emsp;&emsp;&emsp;本地镜像网站的存储文件名称,默认同域名    
  **&emsp;-dl int**    
        &emsp;&emsp;&emsp;Web网页链接的深度,[1-100] (default 2)    
  **&emsp;-tn in**t    
        &emsp;&emsp;&emsp;网页下载时的并行任务数,[1-100] (default 1)    
  **&emsp;-wl int**    
        &emsp;&emsp;&emsp;等待下载的任务队列长度,[10000-10000000] (default 100000)    
  **&emsp;-wm string**    
        &emsp;&emsp;&emsp;工作模式：[NONE|AUTO|APPEND|RELOAD] (default "NONE")    
  **&emsp;-lp int**    
        &emsp;&emsp;&emsp;本地镜像网站的服务端口,[80-65535] (default 80)    
  **&emsp;-ld string**    
        &emsp;&emsp;&emsp;日志文件存储目录    
  **&emsp;-ll string**    
        &emsp;&emsp;&emsp;日志级别,[DEBUG|INFO|WARNING|ERROR] (default "WARNING")    
  **&emsp;-nt**    
        &emsp;&emsp;&emsp;部分网站声明的编码格式有误，允许不做编码格式转换(默认转换)    
  **&emsp;-st int**    
        &emsp;&emsp;&emsp;遇503错误，暂停运行的时长（毫秒）,[1-1000000] (default 1000)    

</font>
