package model

type Table struct {
	Name   string
	AsName string
}

type Analyzer struct {
	SQL            string
	Tables         []Table
	StmtType       StmtType
	InsertColumns  []string
	NotNullColumns []string
	Operations     []AnalyzerOperation
}

type AnalyzerOperation struct {
	Type   OpType
	Column string
	Value  interface{}
}

type Validator struct {
	Ignores []string
	Nodes   []ValidatorNode
}

type ValidatorNode struct {
	Messages        []string
	TableName       string
	StmtTypePattern []StmtType
	Operations      []string
	NotNullColumns  []string
	InsertColumns   []string
}

type ValidateOperation struct {
	Type   OpType
	Column string
}

type ValidateColumn struct {
	Column string
}

type Report struct {
	Analyzer      Analyzer
	ValidatorNode *ValidatorNode
}

type OpType string

const OpTypeEq OpType = "eq"
const OpTypeIn OpType = "in"
const OpTypeLike OpType = "like"

type StmtType string

const StmtTypeInsert StmtType = "insert"
const StmtTypeSelect StmtType = "select"
const StmtTypeUpdate StmtType = "update"
const StmtTypeDelete StmtType = "delete"

type GeneralLog struct {
	CommandType string
	Argument    string
}

type DatabaseDescResult struct {
	Field   string
	Type    string
	Null    string
	Key     string
	Default string
	Extra   string
}
