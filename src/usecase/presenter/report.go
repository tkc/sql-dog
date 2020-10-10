package presenter

import (
	"sql-dog/src/domain/model"

	"github.com/fatih/color"
	"github.com/kyokomi/emoji"
)

type reportPresenter struct{}

func NewReportPresenter() ReportPresenter {
	return reportPresenter{}
}

func (p reportPresenter) Show(reports []model.Report) {
	if len(reports) == 0 {
		beer := emoji.Sprint(":beer:")
		color.Magenta(beer + beer + beer + "No Report !!")
	}

	for _, report := range reports {
		color.Blue("----------Report----------")
		color.Cyan(report.ValidatorNode.TableName)
		for _, message := range report.ValidatorNode.Messages {
			color.Red(message)
		}

		color.White(report.Analyzer.SQL)

		if report.ValidatorNode.HasOperationsError() {
			color.Magenta("Operations Expect")
			for _, operation := range report.ValidatorNode.Operations {
				color.White(operation.Column)
			}
			color.Magenta("Operations Actual")
			for _, operation := range report.Analyzer.Operations {
				color.White(operation.Column)
			}
		}

		if report.ValidatorNode.HasInsertColumnsError() {
			color.Magenta("Inserts Expect")
			for _, insertColumn := range report.ValidatorNode.InsertColumns {
				color.White(insertColumn.Column)
			}
			color.Magenta("Inserts Actual")
			for _, insertColumn := range report.Analyzer.InsertColumns {
				color.White(insertColumn)
			}
		}

		if report.ValidatorNode.HasNotNullColumnsError() {
			color.Magenta("NotNullColumn Expect")
			for _, notNullColumn := range report.ValidatorNode.NotNullColumns {
				color.White(notNullColumn.Column)
			}
		}

		if len(report.ValidatorNode.NullValueOperation) > 0 {
			color.Magenta("NullValueOperation")
		}
	}
}
