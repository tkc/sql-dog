package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tkc/sql-dog/src/domain/model"
)

func TestReadLintConfig(t *testing.T) {
	v := model.Validator{
		Ignores: []string{
			"DELETE FROM table_1",
		},
		Nodes: []model.ValidatorNode{
			{
				TableName: "table_1",
				StmtTypePattern: []model.StmtType{
					model.StmtTypeSelect,
					model.StmtTypeInsert,
					model.StmtTypeUpdate,
					model.StmtTypeDelete,
				},
				Operations: []string{
					"c_1",
					"c_2",
				},
				NotNullColumns: []string{
					"deleted_at",
				},
			},
			{
				TableName: "table_2",
				StmtTypePattern: []model.StmtType{
					model.StmtTypeSelect,
				},
				Operations: []string{
					"c_1",
					"c_2",
				},
				NotNullColumns: []string{
					"deleted_at",
				},
			},
			{
				TableName: "table_3",
				StmtTypePattern: []model.StmtType{
					model.StmtTypeInsert,
				},
				InsertColumns: []string{
					"c_1",
					"c_2",
				},
			},
		},
	}

	validator, err := ReadLintConfig("./linter_test.yaml")
	assert.Equal(t, nil, err)
	assert.Equal(t, v, *validator)
}
