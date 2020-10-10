package main

import (
	"sql-dog/src/domain/model"
	"sql-dog/src/infrastructure/datastore/mysql"
	"sql-dog/src/usecase/presenter"
	"sql-dog/src/usecase/services"
)

func createValidation() model.Validator {
	return model.Validator{
		Ignores: []string{
			"DELETE FROM Ignores table_name",
		},
		Nodes: []model.ValidatorNode{
			{
				TableName: "table_name",
				Operations: []model.ValidateOperation{
					{
						Type:   model.OpTypeEq,
						Column: "require_column_a",
					},
					{
						Type:   model.OpTypeEq,
						Column: "require_column_b",
					},
				},
				StmtTypePattern: []model.StmtType{
					model.StmtTypeSelect,
					model.StmtTypeDelete,
				},
			},
		},
	}
}

func main() {
	v := createValidation()

	handler, _, _ := mysql.NewMySQLHandler(
		"root",
		"password",
		"localhost",
		3306)

	reportService := services.NewReportService(
		mysql.NewGeneralLogRepository(handler),
		services.NewAnalyzerService(),
		services.NewValidatesService(),
		presenter.NewReportPresenter(),
	)

	reportService.Show(v)
}
