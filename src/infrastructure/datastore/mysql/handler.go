package mysql

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func NewMySQLHandler(username, password, host string, port int) (*gorm.DB, func() error, error) {
	return newSQLHandler(
		username,
		password,
		host,
		port,
		"mysql",
	)
}

//func NewMySQLHandler() (*gorm.DB, func() error, error) {
//	return newSQLHandler(
//		"root",
//		"password",
//		"localhost",
//		3306,
//		"mysql",
//	)
//}

func newSQLHandler(userName, password, host string, port int, dbname string) (*gorm.DB, func() error, error) {
	if _, err := time.LoadLocation("UTC"); err != nil {
		return nil, nil, err
	}

	connectionString := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?parseTime=true&loc=Local",
		userName,
		password,
		host,
		port,
		dbname,
	)

	conn, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{})
	if err != nil {
		return nil, nil, err
	}

	db, err := conn.DB()
	if err != nil {
		return nil, nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, nil, err
	}

	conn.Set("gorm:table_options", "ENGINE=InnoDB")
	return conn, db.Close, nil
}
