package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tkc/sql-dog/src/domain/model"
)

func TestSelectQueries(t *testing.T) {
	cases := []struct {
		query  string
		expect []*model.Analyzer
	}{
		{
			query: `select * from table_a WHERE deleted_at IS NULL`,
			expect: []*model.Analyzer{{
				Tables: []model.Table{{
					Name: "table_a",
				}},
				Operations:     nil,
				StmtType:       model.StmtTypeSelect,
				NotNullColumns: []string{"deleted_at"},
			},
			},
		},
		{
			query: `select * from table_a where c_a = 1 and c_b in (1)`,
			expect: []*model.Analyzer{{
				Tables: []model.Table{{Name: "table_a"}},
				Operations: []model.AnalyzerOperation{
					{
						Type:   model.OpTypeEq,
						Column: "c_a",
						Value:  "1",
					},
					{
						Type:   model.OpTypeIn,
						Column: "c_b",
						Value:  "1",
					},
				},
				StmtType:       model.StmtTypeSelect,
				NotNullColumns: nil,
			},
			},
		},
		{
			query: `select * from table_a where deleted_at IS NULL AND ((c_a = 1) AND (c_b = 1) AND (c_c = 1)) ORDER BY id ASC LIMIT 1`,
			expect: []*model.Analyzer{{
				Tables: []model.Table{{Name: "table_a"}},
				Operations: []model.AnalyzerOperation{
					{
						Type:   model.OpTypeEq,
						Column: "c_a",
						Value:  "1",
					},
					{
						Type:   model.OpTypeEq,
						Column: "c_b",
						Value:  "1",
					},
					{
						Type:   model.OpTypeEq,
						Column: "c_c",
						Value:  "1",
					},
				},
				StmtType:       model.StmtTypeSelect,
				NotNullColumns: []string{"deleted_at"},
			}},
		},
		{
			query: `select * from table_a where c_a = 1 and c_b in (1) and c_c = 1`,
			expect: []*model.Analyzer{{
				Tables: []model.Table{{Name: "table_a"}},
				Operations: []model.AnalyzerOperation{
					{
						Type:   model.OpTypeEq,
						Column: "c_a",
						Value:  "1",
					},
					{
						Type:   model.OpTypeEq,
						Column: "c_c",
						Value:  "1",
					},
					{
						Type:   model.OpTypeIn,
						Column: "c_b",
						Value:  "1",
					},
				},
				StmtType:       model.StmtTypeSelect,
				NotNullColumns: nil,
			}},
		},
		{
			query: `select * from table_a as test WHERE deleted_at IS NULL`,
			expect: []*model.Analyzer{{
				Tables:         []model.Table{{Name: "table_a", AsName: "test"}},
				Operations:     nil,
				StmtType:       model.StmtTypeSelect,
				NotNullColumns: []string{"deleted_at"},
			}},
		},
		{
			// joinQuery
			query: `SELECT c_a FROM table_a as t_a
		INNER JOIN table_b as t_b ON table_a.join_c_a=table_b.join_c_b
		where table_a.c_a = 1 and table_b.c_b = 2 and deleted_at IS NULL`,
			expect: []*model.Analyzer{{
				Tables: []model.Table{
					{Name: "table_a", AsName: "t_a"},
					{Name: "table_b", AsName: "t_b"},
				},
				Operations: []model.AnalyzerOperation{
					{
						Type:   model.OpTypeEq,
						Column: "table_a.join_c_a",
						Value:  "table_b.join_c_b",
					},
					{
						Type:   model.OpTypeEq,
						Column: "table_a.c_a",
						Value:  "1",
					},
					{
						Type:   model.OpTypeEq,
						Column: "table_b.c_b",
						Value:  "2",
					},
				},
				StmtType:       model.StmtTypeSelect,
				NotNullColumns: []string{"deleted_at"},
			}},
		},
		{
			// joinQuery
			query: `SELECT c_a FROM table_a as t_a
		INNER JOIN table_b as t_b ON table_a.join_c_a=table_b.join_c_b
		where table_a.c_a = 'a' and table_b.c_b = 'b' and deleted_at IS NULL`,
			expect: []*model.Analyzer{{
				Tables: []model.Table{
					{Name: "table_a", AsName: "t_a"},
					{Name: "table_b", AsName: "t_b"},
				},
				Operations: []model.AnalyzerOperation{
					{
						Type:   model.OpTypeEq,
						Column: "table_a.join_c_a",
						Value:  "table_b.join_c_b",
					},
					{
						Type:   model.OpTypeEq,
						Column: "table_a.c_a",
						Value:  "a",
					},
					{
						Type:   model.OpTypeEq,
						Column: "table_b.c_b",
						Value:  "b",
					},
				},
				StmtType:       model.StmtTypeSelect,
				NotNullColumns: []string{"deleted_at"},
			}},
		},
		{
			// joinQuery
			query: `SELECT c_a FROM table_a as t_a
		INNER JOIN table_b as t_b ON table_a.join_c_a=table_b.join_c_b
		where table_a.c_a = 'a' and table_b.c_b = 'b' and deleted_at IS NULL`,
			expect: []*model.Analyzer{{
				Tables: []model.Table{
					{Name: "table_a", AsName: "t_a"},
					{Name: "table_b", AsName: "t_b"},
				},
				Operations: []model.AnalyzerOperation{
					{
						Type:   model.OpTypeEq,
						Column: "table_a.join_c_a",
						Value:  "table_b.join_c_b",
					},
					{
						Type:   model.OpTypeEq,
						Column: "table_a.c_a",
						Value:  "a",
					},
					{
						Type:   model.OpTypeEq,
						Column: "table_b.c_b",
						Value:  "b",
					},
				},
				StmtType:       model.StmtTypeSelect,
				NotNullColumns: []string{"deleted_at"},
			}},
		},
		{
			// subQuery
			query: `SELECT * FROM
		(SELECT id FROM t_a WHERE c_a = 'a' AND c_b = 'b') AS as_t_a
		INNER JOIN (SELECT id FROM t_b WHERE c_a = 'a' AND c_b = 'b' AND t_b.deleted_at IS NULL) AS t_c
		ON t_c.id = t_b.id`,
			expect: []*model.Analyzer{
				{
					Tables: []model.Table{
						{
							Name: "t_a",
						},
					},
					Operations: []model.AnalyzerOperation{
						{
							Type:   model.OpTypeEq,
							Column: "c_a",
							Value:  "a",
						},
						{
							Type:   model.OpTypeEq,
							Column: "c_b",
							Value:  "b",
						},
					},
					StmtType:       model.StmtTypeSelect,
					NotNullColumns: nil,
				},
				{
					Tables: []model.Table{{Name: "t_b"}},
					Operations: []model.AnalyzerOperation{
						{
							Type:   model.OpTypeEq,
							Column: "c_a",
							Value:  "a",
						},
						{
							Type:   model.OpTypeEq,
							Column: "c_b",
							Value:  "b",
						},
					},
					StmtType:       model.StmtTypeSelect,
					NotNullColumns: []string{"t_b.deleted_at"},
				}},
		},
		{
			// subQuery
			query: `SELECT * FROM t_a as as_t_a,
				(SELECT id from t_b where c_a = 'a') as as_t_b
		WHERE as_t_a.c_a = 'a' and as_t_a.c_b = sub_ci.c_b
		`,
			expect: []*model.Analyzer{
				{
					Tables: []model.Table{
						{
							Name:   "t_a",
							AsName: "as_t_a",
						},
						{
							Name:   "",
							AsName: "as_t_b",
						},
						{
							Name:   "t_b",
							AsName: "",
						},
					},
					Operations: []model.AnalyzerOperation{
						{
							Type:   model.OpTypeEq,
							Column: "c_a",
							Value:  "a",
						},
						{
							Type:   model.OpTypeEq,
							Column: "as_t_a.c_a",
							Value:  "a",
						},
						{
							Type:   model.OpTypeEq,
							Column: "as_t_a.c_b",
							Value:  "sub_ci.c_b",
						},
					},
					StmtType:       model.StmtTypeSelect,
					NotNullColumns: nil,
				},
				{
					Tables: []model.Table{
						{
							Name:   "t_b",
							AsName: "",
						},
					},
					Operations: []model.AnalyzerOperation{
						{
							Type:   model.OpTypeEq,
							Column: "c_a",
							Value:  "a",
						},
					},
					StmtType:       model.StmtTypeSelect,
					NotNullColumns: nil,
				},
			},
		},
		{
			// subQuery
			query: "SELECT * FROM `t_a` WHERE c_1 in ('v_1','v_2') AND (t_a.c_2 in (1) and t_a.c_3 in ('v_3'))",
			expect: []*model.Analyzer{
				{
					Tables: []model.Table{
						{
							Name: "t_a",
						},
					},
					Operations: []model.AnalyzerOperation{
						{
							Type:   model.OpTypeIn,
							Column: "c_1",
							Value:  "v_1,v_2",
						},
						{
							Type:   model.OpTypeIn,
							Column: "t_a.c_2",
							Value:  "1",
						},
						{
							Type:   model.OpTypeIn,
							Column: "t_a.c_3",
							Value:  "v_3",
						},
					},
					StmtType:       model.StmtTypeSelect,
					NotNullColumns: nil,
				}},
		},
	}

	analyzerService := NewAnalyzerService()

	for _, c := range cases {
		astNode, err := analyzerService.Parse(c.query)
		assert.Equal(t, err, nil)
		v := &StmtVisitor{}
		(astNode).Accept(v)
		assert.Equal(t, c.expect, v.Analyzers)
	}
}

func TestInsertQueries(t *testing.T) {
	cases := []struct {
		query  string
		expect []*model.Analyzer
	}{
		{
			query: "INSERT INTO `table_a` (`c_a`,`c_b`,`c_c`) VALUES (1,2,3)",
			expect: []*model.Analyzer{{
				Tables:        []model.Table{{Name: "table_a"}},
				Operations:    nil,
				StmtType:      model.StmtTypeInsert,
				InsertColumns: []string{"c_a", "c_b", "c_c"},
			}},
		},
	}

	analyzerService := NewAnalyzerService()

	for _, c := range cases {
		astNode, err := analyzerService.Parse(c.query)
		assert.Equal(t, err, nil)
		v := &StmtVisitor{}
		(astNode).Accept(v)
		assert.Equal(t, c.expect, v.Analyzers)
	}
}

func TestOperationValue(t *testing.T) {
	cases := []struct {
		query  string
		expect []*model.Analyzer
	}{
		{
			query: `select * from table_a WHERE c_a = '' and deleted_at IS NULL`,
			expect: []*model.Analyzer{{
				Tables: []model.Table{{
					Name: "table_a",
				}},
				Operations: []model.AnalyzerOperation{
					{
						Type:   model.OpTypeEq,
						Column: "c_a",
						Value:  "",
					},
				},
				StmtType:       model.StmtTypeSelect,
				NotNullColumns: []string{"deleted_at"},
			},
			},
		},
		{
			query: `select * from table_a WHERE c_a = 0 and deleted_at IS NULL`,
			expect: []*model.Analyzer{{
				Tables: []model.Table{{
					Name: "table_a",
				}},
				Operations: []model.AnalyzerOperation{
					{
						Type:   model.OpTypeEq,
						Column: "c_a",
						Value:  "0",
					},
				},
				StmtType:       model.StmtTypeSelect,
				NotNullColumns: []string{"deleted_at"},
			},
			},
		},
		{
			query: `select * from table_a WHERE c_a in (0) and deleted_at IS NULL`,
			expect: []*model.Analyzer{{
				Tables: []model.Table{{
					Name: "table_a",
				}},
				Operations: []model.AnalyzerOperation{
					{
						Type:   model.OpTypeIn,
						Column: "c_a",
						Value:  "0",
					},
				},
				StmtType:       model.StmtTypeSelect,
				NotNullColumns: []string{"deleted_at"},
			},
			},
		},
	}

	analyzerService := NewAnalyzerService()

	for _, c := range cases {
		astNode, err := analyzerService.Parse(c.query)
		assert.Equal(t, err, nil)
		v := &StmtVisitor{}
		(astNode).Accept(v)
		assert.Equal(t, c.expect, v.Analyzers)
	}
}
