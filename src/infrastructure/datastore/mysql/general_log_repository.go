package mysql

import (
	"gorm.io/gorm"
	"sql-dog/src/domain/model"
)

const logTableName = "general_log"

type generalLogRepository struct {
	db *gorm.DB
}

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
