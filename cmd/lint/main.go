package main

import (
	"log"

	"github.com/tkc/sql-dog/config"
	"github.com/tkc/sql-dog/src/infrastructure/datastore/mysql"
	"github.com/tkc/sql-dog/src/usecase/presenter"
	"github.com/tkc/sql-dog/src/usecase/services"
)

func main() {
	validation, err := config.ReadLintConfig("./linter.yaml")
	if err != nil {
		log.Fatalf("Failed to read lint config: %v", err)
	}

	conf, err := config.ReadConfig()
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	// データベース接続を確立
	handler, closer, err := mysql.NewMySQLHandler(
		conf.Username,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.RootDatabase)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// レポートサービスを初期化
	reportService := services.NewReportService(
		mysql.NewGeneralLogRepository(handler),
		services.NewAnalyzerService(),
		services.NewValidatesService(),
		presenter.NewReportPresenter(),
	)

	// レポートを表示
	if err := reportService.Show(*validation); err != nil {
		// defer呼び出しを確保するために、ここでクローズしてからエラー終了
		if closeErr := closer(); closeErr != nil {
			log.Printf("Warning: Failed to close database connection: %v", closeErr)
		}
		log.Fatalf("Failed to show report: %v", err)
	}

	// DBコネクションをクローズ
	if err := closer(); err != nil {
		log.Printf("Warning: Failed to close database connection: %v", err)
	}
}
