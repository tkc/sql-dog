package services

import (
	"github.com/pingcap/parser/ast"
	"github.com/tkc/sql-dog/src/domain/model"
)

type AnalyzerService interface {
	Parse(sql string) (ast.StmtNode, error)
	Extract(rootNode *ast.StmtNode, sql string) []*model.Analyzer
}

type EmulateService interface {
	Insert()
}

type ReportService interface {
	Show(validator model.Validator)
	CreateReport(queries []string, validator model.Validator) []model.Report
}

type ValidateService interface {
	Validates(analyzers []*model.Analyzer, validator model.Validator) []model.Report
}
