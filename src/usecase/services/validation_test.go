package services

import (
	"sql-dog/src/domain/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidators(t *testing.T) {
	cases := []struct {
		analyzer  model.Analyzer
		validator model.ValidatorNode
		want      bool
	}{
		{
			analyzer: model.Analyzer{
				TableName:      "table_a",
				Operations:     nil,
				StmtType:       model.StmtTypeSelect,
				NotNullColumns: []string{"deleted_at"},
			},
			validator: model.ValidatorNode{},
			want:      true,
		},
		{
			analyzer: model.Analyzer{
				TableName:      "table_a",
				Operations:     nil,
				StmtType:       model.StmtTypeSelect,
				NotNullColumns: nil,
			},
			validator: model.ValidatorNode{
				TableName:       "table_a",
				Operations:      nil,
				StmtTypePattern: []model.StmtType{model.StmtTypeSelect},
				NotNullColumns: []model.ValidateColumn{
					{
						Column: "deleted_at",
					},
				},
			},
			want: false,
		},
		{
			analyzer: model.Analyzer{
				TableName: "table_a",
				Operations: []model.AnalyzerOperation{
					{
						Type:   model.OpTypeEq,
						Column: "c_a",
					},
					{
						Type:   model.OpTypeIn,
						Column: "c_b",
					},
				},
				StmtType:       model.StmtTypeSelect,
				NotNullColumns: nil,
			},
			validator: model.ValidatorNode{
				TableName: "table_a",
				Operations: []model.ValidateOperation{
					{
						Type:   model.OpTypeEq,
						Column: "c_a",
					},
					{
						Type:   model.OpTypeIn,
						Column: "c_b",
					},
				},
				StmtTypePattern: []model.StmtType{model.StmtTypeSelect},
				NotNullColumns:  nil,
			},
			want: true,
		},
		{
			analyzer: model.Analyzer{
				TableName: "table_a",
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
				Operations: []model.ValidateOperation{
					{
						Type:   model.OpTypeEq,
						Column: "c_a",
					},
					{
						Type:   model.OpTypeIn,
						Column: "c_b",
					},
				},
				StmtTypePattern: []model.StmtType{model.StmtTypeSelect},
				NotNullColumns:  nil,
			},
			want: false,
		},
	}

	v := NewValidatesService()
	for _, c := range cases {
		result := v.Validate(c.analyzer, &c.validator, nil)
		assert.Equal(t, c.want, result == nil)
	}
}
