package mysql

import (
	"sql-dog/src/domain/model"

	"gorm.io/gorm"
)

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
	return r.db.Exec("truncate table general_log").Error
}

func (r *generalLogRepository) GetQueries() ([]string, error) {
	var logs []model.GeneralLog
	if err := r.db.
		Table("general_log").
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
