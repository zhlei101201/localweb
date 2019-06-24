package db

import (
	"github.com/jinzhu/gorm"
	"sync"
)

type config struct {
	Name		string
}

type Command struct {
	gorm.Model

	Options		string	`gorm:"type:VARCHAR(500)"`
}

type Page struct {
	gorm.Model

	Path		string	`gorm:"not null;index;type:VARCHAR(200)"`
	Title		string	`gorm:"index;type:VARCHAR(200)"`
	KeyWords	string	`gorm:"type:VARCHAR(200)"`
	Description	string	`gorm:"type:VARCHAR(500)"`
	Mode		string	`gorm:"index;type:VARCHAR(50)"`
	Parse		bool	`gorm:"index"`
	Data		[]byte	`gorm:"type:BLOB"`
}

var (
	dbConfig = &config{}

	DB *gorm.DB
	Lock sync.Mutex
)