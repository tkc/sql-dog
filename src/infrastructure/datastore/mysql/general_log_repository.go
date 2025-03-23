package mysql

import (
	"context"

	"github.com/tkc/sql-dog/src/domain/model"
	"gorm.io/gorm"
)

const logTableName = "general_log"

type generalLogRepository struct {
	db *gorm.DB
}

//go:generate mockgen -destination mock/general_log_repository.go github.com/tkc/sql-dog/src/infrastructure/datastore/mysql GeneralLogRepository
type GeneralLogRepository interface {
	Clear(ctx context.Context) error
	GetQueries(ctx context.Context) ([]string, error)
}

func NewGeneralLogRepository(
	db *gorm.DB,
) GeneralLogRepository {
	return &generalLogRepository{
		db: db,
	}
}

func (r *generalLogRepository) Clear(ctx context.Context) error {
	if err := r.db.WithContext(ctx).Exec("TRUNCATE TABLE general_log").Error; err != nil {
		return err
	}
	return nil
}

func (r *generalLogRepository) GetQueries(ctx context.Context) ([]string, error) {
	var logs []model.GeneralLog
	if err := r.db.WithContext(ctx).
		Table(logTableName).
		Select("command_type, argument").
		Where("command_type IN (?, ?)", "Execute", "Query").
		Find(&logs).Error; err != nil {
		return nil, err
	}

	// より効率的なメモリ使用のために事前にキャパシティを指定
	queries := make([]string, 0, len(logs))
	for _, l := range logs {
		queries = append(queries, l.Argument)
	}
	return queries, nil
}
