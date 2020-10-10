package services

import (
	"sql-dog/src/domain/model"

	"github.com/pingcap/parser/ast"
)

type AnalyzerService interface {
	Parse(sql string) (ast.StmtNode, error)
	Extract(rootNode *ast.StmtNode, sql string) *model.Analyzer
}

type ReportService interface {
	Show(validator model.Validator)
}

type ValidateService interface {
	Validates(analyzers []model.Analyzer, validator model.Validator) []model.Report
	Validate(analyzer model.Analyzer, node *model.ValidatorNode, ignores []string) *model.Report
}
