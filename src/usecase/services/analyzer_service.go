package services

import (
	"sql-dog/src/domain/model"
	"strings"

	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	"github.com/pingcap/parser/test_driver"
)

type analyzerService struct{}

func NewAnalyzerService() AnalyzerService {
	return analyzerService{}
}

func (a analyzerService) Parse(sql string) (ast.StmtNode, error) {
	p := parser.New()
	stmtNodes, _, err := p.Parse(sql, "", "")
	if err != nil {
		return nil, err
	}
	return stmtNodes[0], nil
}

func (a analyzerService) Extract(rootNode *ast.StmtNode, sql string) *model.Analyzer {
	v := &Visitor{}
	v.Analyzer.SQL = sql
	(*rootNode).Accept(v)
	return &v.Analyzer
}

type Visitor struct {
	Analyzer model.Analyzer
}

func (v *Visitor) checkNilValue(in ast.Node) {
	//if patternInExpr, ok := in.(*ast.PatternInExpr); ok {
	//	if valueExpr, ok := patternInExpr.List[0].(*test_driver.ValueExpr); ok {
	//		if valueExpr.Datum.GetValue() == nil {
	//			var nullValueOperation model.AnalyzerNullValueOperation
	//			nullValueOperation.TableName = v.Analyzer.TableName
	//			nullValueOperation.Value = valueExpr.Datum.GetValue()
	//			nullValueOperation.Type = model.OpTypeIn
	//			if columnNameExpr, ok := patternInExpr.Expr.(*ast.ColumnNameExpr); ok {
	//				nullValueOperation.Column = columnNameExpr.Name.String()
	//			}
	//			v.Analyzer.NullValueOperation = append(v.Analyzer.NullValueOperation, nullValueOperation)
	//		}
	//	}
	//}

	if binaryOperationExpr, ok := in.(*ast.BinaryOperationExpr); ok {
		if binaryOperationExpr.Op.String() == string(model.OpTypeEq) {
			if valueExpr, ok := binaryOperationExpr.R.(*test_driver.ValueExpr); ok {
				if valueExpr.Datum.GetValue() == nil {
					var nullValueOperation model.AnalyzerNullValueOperation
					nullValueOperation.TableName = v.Analyzer.TableName
					nullValueOperation.Value = valueExpr.Datum.GetValue()
					nullValueOperation.Type = model.OpTypeEq
					if columnNameExpr, ok := binaryOperationExpr.R.(*ast.ColumnNameExpr); ok {
						nullValueOperation.Column = columnNameExpr.Name.String()
					}
					v.Analyzer.NullValueOperation = append(v.Analyzer.NullValueOperation, nullValueOperation)
				}
			}
		}
	}
}

func (v *Visitor) tableName(in ast.Node) {
	if TableSource, ok := in.(*ast.TableSource); ok {
		if len(v.Analyzer.TableName) == 0 {
			if tableName, ok := TableSource.Source.(*ast.TableName); ok {
				if len(v.Analyzer.TableName) == 0 {
					v.Analyzer.TableName = tableName.Name.String()
				}
			}
			if len(TableSource.AsName.String()) > 0 {
				v.Analyzer.TableName = TableSource.AsName.String()
			}
		}
	}
}

func (v *Visitor) stmtType(in ast.Node) {
	if insertStmt, ok := in.(*ast.InsertStmt); ok {
		// TODO : sub query
		if len(v.Analyzer.StmtType) == 0 {
			v.Analyzer.StmtType = model.StmtTypeInsert
			for _, column := range insertStmt.Columns {
				v.Analyzer.InsertColumns = append(v.Analyzer.InsertColumns, column.Name.String())
			}
		}
	}
	if _, ok := in.(*ast.SelectStmt); ok {
		// TODO : sub query
		if len(v.Analyzer.StmtType) == 0 {
			v.Analyzer.StmtType = model.StmtTypeSelect
		}
	}
	if _, ok := in.(*ast.UpdateStmt); ok {
		// TODO : sub query
		if len(v.Analyzer.StmtType) == 0 {
			v.Analyzer.StmtType = model.StmtTypeUpdate
		}
	}
	if _, ok := in.(*ast.DeleteStmt); ok {
		// TODO : sub query
		if len(v.Analyzer.StmtType) == 0 {
			v.Analyzer.StmtType = model.StmtTypeDelete
		}
	}
}

func (v *Visitor) isNullExpr(in ast.Node) {
	if isNullExpr, ok := in.(*ast.IsNullExpr); ok {
		if columnNameExpr, ok := isNullExpr.Expr.(*ast.ColumnNameExpr); ok {
			v.Analyzer.NotNullColumns = append(
				v.Analyzer.NotNullColumns,
				formatColumnName(columnNameExpr.Name.String(), v.Analyzer.TableName))
		}
	}
}

func (v *Visitor) patternInExpr(in ast.Node) {
	if patternInExpr, ok := in.(*ast.PatternInExpr); ok {
		var operation model.AnalyzerOperation
		operation.Type = model.OpTypeIn
		if columnNameExpr, ok := patternInExpr.Expr.(*ast.ColumnNameExpr); ok {
			operation.Column = formatColumnName(columnNameExpr.Name.String(), v.Analyzer.TableName)
		}
		if valueExpr, ok := patternInExpr.List[0].(*test_driver.ValueExpr); ok {
			operation.Value = valueExpr.Datum.GetInt64()
		}
		//if columnNameExpr, ok := patternInExpr.List[0].(*ast.ColumnNameExpr); ok {
		//operation.Value = columnNameExpr..GetInt64()
		//}
		v.Analyzer.Operations = append(v.Analyzer.Operations, operation)
	}
}

func (v *Visitor) binaryOperationExpr(in ast.Node) {
	if binaryOperationExpr, ok := in.(*ast.BinaryOperationExpr); ok {
		if binaryOperationExpr.Op.String() == string(model.OpTypeEq) {
			var operation model.AnalyzerOperation
			operation.Type = model.OpType(binaryOperationExpr.Op.String())
			if columnNameExpr, ok := binaryOperationExpr.L.(*ast.ColumnNameExpr); ok {
				operation.Column = formatColumnName(columnNameExpr.Name.String(), v.Analyzer.TableName)
			}
			if valueExpr, ok := binaryOperationExpr.R.(*test_driver.ValueExpr); ok {
				operation.Value = valueExpr.Datum.GetValue()
			}
			if columnNameExpr, ok := binaryOperationExpr.R.(*ast.ColumnNameExpr); ok {
				operation.Value = columnNameExpr.Name.String()
			}
			v.Analyzer.Operations = append(v.Analyzer.Operations, operation)
		}
	}
}

func (v *Visitor) Enter(in ast.Node) (ast.Node, bool) {
	v.tableName(in)
	v.stmtType(in)
	v.isNullExpr(in)
	v.patternInExpr(in)
	v.binaryOperationExpr(in)
	v.checkNilValue(in)
	return in, false
}

func (v *Visitor) Leave(in ast.Node) (ast.Node, bool) {
	return in, true
}

func formatColumnName(column string, tableName string) string {
	slice := strings.Split(column, ".")
	if len(slice) == 1 {
		return column
	}
	if slice[0] == tableName {
		return slice[1]
	}
	return column
}

//func debug(in interface{}) {
//	log.Print("---debug---")
//	log.Print(reflect.TypeOf(in))
//}
