package services

import (
	"github.com/tkc/sql-dog/src/domain/model"
	"github.com/tkc/sql-dog/src/infrastructure/datastore/mysql"
	"github.com/tkc/sql-dog/src/usecase/presenter"
)

type reportService struct {
	generalLogRepository mysql.GeneralLogRepository
	analyzerService      AnalyzerService
	validatesService     ValidateService
	reportPresenter      presenter.ReportPresenter
}

func NewReportService(
	generalLogRepository mysql.GeneralLogRepository,
	analyzerService AnalyzerService,
	validatesService ValidateService,
	reportPresenter presenter.ReportPresenter) ReportService {
	return reportService{
		generalLogRepository,
		analyzerService,
		validatesService,
		reportPresenter,
	}
}

func (s reportService) Show(validator model.Validator) {
	queries, _ := s.generalLogRepository.GetQueries()
	reportPresenter := presenter.NewReportPresenter()
	reportPresenter.Show(s.CreateReport(queries, validator))
}

func (s reportService) CreateReport(queries []string, validator model.Validator) []model.Report {
	var analyzers = make([]*model.Analyzer, 0)
	for _, query := range queries {
		query := query
		astNode, err := s.analyzerService.Parse(query)
		if err != nil {
			panic(err)
		}
		analyzers = append(analyzers, s.analyzerService.Extract(&astNode, query)...)
	}
	return s.validatesService.Validates(analyzers, validator)
}
