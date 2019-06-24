package web

import (
	"fmt"
	"golang.org/x/text/encoding/simplifiedchinese"
	"localweb/log"
	"regexp"
	"strings"
)

func NewContext(task *urlInfo, body []byte) *Context {
	return &Context{
		task: task,
		body: string(body),
	}
}

func PageParse(task *urlInfo, body []byte) *Context {
	c := NewContext(task, body)

	//检测并转换中午字符集至UTF-8
	if !webConfig.notrans && (task.source && (strings.Contains(task.mode, "text/html") || strings.Contains(task.mode, "text/css") || strings.Contains(task.mode, "javascript"))) {
		c.FindCharest().TransCharset()
	}
	if strings.Contains(task.mode, "text/css") { //css中包含图片：url()
		if task.source {
			c.GetDir().CssUrlBefore().CssUrlFind().CssUrlAfter()
		} else {
			c.GetDir().CssUrlFind()
		}
	}
	if strings.Contains(task.mode, "text/html") == false {
		return c
	}

	//下载允css、js、图片，允许外站引用，需要优先调用GetDir；href链接下载
	if task.source {
		c.GetDir().CssBefore().CssFind().CssAfter().JsBefore().JsFind().JsAfter().ImgBefore().ImgFind().ImgAfter().ReplaceDomain().HrefFind().DeleteBaidu()
		c.GetTitle().GetKeywords().GetDescription()  //替换本站域名，查询title、keywords、description，a.href链接分析
	} else { //其中，原网页中的外站链接，若未及时下载，会因after替换而无法下载，暂时忽略
		c.GetDir().CssFind().JsFind().ImgFind().HrefFind()
	}

	return c
}

func (c *Context) FindCharest() *Context {	//检测编码格式
	Pattern := regexp.MustCompile(`<meta ([^>]*?)(charset=[\w-]+)([^>]*?)>`)
	matcher := Pattern.FindStringSubmatch(strings.ToLower(c.body))
	if matcher != nil {
		if matcher[2] == "charset=gbk" {
			c.task.charset = CHARSET_GBK
		} else if matcher[2] == "charset=gb18030" {
			c.task.charset = CHARSET_GB18030
		} else {
			c.task.charset = CHARSET_UTF8
		}
	} else {
		c.task.charset = CHARSET_UTF8
	}

	return c
}

func (c *Context) TransCharset() *Context {	//转换字符编码为UTF-8

	if c.task.charset == CHARSET_GBK {
		tr := simplifiedchinese.GBK.NewDecoder()
		dst := make([]byte, len([]byte(c.body)) * 2)
		nDst, _, err := tr.Transform(dst, []byte(c.body), true)
		if err != nil {
			log.Log.Warningf("GBK.decode(%v) failed :%s\n", c.task, err.Error())
		} else if strings.Contains(c.task.mode, "text/html") {
			Pattern := regexp.MustCompile(`<meta ([^>]*?)(charset=[\w-]+)([^>]*?)>`)
			c.body = Pattern.ReplaceAllString(string(dst[:nDst]), `<meta ${1}charset=UTF-8${3}>`)
			log.Log.Infof("Translate(from GBK to UTF-8) sucess!")
		}
	} else if c.task.charset == CHARSET_GB18030 {
		tr := simplifiedchinese.GB18030.NewDecoder()
		dst := make([]byte, len([]byte(c.body)) * 2)
		nDst, _, err := tr.Transform(dst, []byte(c.body), true)
		if err != nil {
			log.Log.Warningf("GB18030.decode(%v) failed :%s\n", c.task, err.Error())
		} else if strings.Contains(c.task.mode, "text/html") {
			Pattern := regexp.MustCompile(`<meta ([^>]*?)(charset=[\w-]+)([^>]*?)>`)
			c.body = Pattern.ReplaceAllString(string(dst[:nDst]), `<meta ${1}charset=UTF-8${3}/>`)
			log.Log.Infof("Translate(from GB18030 to UTF-8) sucess!")
		}
	}

	return c
}

func (c *Context) GetDir() *Context {
	idx := strings.LastIndexByte(c.task.path, '/')
	if idx >= 0 {
		c.dir = c.task.path[:idx+1]
	} else {
		c.dir = "/"
	}

	return c
}

func (c *Context) CssBefore() *Context {
	//优先删除链接中包含的本站域名，避免影响后续的路径判断
	//regstr := fmt.Sprintf(`<link([^>]*?) rel=[\"|\']stylesheet[\"|\']([^>]*?) href=([\"|\'])https?://%s([^>]*?)>`, webConfig.domain)
	regstr := fmt.Sprintf(`<link([^>]*?) href=([\"|\'])https?://%s([^>]*?)>`, webConfig.domain)
	Pattern := regexp.MustCompile(regstr)
	c.body = Pattern.ReplaceAllString(c.body, `<link${1} href=${2}${3}>`)

	return c
}

func (c *Context) CssFind() *Context {
	//Pattern := regexp.MustCompile(`<link([^>]*?) rel=[\"|\']stylesheet[\"|\']([^>]*?) href=[\"|\']((https?:/)?([^>]*?))[\"|\']([^>]*?)>`)
	Pattern := regexp.MustCompile(`<link([^>]*?) href=[\"|\']((https?:/)?([^>]*?))[\"|\']([^>]*?)>`)
	matcher := Pattern.FindAllStringSubmatch(c.body, -1)
	log.Log.Infof("Find link.css=%d, in %v\n", len(matcher), c.task)
	for i, v := range matcher {
		log.Log.Infof("  css.link[%d]: %s\n", i, v[2])
		if !strings.Contains(v[1], "stylesheet") && !strings.Contains(v[5], "stylesheet") {
			continue
		} else if len(v[3]) > 0 { //含外站域名
			NewTasks(v[2], v[4], 1, true)
		} else if len(v[4]) == 0 { //空连接
			continue
		} else if v[4][0] == '/' { //绝对路径（本站）
			NewTasks(v[4], v[4], 1, true)
		} else { //相对路径，添加父页面的路径
			NewTasks(c.dir+v[4], c.dir+v[4], 1, true)
		}
	}

	return c
}

func (c *Context) CssAfter() *Context {
	//替换外站的链接地址，若外站的链接未及时下载，append模式下可能无法下载
	//Pattern := regexp.MustCompile(`<link([^>]*?) rel=[\"|\']stylesheet[\"|\']([^>]*?) href=([\"|\'])https?:/([^>]*?)>`)
	Pattern := regexp.MustCompile(`<link([^>]*?) href=([\"|\'])https?:/([^>]*?)>`)
	c.body = Pattern.ReplaceAllString(c.body, `<link${1} href=${2}${3}>`)

	return c
}

func (c *Context) JsBefore() *Context {
	//优先删除链接中包含的本站域名，避免影响后续的路径判断
	regstr := fmt.Sprintf(`<script([^>]*?) src=([\"|\'])https?://%s([^>]*?)>`, webConfig.domain)
	Pattern := regexp.MustCompile(regstr)
	c.body = Pattern.ReplaceAllString(c.body, `<script${1} src=${2}${3}>`)

	return c
}

func (c *Context) JsFind() *Context {
	Pattern := regexp.MustCompile(`<script([^>]*?) src=[\"|\']((https?:/)?([\w./-]+))[\"|\']([^>]*?)>`)
	matcher := Pattern.FindAllStringSubmatch(c.body, -1)
	log.Log.Infof("Find link.js=%d, in %v\n", len(matcher), c.task)
	for i, v := range matcher {
		log.Log.Infof("  js.link[%d]: %s\n", i, v[2])
		if len(v[3]) > 0 { //含外站域名
			NewTasks(v[2], v[4], 1, false)
		} else if len(v[4]) == 0 { //空链接
			continue
		} else if v[4][0] == '/' { //绝对路径（本站）
			NewTasks(v[4], v[4], 1, false)
		} else { //相对路径，添加父页面的路径
			NewTasks(c.dir+v[4], c.dir+v[4], 1, false)
		}
	}

	return c
}

func (c *Context) JsAfter() *Context {
	//替换外站的链接地址
	Pattern := regexp.MustCompile(`<script([^>]*?) src=([\"|\'])https?:/([^>]*?)>`)
	c.body = Pattern.ReplaceAllString(c.body, `<script${1} src=${2}${3}>`)

	return c
}

func (c *Context) ImgBefore() *Context {
	//优先删除链接中包含的本站域名，避免影响后续的路径判断
	regstr := fmt.Sprintf(`<img([^>]*?) src=([\"|\'])https?://%s([^>]*?)>`, webConfig.domain)
	Pattern := regexp.MustCompile(regstr)
	c.body = Pattern.ReplaceAllString(c.body, `<img${1} src=${2}${3}>`)

	return c
}

func (c *Context) ImgFind() *Context {
	Pattern := regexp.MustCompile(`<img([^>]*?) src=[\"|\']((https?:/)?([^>]*?))[\"|\']([^>]*?)>`)
	matcher := Pattern.FindAllStringSubmatch(c.body, -1)
	log.Log.Infof("Find link.img=%d, in %v\n", len(matcher), c.task)
	for i, v := range matcher {
		log.Log.Infof("  img.link[%d]: %s\n", i, v[2])
		if len(v[3]) > 0 { //含外站域名
			NewTasks(v[2], v[4], 1, false)
		} else if len(v[4]) == 0 { //空链接
			continue
		} else if v[4][0] == '/' { //绝对路径（本站）
			NewTasks(v[4], v[4], 1, false)
		} else { //相对路径，添加父页面的路径
			NewTasks(c.dir+v[4], c.dir+v[4], 1, false)
		}
	}

	return c
}

func (c *Context) ImgAfter() *Context {
	//替换外站的链接地址
	Pattern := regexp.MustCompile(`<img([^>]*?) src=([\"|\'])https?:/([^>]*?)>`)
	c.body = Pattern.ReplaceAllString(c.body, `<img${1} src=${2}${3}>`)

	return c
}

func (c *Context) ReplaceDomain() *Context {
	regstr := fmt.Sprintf(`<a ([^>]*?)href=[\"|\']https?://%s/?[\"|\']([^>]*?)>`, webConfig.domain)
	Pattern := regexp.MustCompile(regstr)
	c.body = Pattern.ReplaceAllString(c.body, `<a ${1}href="/index.html"${2}>`)

	regstr = fmt.Sprintf(`<a ([^>]*?)href=([\"|\'])https?://%s([^>]*?)>`, webConfig.domain)
	Pattern = regexp.MustCompile(regstr)
	c.body = Pattern.ReplaceAllString(c.body, `<a ${1}href=${2}${3}>`)

	return c
}

func (c *Context) DeleteBaidu() *Context {
	Pattern := regexp.MustCompile(`<[^<]*?hm.baidu.com/hm.js[^>]*?>`)
	c.body = Pattern.ReplaceAllString(c.body, ``)

	return c
}

func (c *Context) GetTitle() *Context {
	Pattern := regexp.MustCompile(`<title>([^<|>]*?)</title>`)
	matcher := Pattern.FindStringSubmatch(c.body)
	if matcher != nil {
		c.title = matcher[1]
	}

	return c
}

func (c *Context) GetKeywords() *Context {
	Pattern := regexp.MustCompile(`<meta([^>]*?) name=[\"|\']keywords[\"|\']([^>]*?) content=[\"|\']([^>]*?)[\"|\']([^>]*?)>`)
	matcher := Pattern.FindStringSubmatch(c.body)
	if matcher != nil {
		c.keywords = matcher[3]
	}

	return c
}

func (c *Context) GetDescription() *Context {
	Pattern := regexp.MustCompile(`<meta([^>]*?) name=[\"|\']description[\"|\']([^>]*) content=[\"|\']([^>]*?)[\"|\']([^>]*?)/>`)
	matcher := Pattern.FindStringSubmatch(c.body)
	if matcher != nil {
		c.title = matcher[3]
	}

	return c
}

func (c *Context) HrefFind() *Context {
	//优先删除链接中包含的本站域名，避免影响后续的路径判断 ReplaceDomain()
	if c.task.depth <= 1 {
		return c
	}

	Pattern := regexp.MustCompile(`<a ([^>]*?)href=[\"|\'](https?://)?([^>]*?)[\"|\']([^>]*?)>`)
	matcher := Pattern.FindAllStringSubmatch(c.body, -1)
	log.Log.Infof("Find link.href=%d, in %v\n", len(matcher), c.task)
	for i, v := range matcher {
		log.Log.Infof("  href.link[%d]: %s\n", i, v[3])
		Pattern = regexp.MustCompile(`[\s]+`) //url中可能包含空格
		v[3] = Pattern.ReplaceAllString(v[3], ``)
		if len(v[2]) > 0 { //含外站域名，忽略
			continue
		} else if len(v[3]) == 0 { //空连接
			continue
		} else if v[3][0] == '/'{ //绝对路径（本站）
			if len(v[3]) < len(webConfig.dir) || webConfig.dir != v[3][:len(webConfig.dir)] { //非入参指定路径
				continue
			}
			NewTasks(v[3], v[3], c.task.depth - 1, true)
		} else if strings.Contains(v[3], "javascript:") { //忽略js代码
			continue
		} else{ //相对路径，添加父页面的路径
			NewTasks(c.dir + v[3], c.dir + v[3], c.task.depth - 1, true)
		}
	}

	//保留外站的链接地址，不做处理

	return c
}

func (c *Context) CssUrlBefore() *Context {
	//优先删除链接中包含的本站域名，避免影响后续的路径判断
	regstr := fmt.Sprintf(`url\([\"|\']?https?://%s([^)]*?)[\"|\']?\)`, webConfig.domain)
	Pattern := regexp.MustCompile(regstr)
	c.body = Pattern.ReplaceAllString(c.body, `url\(${1}\)`)

	return c
}

func (c *Context) CssUrlFind() *Context {
	Pattern := regexp.MustCompile(`url\([\"|\']?((https?:/)?([^)]*?))[\"|\']?\)`)
	matcher := Pattern.FindAllStringSubmatch(c.body, -1)
	log.Log.Infof("Find css.url=%d, in %v\n", len(matcher), c.task)
	for i, v := range matcher {
		log.Log.Infof("  url.link[%d]: %s\n", i, v[1])
		if len(v[2]) > 0 { //含外站域名
			NewTasks(v[1], v[3], 1, false)
		} else if len(v[3]) == 0 { //空连接
			continue
		} else if v[3][0] == '/'{ //绝对路径（本站）
			NewTasks(v[3], v[3], 1, false)
		} else { //相对路径，添加父页面的路径
			NewTasks(c.dir + v[3], c.dir + v[3], 1, false)
		}
	}

	return c
}

func (c *Context) CssUrlAfter() *Context {
	//替换外站的链接地址
	Pattern := regexp.MustCompile(`url\([\"|\']?https?:/([^)]*?)[\"|\']?\)`)
	c.body = Pattern.ReplaceAllString(c.body, `url\(${1}\)`)

	return c
}

func LocalPageParse(path, mode string, body []byte) {
	task := &urlInfo{
		path: path,
		mode: mode,
		depth: 2,
	}

	NewContext(task, body).GetDir().HrefFind() //检测和下载href链接

	return
}
