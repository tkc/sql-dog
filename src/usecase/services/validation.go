package services

import (
	"sql-dog/src/domain/model"
)

const errorMessageNullColumn = "fail mast not null column check : "
const errorMessageInsertColumn = "fail required insert column check :"
const errorMessageOperation = "fail required where operations check : "

// const errorMessageNilValueOperation = "fail nil value check : "

type validateService struct{}

func NewValidatesService() ValidateService {
	return validateService{}
}

func (s validateService) Validates(analyzers []model.Analyzer, validator model.Validator) []model.Report {
	var reports []model.Report
	for _, analyzer := range analyzers {
		for _, node := range validator.Nodes {
			node := node
			res := s.Validate(analyzer, &node, validator.Ignores)
			if res != nil {
				reports = append(reports, *res)
			}
		}
	}

	//for _, analyzer := range analyzers {
	//	if len(analyzer.NullValueOperation) > 0 {
	//		reports = append(reports, model.Report{
	//			Analyzer: analyzer,
	//			ValidatorNode: &model.ValidatorNode{
	//				TableName:          analyzer.TableName,
	//				NullValueOperation: analyzer.NullValueOperation,
	//			},
	//		})
	//	}
	//}

	return reports
}

func (s validateService) Validate(analyzer model.Analyzer, node *model.ValidatorNode, ignores []string) *model.Report {
	if isIgnoreQuery(analyzer, node, ignores) {
		return nil
	}

	for i := range node.NotNullColumns {
		node.NotNullColumns[i].Valid = true
		valid := false
		for _, analyzerNotNullColumn := range analyzer.NotNullColumns {
			if analyzerNotNullColumn == node.NotNullColumns[i].Column {
				valid = true
			}
		}
		if !valid {
			node.NotNullColumns[i].Valid = false
			node.Messages = append(node.Messages, errorMessageNullColumn+node.NotNullColumns[i].Column)
		}
	}

	for i := range node.InsertColumns {
		node.InsertColumns[i].Valid = true
		valid := false
		for _, analyzerInsertColumns := range analyzer.InsertColumns {
			if analyzerInsertColumns == node.InsertColumns[i].Column {
				valid = true
			}
		}
		if !valid {
			node.InsertColumns[i].Valid = false
			node.Messages = append(node.Messages, errorMessageInsertColumn+node.InsertColumns[i].Column)
		}
	}

	for i := range node.Operations {
		node.Operations[i].Valid = true
		valid := false
		for _, analyzerOperation := range analyzer.Operations {
			if analyzerOperation.Column == node.Operations[i].Column {
				valid = true
			}
		}
		if !valid {
			node.Operations[i].Valid = false
			node.Messages = append(node.Messages, errorMessageOperation+node.Operations[i].Column)
		}
	}

	if node.HasError() {
		return &model.Report{
			Analyzer:      analyzer,
			ValidatorNode: node,
		}
	}
	return nil
}

func isIgnoreQuery(analyzer model.Analyzer, node *model.ValidatorNode, ignores []string) bool {
	for _, ignore := range ignores {
		if analyzer.SQL == ignore {
			return true
		}
	}
	if analyzer.TableName != node.TableName {
		return true
	}

	// Checking the Type of stmt
	stmtTypeMatch := false
	for _, stmtTypes := range node.StmtTypePattern {
		if analyzer.StmtType == stmtTypes {
			stmtTypeMatch = true
		}
	}

	return !stmtTypeMatch
}
