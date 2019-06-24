package db

import (
	"bytes"
	"compress/zlib"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"io"
)

func Open(v string) (err error) {
	DB, err = gorm.Open("sqlite3", dbConfig.Name)
	if err != nil {
		return
	}

	DB.AutoMigrate(&Page{}, &Command{})

	cmd := &Command{
		Options: v,
	}
	DB.Create(cmd)

	return
}

func Close() {
	DB.Close()
}

func AddPath(path, title, keywords, description string, mode string, parse bool, data []byte) (err error) {
	//zlib压缩
	var in bytes.Buffer
	w, _ := zlib.NewWriterLevel(&in, zlib.BestCompression)
	w.Write(data)
	w.Close()

	page := &Page{}
	err = DB.Select("id").Where("path = ?", path).First(page).Error
	if err == nil { //path已经存在，更新
		page.Path = path
		page.Title = title
		page.KeyWords = keywords
		page.Description = description
		page.Mode = mode
		page.Parse = parse
		page.Data = in.Bytes()

		Lock.Lock()
		err = DB.Save(page).Error
		Lock.Unlock()
	} else { //path，新增
		page.Path = path
		page.Title = title
		page.KeyWords = keywords
		page.Description = description
		page.Mode = mode
		page.Parse = parse
		page.Data = in.Bytes()

		Lock.Lock()
		err = DB.Create(page).Error
		Lock.Unlock()
	}

	return
}

func GetPath(path string) (context []byte, mode string, err error) {
	page := &Page{}
	err = DB.Where("path = ?", path).First(page).Error
	if err == nil {
		var out bytes.Buffer
		in := bytes.NewBuffer(page.Data)
		r, _ := zlib.NewReader(in)
		io.Copy(&out, r)
		context = out.Bytes()
		mode = page.Mode
	}

	return
}

func GetAllPath() (pages []*Page, err error) {
	err = DB.Select("path").Find(&pages).Error
	return
}

func GetCommands() (cmd []*Command, err error) {
	err = DB.Select("id, created_at, options").Order("id desc").Find(&cmd).Error
	return
}