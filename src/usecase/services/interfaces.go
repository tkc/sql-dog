package services

import (
	"github.com/pingcap/parser/ast"
	"sql-dog/src/domain/model"
)

type AnalyzerService interface {
	Parse(sql string) (ast.StmtNode, error)
	Extract(rootNode *ast.StmtNode, sql string) *model.Analyzer
}

type ReportService interface {
	Show(validator model.Validator)
}

type ValidatesService interface {
	Validates(analyzers []model.Analyzer, validator model.Validator) []model.Report
}
