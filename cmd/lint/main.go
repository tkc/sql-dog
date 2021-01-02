package main

import (
	"github.com/tkc/sql-dog/config"
	"github.com/tkc/sql-dog/src/infrastructure/datastore/mysql"
	"github.com/tkc/sql-dog/src/usecase/presenter"
	"github.com/tkc/sql-dog/src/usecase/services"
)

func main() {
	validation, err := config.ReadLintConfig("./linter.yaml")
	if err != nil {
		panic(err)
	}

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

	reportService := services.NewReportService(
		mysql.NewGeneralLogRepository(handler),
		services.NewAnalyzerService(),
		services.NewValidatesService(),
		presenter.NewReportPresenter(),
	)

	reportService.Show(*validation)
}
