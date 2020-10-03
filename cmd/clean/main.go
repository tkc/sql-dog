package main

import (
	"sql-dog/src/infrastructure/datastore/mysql"
)

func main() {
	handler, _, _ := mysql.NewMySQLHandler(
		"root",
		"password",
		"localhost",
		3306)
	repo := mysql.NewGeneralLogRepository(handler)
	repo.Clear()
}
