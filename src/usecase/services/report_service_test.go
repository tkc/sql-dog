package services

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/tkc/sql-dog/src/domain/model"
	mock_mysql "github.com/tkc/sql-dog/src/infrastructure/datastore/mysql/mock"
	"github.com/tkc/sql-dog/src/usecase/presenter"
)

func TestReportService_CreateReport(t *testing.T) {
	cases := []struct {
		queries   []string
		validator model.Validator
		expect    *model.Report
	}{
		{
			queries: []string{
				"select * from table_a WHERE c_1 = 1",
			},
			validator: model.Validator{
				Ignores: nil,
				Nodes: []model.ValidatorNode{
					{
						TableName: "table_a",
						Operations: []string{
							"c_2",
						},
						StmtTypePattern: []model.StmtType{
							model.StmtTypeSelect,
						},
					},
				},
			},
			expect: &model.Report{
				ValidatorNode: &model.ValidatorNode{
					Operations: []string{
						"c_2",
					},
				},
			},
		},
		{
			queries: []string{
				"select * from table_a WHERE c_1 = 1",
			},
			validator: model.Validator{
				Ignores: nil,
				Nodes: []model.ValidatorNode{
					{
						TableName: "table_a",
						Operations: []string{
							"c_1",
						},
						StmtTypePattern: []model.StmtType{
							model.StmtTypeSelect,
						},
					},
				},
			},
			expect: nil,
		},
		{
			queries: []string{
				"SELECT c_a FROM table_a as t_a JOIN table_b as t_b ON table_a.join_c_a=table_b.join_c_b where table_a.c_a = 1 and table_b.c_b = 2 and `deleted_at` IS NULL",
			},
			validator: model.Validator{
				Ignores: nil,
				Nodes: []model.ValidatorNode{
					{
						TableName: "table_a",
						Operations: []string{
							"c_1_none",
						},
						StmtTypePattern: []model.StmtType{
							model.StmtTypeSelect,
						},
					},
				},
			},
			expect: &model.Report{
				ValidatorNode: &model.ValidatorNode{
					TableName: "table_a",
					Operations: []string{
						"c_1_none",
					},
				},
			},
		},
		{
			queries: []string{
				"SELECT c_a FROM table_a as t_a JOIN table_b as t_b ON table_a.join_c_a=table_b.join_c_b where table_a.c_a = 1 and table_b.c_b = 2 and `deleted_at` IS NULL",
			},
			validator: model.Validator{
				Ignores: nil,
				Nodes: []model.ValidatorNode{
					{
						TableName: "table_b",
						Operations: []string{
							"c_2_none",
						},
						StmtTypePattern: []model.StmtType{
							model.StmtTypeSelect,
						},
					},
				},
			},
			expect: &model.Report{
				ValidatorNode: &model.ValidatorNode{
					TableName: "table_b",
					Operations: []string{
						"c_2_none",
					},
				},
			},
		},
	}

	controller := gomock.NewController(t)
	mock_mysql.NewMockGeneralLogRepository(controller)

	reportService := NewReportService(
		mock_mysql.NewMockGeneralLogRepository(controller),
		NewAnalyzerService(),
		NewValidatesService(),
		presenter.NewReportPresenter())

	for _, c := range cases {
		reports := reportService.CreateReport(c.queries, c.validator)
		for _, r := range reports {
			assert.Equal(t, c.expect.ValidatorNode.Operations, r.ValidatorNode.Operations)
		}
	}

	controller.Finish()
}
