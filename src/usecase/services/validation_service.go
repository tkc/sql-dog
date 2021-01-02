package services

import (
	"strings"

	"github.com/tkc/sql-dog/src/domain/model"
)

const errorMessageNullColumn = "fail null column check : "
const errorMessageInsertColumn = "fail insert column check :"
const errorMessageOperation = "fail operations check : "

type validateService struct{}

func NewValidatesService() ValidateService {
	return validateService{}
}

func (v validateService) Validates(analyzers []*model.Analyzer, validator model.Validator) []model.Report {
	var reports []model.Report
	for _, analyzer := range analyzers {
		for _, node := range validator.Nodes {
			node := node
			res := validate(analyzer, &node, validator.Ignores)
			if res != nil {
				reports = append(reports, *res)
			}
		}
	}
	return reports
}

func isIgnoreTable(analyzer *model.Analyzer, node *model.ValidatorNode, ignores []string) bool {
	for _, ignore := range ignores {
		if analyzer.SQL == ignore {
			return true
		}
	}
	if isAllIgnoreTables(analyzer.Tables, node.TableName) {
		return true
	}
	stmtTypeMatch := false
	for _, stmtTypes := range node.StmtTypePattern {
		if analyzer.StmtType == stmtTypes {
			stmtTypeMatch = true
		}
	}
	return !stmtTypeMatch
}

func isAllIgnoreTables(tables []model.Table, tableName string) bool {
	var ignore = true
	for _, table := range tables {
		if table.Name == tableName {
			ignore = false
		}
	}
	return ignore
}

func validate(analyzer *model.Analyzer, validators *model.ValidatorNode, ignores []string) *model.Report {
	if isIgnoreTable(analyzer, validators, ignores) {
		return nil
	}

	for i := range validators.NotNullColumns {
		valid := false
		for _, analyzerNotNullColumn := range analyzer.NotNullColumns {
			if isTargetColumn(
				validators.TableName,
				validators.NotNullColumns[i],
				analyzerNotNullColumn,
				analyzer.Tables) {
				valid = true
			}
		}
		if !valid {
			validators.Messages = append(validators.Messages, errorMessageNullColumn+validators.NotNullColumns[i])
		}
	}

	for i := range validators.InsertColumns {
		valid := false
		for _, analyzerInsertColumns := range analyzer.InsertColumns {
			if isTargetColumn(
				validators.TableName,
				validators.InsertColumns[i],
				analyzerInsertColumns,
				analyzer.Tables) {
				valid = true
			}
		}
		if !valid {
			validators.Messages = append(validators.Messages, errorMessageInsertColumn+validators.InsertColumns[i])
		}
	}

	for i, validatorOperation := range validators.Operations {
		valid := false
		for _, analyzerOperation := range analyzer.Operations {
			if isTargetColumn(
				validators.TableName,
				validatorOperation,
				analyzerOperation.Column,
				analyzer.Tables) {
				valid = true

				// NullValue check
				// if ok, message := isNullValue(analyzerOperation); ok {
				//	validators.Messages = append(
				//		validators.Messages,
				//		errorMessageNilValueOperation+analyzerOperation.Column+": "+message)
				//}
			}
		}
		if !valid {
			validators.Messages = append(validators.Messages, errorMessageOperation+validators.Operations[i])
		}
	}

	if len(validators.Messages) > 0 {
		return &model.Report{
			Analyzer:      *analyzer,
			ValidatorNode: validators,
		}
	}
	return nil
}

func isTargetColumn(
	expectTableName string,
	expectColumnName string,
	actualColumnName string,
	actualTables []model.Table) bool {
	if expectColumnName == actualColumnName {
		return true
	}
	if expectTableName+","+expectColumnName == actualColumnName {
		return true
	}
	for _, table := range actualTables {
		if table.Name == expectTableName {
			slice := strings.Split(actualColumnName, ".")
			if len(slice) > 1 {
				if expectColumnName == slice[1] {
					return true
				}
			}
		}
	}
	return false
}

// func isNullValue(operation model.AnalyzerOperation) (bool, string) {
//	if operation.Value == nil {
//		return true, "nil"
//	}
//
//	if operation.Type == model.OpTypeEq {
//		if value, ok := operation.Value.(string); ok {
//			if len(value) == 0 {
//				return true, "empty"
//			}
//			if value == "0" {
//				return true, "0"
//			}
//		}
//		if value, ok := operation.Value.(int64); ok {
//			if value == 0 {
//				return true, "0"
//			}
//		}
//	}
//
//	if operation.Type == model.OpTypeIn {
//		if value, ok := operation.Value.(string); ok {
//			if len(value) == 0 {
//				return true, "empty"
//			}
//		}
//		if value, ok := operation.Value.(int64); ok {
//			if value == 0 {
//				return true, "0"
//			}
//			// TODO : 値に0を含む場合む条件分岐
//		}
//	}
//
//	return false, "not nil"
// }
