![sql-dog](https://github.com/tkc/sql-dog/workflows/sql-dog/badge.svg?branch=master)

# sql-dog

Parse the query log in sql and output a report with specific conditions.

## MySQL Table Setting

```sql
SET GLOBAL general_log = 'ON';
SET GLOBAL log_output='TABLE';
```

## Run

run query analyzer and show report.

```bash
$ go run ./cmd/lint
```

clear general_log table records.

```bash
$  go run ./cmd/clean 
```

### Features

- Check if a specific where condition set for a query to the target table.
- Check if the value of the query specified as null to the target table.
- Check if NOT NUll constraint attached to the target table.

## Todo

- [ ] read sql setting from config fail. it's hard-coded now.
- [ ] read yaml config and convert validate struct.
- [ ] read other format query log. ex / http request, text log and other.
- [ ] separate analyzer service depend on report type.


## Architecture

### Analyzer
Parsing the sql query and annotating it for parsing

### Validator
Compares the results of the analysis of the query with user-defined validations and outputs an analysis report.



