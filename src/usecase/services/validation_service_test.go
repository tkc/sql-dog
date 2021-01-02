package services

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tkc/sql-dog/src/domain/model"
)

func TestValidators(t *testing.T) {
	cases := []struct {
		analyzer  *model.Analyzer
		validator model.ValidatorNode
		want      bool
	}{
		{
			analyzer: &model.Analyzer{
				Tables:         []model.Table{{Name: "table_a"}},
				Operations:     nil,
				StmtType:       model.StmtTypeSelect,
				NotNullColumns: []string{"deleted_at"},
			},
			validator: model.ValidatorNode{},
			want:      true,
		},
		{
			analyzer: &model.Analyzer{
				Tables:         []model.Table{{Name: "table_a"}},
				Operations:     nil,
				StmtType:       model.StmtTypeSelect,
				NotNullColumns: nil,
			},
			validator: model.ValidatorNode{
				TableName:       "table_a",
				Operations:      nil,
				StmtTypePattern: []model.StmtType{model.StmtTypeSelect},
				NotNullColumns: []string{

					"deleted_at",
				},
			},
			want: false,
		},
		{
			analyzer: &model.Analyzer{
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
			validator: model.ValidatorNode{
				TableName: "table_a",
				Operations: []string{
					"c_a",
					"c_b",
				},
				StmtTypePattern: []model.StmtType{model.StmtTypeSelect},
				NotNullColumns:  nil,
			},
			want: true,
		},
		{
			analyzer: &model.Analyzer{
				Tables: []model.Table{{Name: "table_a"}},
				Operations: []model.AnalyzerOperation{
					{
						Type:   model.OpTypeEq,
						Column: "c_a",
					},
				},
				StmtType:       model.StmtTypeSelect,
				NotNullColumns: nil,
			},
			validator: model.ValidatorNode{
				TableName: "table_a",
				Operations: []string{
					"c_a",
					"c_b",
				},
				StmtTypePattern: []model.StmtType{model.StmtTypeSelect},
				NotNullColumns:  nil,
			},
			want: false,
		},
	}

	for _, c := range cases {
		result := validate(c.analyzer, &c.validator, nil)
		assert.Equal(t, c.want, result == nil)
	}
}

// func TestNullValueValidators(t *testing.T) {
//	cases := []struct {
//		analyzer  *model.Analyzer
//		validator model.ValidatorNode
//		want      bool
//	}{
//		{
//			analyzer: &model.Analyzer{
//				Tables: []model.Table{
//					{
//						Name: "table_a",
//					},
//				},
//				Operations: []model.AnalyzerOperation{
//					{
//						Type:   model.OpTypeEq,
//						Column: "c_a",
//						Value:  int64(0), // nil value
//					},
//				},
//				StmtType:       model.StmtTypeSelect,
//				NotNullColumns: nil,
//			},
//			validator: model.ValidatorNode{
//				TableName: "table_a",
//				Operations: []string{
//					"c_a",
//				},
//				StmtTypePattern: []model.StmtType{model.StmtTypeSelect},
//				NotNullColumns:  nil,
//			},
//			want: false,
//		},
//		{
//			analyzer: &model.Analyzer{
//				Tables: []model.Table{
//					{
//						Name: "table_a",
//					},
//				},
//				Operations: []model.AnalyzerOperation{
//					{
//						Type:   model.OpTypeEq,
//						Column: "c_a",
//						Value:  "", // nil value
//					},
//				},
//				StmtType:       model.StmtTypeSelect,
//				NotNullColumns: nil,
//			},
//			validator: model.ValidatorNode{
//				TableName: "table_a",
//				Operations: []string{
//					"c_a",
//				},
//				StmtTypePattern: []model.StmtType{model.StmtTypeSelect},
//				NotNullColumns:  nil,
//			},
//			want: false,
//		},
//		{
//			analyzer: &model.Analyzer{
//				Tables: []model.Table{
//					{
//						Name: "table_a",
//					},
//				},
//				Operations: []model.AnalyzerOperation{
//					{
//						Type:   model.OpTypeIn,
//						Column: "c_a",
//						Value:  int64(0), // nil value
//					},
//				},
//				StmtType:       model.StmtTypeSelect,
//				NotNullColumns: nil,
//			},
//			validator: model.ValidatorNode{
//				TableName: "table_a",
//				Operations: []string{
//					"c_a",
//				},
//				StmtTypePattern: []model.StmtType{model.StmtTypeSelect},
//				NotNullColumns:  nil,
//			},
//			want: false,
//		},
//		{
//			analyzer: &model.Analyzer{
//				Tables: []model.Table{
//					{
//						Name: "table_a",
//					},
//				},
//				Operations: []model.AnalyzerOperation{
//					{
//						Type:   model.OpTypeIn,
//						Column: "c_a",
//						Value:  "", // nil value
//					},
//				},
//				StmtType:       model.StmtTypeSelect,
//				NotNullColumns: nil,
//			},
//			validator: model.ValidatorNode{
//				TableName: "table_a",
//				Operations: []string{
//					"c_a",
//				},
//				StmtTypePattern: []model.StmtType{model.StmtTypeSelect},
//				NotNullColumns:  nil,
//			},
//			want: false,
//		},
//	}
//
//	for _, c := range cases {
//		result := validate(c.analyzer, &c.validator, nil)
//		assert.Equal(t, c.want, result == nil)
//	}
// }
