package services

import (
	"log"
	"sql-dog/src/domain/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSelectQueries(t *testing.T) {
	cases := []struct {
		query    string
		analyzer model.Analyzer
	}{
		{
			query: "select * from `table_a` WHERE `deleted_at` IS NULL",
			analyzer: model.Analyzer{
				TableName:      "table_a",
				Operations:     nil,
				StmtType:       model.StmtTypeSelect,
				NotNullColumns: []string{"deleted_at"},
			},
		},
		{
			query: "select * from table_a where c_a in (1)",
			analyzer: model.Analyzer{
				TableName: "table_a",
				Operations: []model.AnalyzerOperation{
					{
						Type:   model.OpTypeIn,
						Column: "c_a",
						Value:  int64(1),
					},
				},
				StmtType:       model.StmtTypeSelect,
				NotNullColumns: nil,
			},
		},
		{
			query: "select * from table_a where c_a = 1 and c_b in (1)",
			analyzer: model.Analyzer{
				TableName: "table_a",
				Operations: []model.AnalyzerOperation{
					{
						Type:   model.OpTypeEq,
						Column: "c_a",
						Value:  int64(1),
					},
					{
						Type:   model.OpTypeIn,
						Column: "c_b",
						Value:  int64(1),
					},
				},
				StmtType:       model.StmtTypeSelect,
				NotNullColumns: nil,
			},
		},
		{
			query: "select * from `table_a` where `deleted_at` IS NULL AND ((c_a = 1) AND (c_b = 1) AND (c_c = 1)) ORDER BY `id` ASC LIMIT 1",
			analyzer: model.Analyzer{
				TableName: "table_a",
				Operations: []model.AnalyzerOperation{
					{
						Type:   model.OpTypeEq,
						Column: "c_a",
						Value:  int64(1),
					},
					{
						Type:   model.OpTypeEq,
						Column: "c_b",
						Value:  int64(1),
					},
					{
						Type:   model.OpTypeEq,
						Column: "c_c",
						Value:  int64(1),
					},
				},
				StmtType:       model.StmtTypeSelect,
				NotNullColumns: []string{"deleted_at"},
			},
		},
		{
			query: "select * from table_a where c_a = 1 and c_b in (1) and c_c = 1",
			analyzer: model.Analyzer{
				TableName: "table_a",
				Operations: []model.AnalyzerOperation{
					{
						Type:   model.OpTypeEq,
						Column: "c_a",
						Value:  int64(1),
					},
					{
						Type:   model.OpTypeIn,
						Column: "c_b",
						Value:  int64(1),
					},
					{
						Type:   model.OpTypeEq,
						Column: "c_c",
						Value:  int64(1),
					},
				},
				StmtType:       model.StmtTypeSelect,
				NotNullColumns: nil,
			},
		},
		{
			query: "select * from table_a as test WHERE `deleted_at` IS NULL",
			analyzer: model.Analyzer{
				TableName:      "table_a",
				Operations:     nil,
				StmtType:       model.StmtTypeSelect,
				NotNullColumns: []string{"deleted_at"},
			},
		},
	}

	analyzerService := NewAnalyzerService()

	for _, c := range cases {
		astNode, err := analyzerService.Parse(c.query)
		if err != nil {
			log.Print(err)
			return
		}
		v := &Visitor{}
		(astNode).Accept(v)

		analyzer := v.Analyzer
		assert.Equal(t, c.analyzer.TableName, analyzer.TableName)
		assert.Equal(t, c.analyzer.Operations, analyzer.Operations)
		assert.Equal(t, c.analyzer.StmtType, analyzer.StmtType)
		assert.Equal(t, c.analyzer.NotNullColumns, analyzer.NotNullColumns)
	}
}

func TestInsertQueries(t *testing.T) {
	cases := []struct {
		query    string
		analyzer model.Analyzer
	}{
		{
			query: "INSERT INTO `table_a` (`c_a`,`c_b`,`c_c`) VALUES (1,2,3)",
			analyzer: model.Analyzer{
				TableName:     "table_a",
				Operations:    nil,
				StmtType:      model.StmtTypeInsert,
				InsertColumns: []string{"c_a", "c_b", "c_c"},
			},
		},
	}

	analyzerService := NewAnalyzerService()

	for _, c := range cases {
		astNode, err := analyzerService.Parse(c.query)
		if err != nil {
			log.Print(err)
			return
		}
		v := &Visitor{}
		(astNode).Accept(v)

		analyzer := v.Analyzer
		assert.Equal(t, c.analyzer.TableName, analyzer.TableName)
		assert.Equal(t, c.analyzer.Operations, analyzer.Operations)
		assert.Equal(t, c.analyzer.StmtType, analyzer.StmtType)
		assert.Equal(t, c.analyzer.NotNullColumns, analyzer.NotNullColumns)
		assert.Equal(t, c.analyzer.InsertColumns, analyzer.InsertColumns)
	}
}
