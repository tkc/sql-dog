package services

import (
	"strconv"
	"strings"

	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/test_driver"
	"github.com/tkc/sql-dog/src/domain/model"
)

type StmtVisitor struct {
	SQL       string
	Analyzers []*model.Analyzer
}

func (v *StmtVisitor) Enter(
	in ast.Node) (ast.Node, bool) {
	// InsertStmt
	if stmt, ok := in.(*ast.InsertStmt); ok {
		var analyzer model.Analyzer
		analyzer.SQL = v.SQL
		analyzer.StmtType = model.StmtTypeInsert

		for _, column := range stmt.Columns {
			analyzer.InsertColumns = append(analyzer.InsertColumns, column.Name.String())
		}

		tableNameVisitor := &TableNameVisitor{}
		(*stmt).Accept(tableNameVisitor)
		analyzer.Tables = tableNameVisitor.Tables

		operationExprVisitor := &OperationExprVisitor{}
		(*stmt).Accept(operationExprVisitor)
		analyzer.Operations = append(analyzer.Operations, operationExprVisitor.Operations...)

		patternInExprVisitor := &PatternInExprVisitor{}
		(*stmt).Accept(patternInExprVisitor)
		analyzer.Operations = append(analyzer.Operations, patternInExprVisitor.Operations...)

		isNullExprVisitor := &IsNullExprVisitor{}
		(*stmt).Accept(isNullExprVisitor)
		analyzer.NotNullColumns = isNullExprVisitor.NotNullColumns

		v.Analyzers = append(v.Analyzers, &analyzer)
	}

	// SelectStmt
	if stmt, ok := in.(*ast.SelectStmt); ok {
		var analyzer model.Analyzer
		analyzer.SQL = v.SQL
		analyzer.StmtType = model.StmtTypeSelect

		tableNameVisitor := &TableNameVisitor{}
		(*stmt).Accept(tableNameVisitor)

		// check hasOriginalTableName for subQuery
		if tableNameVisitor.hasOriginalTableName() {
			analyzer.Tables = tableNameVisitor.Tables

			operationExprVisitor := &OperationExprVisitor{}
			(*stmt).Accept(operationExprVisitor)
			analyzer.Operations = append(analyzer.Operations, operationExprVisitor.Operations...)

			patternInExprVisitor := &PatternInExprVisitor{}
			(*stmt).Accept(patternInExprVisitor)
			analyzer.Operations = append(analyzer.Operations, patternInExprVisitor.Operations...)

			isNullExprVisitor := &IsNullExprVisitor{}
			(*stmt).Accept(isNullExprVisitor)
			analyzer.NotNullColumns = isNullExprVisitor.NotNullColumns

			v.Analyzers = append(v.Analyzers, &analyzer)
		}
	}

	// UpdateStmt
	if stmt, ok := in.(*ast.UpdateStmt); ok {
		var analyzer model.Analyzer
		analyzer.SQL = v.SQL
		analyzer.StmtType = model.StmtTypeUpdate

		tableNameVisitor := &TableNameVisitor{}
		(*stmt).Accept(tableNameVisitor)
		analyzer.Tables = tableNameVisitor.Tables

		operationExprVisitor := &OperationExprVisitor{}
		(*stmt).Accept(operationExprVisitor)
		analyzer.Operations = append(analyzer.Operations, operationExprVisitor.Operations...)

		patternInExprVisitor := &PatternInExprVisitor{}
		(*stmt).Accept(patternInExprVisitor)
		analyzer.Operations = append(analyzer.Operations, patternInExprVisitor.Operations...)

		isNullExprVisitor := &IsNullExprVisitor{}
		(*stmt).Accept(isNullExprVisitor)
		analyzer.NotNullColumns = isNullExprVisitor.NotNullColumns

		v.Analyzers = append(v.Analyzers, &analyzer)
	}

	// DeleteStmt
	if stmt, ok := in.(*ast.DeleteStmt); ok {
		var analyzer model.Analyzer
		analyzer.SQL = v.SQL
		analyzer.StmtType = model.StmtTypeDelete

		tableNameVisitor := &TableNameVisitor{}
		(*stmt).Accept(tableNameVisitor)
		analyzer.Tables = tableNameVisitor.Tables

		operationExprVisitor := &OperationExprVisitor{}
		(*stmt).Accept(operationExprVisitor)
		analyzer.Operations = append(analyzer.Operations, operationExprVisitor.Operations...)

		patternInExprVisitor := &PatternInExprVisitor{}
		(*stmt).Accept(patternInExprVisitor)
		analyzer.Operations = append(analyzer.Operations, patternInExprVisitor.Operations...)

		isNullExprVisitor := &IsNullExprVisitor{}
		(*stmt).Accept(isNullExprVisitor)
		analyzer.NotNullColumns = isNullExprVisitor.NotNullColumns

		v.Analyzers = append(v.Analyzers, &analyzer)
	}
	return in, false
}

func (v *StmtVisitor) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

type TableNameVisitor struct {
	Tables []model.Table
}

func (v *TableNameVisitor) hasOriginalTableName() bool {
	return len(v.Tables) > 0 && len(v.Tables[0].Name) > 0
}

func (v *TableNameVisitor) Enter(in ast.Node) (ast.Node, bool) {
	if TableSource, ok := in.(*ast.TableSource); ok {
		var t model.Table
		if tableName, ok := TableSource.Source.(*ast.TableName); ok {
			t.Name = tableName.Name.String()
		}
		if len(TableSource.AsName.String()) > 0 {
			t.AsName = TableSource.AsName.String()
		}
		v.Tables = append(v.Tables, t)
	}
	return in, false
}

func (v *TableNameVisitor) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

type IsNullExprVisitor struct {
	NotNullColumns []string
}

func (v *IsNullExprVisitor) Enter(in ast.Node) (ast.Node, bool) {
	if isNullExpr, ok := in.(*ast.IsNullExpr); ok {
		if columnNameExpr, ok := isNullExpr.Expr.(*ast.ColumnNameExpr); ok {
			v.NotNullColumns = append(v.NotNullColumns, columnNameExpr.Name.String())
		}
	}
	return in, false
}

func (v *IsNullExprVisitor) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

type PatternInExprVisitor struct {
	Operations []model.AnalyzerOperation
}

func (v *PatternInExprVisitor) Enter(in ast.Node) (ast.Node, bool) {
	if patternInExpr, ok := in.(*ast.PatternInExpr); ok {
		var operation model.AnalyzerOperation
		operation.Type = model.OpTypeIn
		if columnNameExpr, ok := patternInExpr.Expr.(*ast.ColumnNameExpr); ok {
			operation.Column = columnNameExpr.Name.String()
		}
		operation.Value = v.getValuesAsString(patternInExpr.List)
		if columnNameExpr, ok := patternInExpr.List[0].(*ast.ColumnNameExpr); ok {
			operation.Value = columnNameExpr.Name
		}
		v.Operations = append(v.Operations, operation)
	}
	return in, false
}

func (v *PatternInExprVisitor) getValuesAsString(list []ast.ExprNode) interface{} {
	var res []string
	for _, l := range list {
		if valueExpr, ok := l.(*test_driver.ValueExpr); ok {
			if v, ok := valueExpr.Datum.GetValue().(string); ok {
				res = append(res, v)
			}
			if v, ok := valueExpr.Datum.GetValue().(int64); ok {
				res = append(res, strconv.FormatInt(v, 10))
			}
		}
	}
	return strings.Join(res, ",")
}

func (v *PatternInExprVisitor) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

type OperationExprVisitor struct {
	Operations []model.AnalyzerOperation
}

func (v *OperationExprVisitor) Enter(in ast.Node) (ast.Node, bool) {
	// BinaryOperationExpr
	if binaryOperationExpr, ok := in.(*ast.BinaryOperationExpr); ok {
		if binaryOperationExpr.Op.String() == string(model.OpTypeEq) {
			var operation model.AnalyzerOperation
			operation.Type = model.OpType(binaryOperationExpr.Op.String())
			if columnNameExpr, ok := binaryOperationExpr.L.(*ast.ColumnNameExpr); ok {
				operation.Column = columnNameExpr.Name.String()
			}
			if valueExpr, ok := binaryOperationExpr.R.(*test_driver.ValueExpr); ok {
				if v, ok := valueExpr.Datum.GetValue().(string); ok {
					operation.Value = v
				}
				if v, ok := valueExpr.Datum.GetValue().(int64); ok {
					operation.Value = strconv.FormatInt(v, 10)
				}
			} else if columnNameExpr, ok := binaryOperationExpr.R.(*ast.ColumnNameExpr); ok {
				operation.Value = columnNameExpr.Name.String()
			}
			v.Operations = append(v.Operations, operation)
		}
	}

	// PatternLikeExpr
	if patternLikeExpr, ok := in.(*ast.PatternLikeExpr); ok {
		var operation model.AnalyzerOperation
		operation.Type = model.OpTypeLike
		if columnNameExpr, ok := patternLikeExpr.Expr.(*ast.ColumnNameExpr); ok {
			operation.Column = columnNameExpr.Name.Name.String()
		}
		if valueExpr, ok := patternLikeExpr.Pattern.(*test_driver.ValueExpr); ok {
			operation.Value = valueExpr.GetString()
		}
		v.Operations = append(v.Operations, operation)
	}

	return in, false
}

func (v *OperationExprVisitor) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}
