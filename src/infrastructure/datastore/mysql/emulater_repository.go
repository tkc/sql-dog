package mysql

import (
	"github.com/tkc/sql-dog/src/domain/model"
	"gorm.io/gorm"
)

type emulateRepository struct {
	db *gorm.DB
}

type EmulateRepository interface {
	Tables() ([]string, error)
	TablesSchemas(tableName string) ([]model.DatabaseDescResult, error)
	Exec(sql string, arg ...interface{}) error
}

func NewEmulateRepository(
	db *gorm.DB,
) EmulateRepository {
	return &emulateRepository{
		db: db,
	}
}

func (r *emulateRepository) Tables() ([]string, error) {
	var tables []string
	if err := r.db.Raw("SHOW TABLES").Scan(&tables).Error; err != nil {
		return nil, err
	}
	return tables, nil
}

func (r *emulateRepository) TablesSchemas(tableName string) ([]model.DatabaseDescResult, error) {
	var results []model.DatabaseDescResult
	if err := r.db.Raw("DESC ?", tableName).Scan(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

func (r *emulateRepository) Exec(sql string, values ...interface{}) error {
	if err := r.db.Exec(sql, values...).Error; err != nil {
		return err
	}
	return nil
}
