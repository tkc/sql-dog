package model

import (
	// test_driver
	_ "github.com/pingcap/parser/test_driver"
)

type Analyzer struct {
	SQL                string
	TableName          string
	StmtType           StmtType
	InsertColumns      []string
	NotNullColumns     []string
	Operations         []AnalyzerOperation
	NullValueOperation []AnalyzerNullValueOperation
}

type AnalyzerOperation struct {
	Type   OpType
	Column string
	Value  interface{}
}

type AnalyzerNullValueOperation struct {
	TableName string
	Type      OpType
	Column    string
	Value     interface{}
}

type Validator struct {
	Ignores []string
	Nodes   []ValidatorNode
}

type ValidatorNode struct {
	Messages           []string
	TableName          string
	StmtTypePattern    []StmtType
	Operations         []ValidateOperation
	NotNullColumns     []ValidateColumn
	InsertColumns      []ValidateColumn
	NullValueOperation []AnalyzerNullValueOperation
}

type ValidateOperation struct {
	Type   OpType
	Column string
	Valid  bool
}

type ValidateColumn struct {
	Column string
	Valid  bool
}

func (v ValidatorNode) HasError() bool {
	if v.HasOperationsError() {
		return true
	}
	if v.HasNotNullColumnsError() {
		return true
	}
	if v.HasInsertColumnsError() {
		return true
	}
	if v.HasNullValueOperationError() {
		return true
	}
	return false
}

func (v ValidatorNode) HasOperationsError() bool {
	for _, c := range v.Operations {
		if !c.Valid {
			return true
		}
	}
	return false
}

func (v ValidatorNode) HasNotNullColumnsError() bool {
	for _, c := range v.NotNullColumns {
		if !c.Valid {
			return true
		}
	}
	return false
}

func (v ValidatorNode) HasInsertColumnsError() bool {
	for _, c := range v.InsertColumns {
		if !c.Valid {
			return true
		}
	}
	return false
}

func (v ValidatorNode) HasNullValueOperationError() bool {
	return len(v.NullValueOperation) > 0
}

type Report struct {
	Analyzer      Analyzer
	ValidatorNode *ValidatorNode
}

type OpType string

const OpTypeEq OpType = "eq"
const OpTypeIn OpType = "in"

type StmtType string

const StmtTypeInsert StmtType = "insert"
const StmtTypeSelect StmtType = "select"
const StmtTypeUpdate StmtType = "update"
const StmtTypeDelete StmtType = "delete"

type GeneralLog struct {
	CommandType string
	Argument    string
}
