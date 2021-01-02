package main

import (
	"github.com/tkc/sql-dog/config"
	"github.com/tkc/sql-dog/src/infrastructure/datastore/mysql"
	"github.com/tkc/sql-dog/src/usecase/services"
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
		conf.ServiceDatabase)

	repo := mysql.NewEmulateRepository(handler)
	emulateService := services.NewEmulateService(repo)
	emulateService.Insert()
}
