package config

import (
	"io/ioutil"

	"github.com/goccy/go-yaml"
	"github.com/tkc/sql-dog/src/domain/model"
)

type LintConfig struct {
	Ignores []string `yaml:"ignores"`
	Tables  []Table  `yaml:"tables"`
}

type Table struct {
	Name           string           `yaml:"name"`
	StmtTypes      []model.StmtType `yaml:"stmtTypePatterns"`
	Operations     []string         `yaml:"mustSelectColumns"`
	InsertColumns  []string         `yaml:"mustInsertColumns"`
	NotNullColumns []string         `yaml:"notNullColumns"`
}

func (c *LintConfig) ConvertToValidators() (*model.Validator, error) {
	var nodes []model.ValidatorNode

	for _, t := range c.Tables {
		nodes = append(nodes,
			model.ValidatorNode{
				TableName:       t.Name,
				StmtTypePattern: t.StmtTypes,
				Operations:      t.Operations,
				InsertColumns:   t.InsertColumns,
				NotNullColumns:  t.NotNullColumns,
			})
	}

	return &model.Validator{
		Ignores: c.Ignores,
		Nodes:   nodes,
	}, nil
}

func ReadLintConfig(path string) (*model.Validator, error) {
	buf, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var config LintConfig
	if err := yaml.Unmarshal(buf, &config); err != nil {
		panic(err)
	}
	return config.ConvertToValidators()
}
