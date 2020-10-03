package services

import (
	"log"
	"sql-dog/src/domain/model"
	"sql-dog/src/infrastructure/datastore/mysql"
	"sql-dog/src/usecase/presenter"
)

type reportService struct {
	generalLogRepository mysql.GeneralLogRepository
	analyzerService      AnalyzerService
	validatesService     ValidatesService
	reportPresenter      presenter.ReportPresenter
}

func NewReportService(
	generalLogRepository mysql.GeneralLogRepository,
	analyzerService AnalyzerService,
	validatesService ValidatesService,
	reportPresenter presenter.ReportPresenter) ReportService {
	return reportService{
		generalLogRepository,
		analyzerService,
		validatesService,
		reportPresenter,
	}
}

func (s reportService) Show(validator model.Validator) {
	res, _ := s.generalLogRepository.GetQueries()
	reportPresenter := presenter.NewReportPresenter()

	var analyzers []model.Analyzer
	for _, query := range res {
		astNode, err := s.analyzerService.Parse(query)
		if err != nil {
			log.Print(err)
		}
		analyzers = append(analyzers, *s.analyzerService.Extract(&astNode, query))
	}

	reportPresenter.Show(s.validatesService.Validates(analyzers, validator))
}
