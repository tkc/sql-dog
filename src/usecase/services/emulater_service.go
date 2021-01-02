package services

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"regexp"
	"strings"
	"time"

	"github.com/bxcodec/faker/v3"
	"github.com/tkc/sql-dog/src/domain/model"
	"github.com/tkc/sql-dog/src/infrastructure/datastore/mysql"
	"golang.org/x/sync/errgroup"
)

type emulateService struct {
	emulateRepository mysql.EmulateRepository
}

func NewEmulateService(emulateRepository mysql.EmulateRepository) EmulateService {
	return emulateService{emulateRepository}
}

func (s emulateService) Insert() {
	tableNames, _ := s.emulateRepository.Tables()
	for _, tableName := range tableNames {
		tableName := tableName
		descResult, _ := s.emulateRepository.TablesSchemas(tableName)

		eg := errgroup.Group{}
		for i := 0; i < 1; i++ {
			eg.Go(func() error {
				err := s.bulkInsert(tableName, descResult)
				if err != nil {
					return err
				}
				return nil
			})
		}

		if err := eg.Wait(); err != nil {
			panic(err)
		}
	}
}

func (s emulateService) bulkInsert(tableName string, schemes []model.DatabaseDescResult) error {
	valueStrings := []string{}
	args := []interface{}{}

	// TODO: 1000
	for i := 0; i < 1000; i++ {
		valueStrings = append(valueStrings, placeholder(len(schemes)-1))
		for _, v := range schemes {
			if v.Key != "PRI" {
				args = append(args, value(v.Type))
			}
		}
	}

	// panic: sql: expected 110 arguments, got 72
	// eg := errgroup.Group{}
	// for i := 0; i < 1000; i++ {
	//	eg.Go(func() error {
	//		valueStrings = append(valueStrings, placeholder(len(schemes)-1))
	//		for _, v := range schemes {
	//			if v.Key != "PRI" {
	//				args = append(args, value(v.Type))
	//			}
	//		}
	//		return nil
	//	})
	// }

	// if err := eg.Wait(); err != nil {
	//	panic(err)
	// }

	sql := fmt.Sprintf(
		"INSERT INTO %s %s VALUES %s",
		tableName,
		columns(schemes),
		strings.Join(valueStrings, ","))

	err := s.emulateRepository.Exec(sql, args...)
	if err != nil {
		log.Print(err)
		panic(err)
	}
	return nil
}

const varcharStr = `varchar`

var regexpVarcharStr, _ = regexp.Compile(varcharStr)

func IsStringFieldType(typeStr string) bool {
	if regexpVarcharStr.MatchString(typeStr) {
		return true
	}
	if typeStr == "text" {
		return true
	}
	return false
}

const intStr = `int`

var regexpIntStr, _ = regexp.Compile(intStr)

func IsIntFieldType(typeStr string) bool {
	return regexpIntStr.MatchString(typeStr)
}

const tinyIntStr = `tinyint`

var regexpTinyIntStr, _ = regexp.Compile(tinyIntStr)

func IsTinyIntFieldType(typeStr string) bool {
	return regexpTinyIntStr.MatchString(typeStr)
}

func value(fieldType string) interface{} {
	if IsStringFieldType(fieldType) {
		return faker.UUIDDigit()
	}
	if IsIntFieldType(fieldType) {
		if IsTinyIntFieldType(fieldType) {
			return 1
		}

		max := 1000000000
		n, err := rand.Int(rand.Reader, big.NewInt(int64(max)))
		if err != nil {
			panic(err)
		}
		return n
	}
	if fieldType == "datetime" {
		return time.Now()
	}
	return nil
}

func placeholder(len int) string {
	placeholder := "("
	for i := 0; i < len; i++ {
		if i < len-1 {
			placeholder += " ?,"
		} else {
			placeholder += " ?"
		}
	}
	return placeholder + ")"
}

func columns(schemes []model.DatabaseDescResult) string {
	len := len(schemes) - 1
	placeholder := "("
	for i, v := range schemes {
		if v.Key == "PRI" {
			continue
		}
		if i < len {
			placeholder += v.Field + ","
		} else {
			placeholder += v.Field
		}
	}
	return placeholder + ")"
}
