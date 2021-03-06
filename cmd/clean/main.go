package main

import (
	"github.com/tkc/sql-dog/config"
	"github.com/tkc/sql-dog/src/infrastructure/datastore/mysql"
)

func main() {
	conf, err := config.ReadConfig()
	if err != nil {
		panic(err)
	}

	handler, _, _ := mysql.NewMySQLHandler(
		conf.Username,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.RootDatabase)
	repo := mysql.NewGeneralLogRepository(handler)
	if err := repo.Clear(); err != nil {
		panic(err)
	}
}
