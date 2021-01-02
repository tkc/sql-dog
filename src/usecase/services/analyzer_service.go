package services

import (
	"github.com/pingcap/parser"
	"github.com/pingcap/parser/ast"
	"github.com/tkc/sql-dog/src/domain/model"
)

type analyzerService struct{}

func NewAnalyzerService() AnalyzerService {
	return analyzerService{}
}

func (a analyzerService) Parse(sql string) (ast.StmtNode, error) {
	p := parser.New()
	stmtNode, err := p.ParseOneStmt(sql, "", "")
	if err != nil {
		return nil, err
	}
	return stmtNode, nil
}

func (a analyzerService) Extract(rootNode *ast.StmtNode, sql string) []*model.Analyzer {
	v := &StmtVisitor{}
	v.SQL = sql
	(*rootNode).Accept(v)
	return v.Analyzers
}

// func debug(in interface{}) {
//	log.Print("---debug---")
//	log.Print(reflect.TypeOf(in))
// }
