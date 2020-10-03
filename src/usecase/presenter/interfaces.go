package presenter

import "sql-dog/src/domain/model"

type ReportPresenter interface {
	Show(reports []model.Report)
}
