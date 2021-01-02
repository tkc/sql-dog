package presenter

import "github.com/tkc/sql-dog/src/domain/model"

type ReportPresenter interface {
	Show(reports []model.Report)
}
