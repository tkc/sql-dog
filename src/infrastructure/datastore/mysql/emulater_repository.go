package mysql

import (
	"github.com/tkc/sql-dog/src/domain/model"
	"gorm.io/gorm"
)

const logTableName = "general_log"

type generalLogRepository struct {
	db *gorm.DB
}

//go:generate mockgen -destination mock/general_log_repository.go github.com/tkc/sql-dog/src/infrastructure/datastore/mysql GeneralLogRepository
type GeneralLogRepository interface {
	Clear() error
	GetQueries() ([]string, error)
}

func NewGeneralLogRepository(
	db *gorm.DB,
) GeneralLogRepository {
	return &generalLogRepository{
		db: db,
	}
}

func (r *generalLogRepository) Clear() error {
	r.db.Exec("truncate table general_log")
	return nil
}

func (r *generalLogRepository) Tables() ([]string, error) {
	var tables []string
	if err := r.db.Raw("SHOW TABLES").Scan(&tables).Error; err != nil {
		return nil, err
	}
	return tables, nil
}

func (r *generalLogRepository) TablesSchemas(tableName string) ([]model.DatabaseDescResult, error) {
	var results []model.DatabaseDescResult
	if err := r.db.Raw("desc ?", tableName).Scan(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

func (r *generalLogRepository) GetQueries() ([]string, error) {
	var logs []model.GeneralLog
	if err := r.db.
		Table(logTableName).
		Select("command_type, argument").
		Where("command_type in ('Execute', 'Query')").
		Find(&logs).Error; err != nil {
		return nil, err
	}
	var res []string
	for _, l := range logs {
		res = append(res, l.Argument)
	}
	return res, nil
}
