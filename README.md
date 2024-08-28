![sql-dog](https://github.com/tkc/sql-dog/workflows/sql-dog/badge.svg?branch=master)
![reviewdog](https://github.com/tkc/sql-dog/workflows/reviewdog/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/tkc/sql-dog)](https://goreportcard.com/report/github.com/tkc/sql-dog)

# sql-dog
Analyzes SQL query logs and triggers a warning if a specified condition is missing from the query.

## MySQL Table Setting

This app will run an analysis on the log of the sql execution, go to mysql database settings and enable logging to the general_log table.

```sql
SET GLOBAL general_log = 'ON';
SET GLOBAL log_output='TABLE';

SET GLOBAL slow_query_log = 'ON';
SET GLOBAL long_query_time = 0;
```

## Database Setting
To set up database 
- rename `config.sample.yaml` to `config.yaml`
- modify the yaml settings

https://github.com/tkc/sql-dog/blob/master/config.sample.yaml

## Validate Setting
To set up validation
- rename `linter.sample.yaml` to `linter.yaml`
- modify the yaml settings

https://github.com/tkc/sql-dog/blob/master/linter.sample.yaml

## Run

run query analyzer and show report.

```bash
$ go run ./cmd/lint/main.go 
```

clear general_log table records.

```bash
$  go run ./cmd/clean/main.go 
```

### Features

- Check if a specific where condition set for a query to the target table.
- Check if NOT NUll constraint attached to the target table.

## Todo

- [ ] read other format query log. ex / http request, text log and other.

## Architecture

### Analyzer
Parsing the sql query and annotating it for parsing

### Validator
Compares the results of the analysis of the query with user-defined validations and outputs an analysis report.



