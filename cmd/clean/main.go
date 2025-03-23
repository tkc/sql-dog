package main

import (
	"context"
	"log"

	"github.com/tkc/sql-dog/config"
	"github.com/tkc/sql-dog/src/infrastructure/datastore/mysql"
)

func main() {
	conf, err := config.ReadConfig()
	if err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	handler, closer, err := mysql.NewMySQLHandler(
		conf.Username,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.RootDatabase)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	repo := mysql.NewGeneralLogRepository(handler)
	ctx := context.Background()
	if err := repo.Clear(ctx); err != nil {
		// defer呼び出しを確保するために、ここでクローズしてからエラー終了
		if closeErr := closer(); closeErr != nil {
			log.Printf("Warning: Failed to close database connection: %v", closeErr)
		}
		log.Fatalf("Failed to clear general log: %v", err)
	}

	// DBコネクションをクローズ
	if err := closer(); err != nil {
		log.Printf("Warning: Failed to close database connection: %v", err)
	}
	log.Println("Successfully cleared general log table")
}
