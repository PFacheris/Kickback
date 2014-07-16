package db

import (
    "fmt"
    "github.com/jinzhu/gorm"
    _ "github.com/go-sql-driver/mysql"
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
    // Get database connection handle [*sql.DB](http://golang.org/pkg/database/sql/#DB)
    defer DB.DB().Close()

    DB.DB().SetMaxIdleConns(10)
    DB.DB().SetMaxOpenConns(100)

    if err != nil {
        panic(fmt.Sprintf("Got error when connect database, the error is '%v'", err))
    }
}
