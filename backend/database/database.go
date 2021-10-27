package database

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB

func InitDb() {
	conn, err := gorm.Open("mysql", "sistema:Covid2020@488@tcp(localhost:3306)/vivircarlospaz?charset=utf8&parseTime=true&loc=Local")

	if err != nil {
		panic("Falló conexión con base de datos")
	}

	db = conn

	db.Exec("set names utf8")

	db.LogMode(true)
	fmt.Println("conn")
}

func GetDb() *gorm.DB {
	return db
}
