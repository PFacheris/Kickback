package db

import (
    "fmt"
    "github.com/jinzhu/gorm"
    _ "github.com/go-sql-driver/mysql"

    // Application Specific Imports
    . "github.com/pfacheris/kickback/models"
)

var DB gorm.DB
func init() {
    var err error

    /*
     * Load MYSQL Database
     */
    DB, err = gorm.Open("mysql", "root@/kickback_development?charset=utf8&parseTime=True")
    if err != nil {
      panic(fmt.Sprintf("Got error when connect database, the error is '%v'", err))
    }

    DB.DB().SetMaxIdleConns(10)
    DB.DB().SetMaxOpenConns(100)

    DB.AutoMigrate(User{})
}
