![Build Status](https://github.com/tkc/sql-dog/workflows/sql-dog/badge.svg)
![Reviewdog](https://github.com/tkc/sql-dog/workflows/reviewdog/badge.svg)
![CodeQL](https://github.com/tkc/sql-dog/workflows/CodeQL/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/tkc/sql-dog)](https://goreportcard.com/report/github.com/tkc/sql-dog)

# SQL Dog

SQL Dog is a tool that analyzes MySQL query logs and triggers warnings when specified conditions (such as WHERE clauses or NOT NULL constraints) are missing from queries. It's designed to enhance database security and performance.

## Table of Contents

- [Overview](#overview)
- [Design Philosophy and Architecture](#design-philosophy-and-architecture)
- [Installation](#installation)
- [Configuration](#configuration)
  - [MySQL Configuration](#mysql-configuration)
  - [Database Settings](#database-settings)
  - [Validation Rules](#validation-rules)
- [Usage](#usage)
- [Features](#features)
- [Technical Design](#technical-design)
- [Developer Information](#developer-information)
- [Troubleshooting](#troubleshooting)
- [Roadmap](#roadmap)

## Overview

SQL Dog helps with:

- **Security**: Detecting queries missing required WHERE conditions
- **Performance**: Identifying unoptimized queries that might cause full table scans
- **Quality Control**: Ensuring proper use of NOT NULL constraints
- **Auditing**: Analyzing database access patterns

It's a powerful tool for preventing bugs and issues in high-load production environments before they occur.

## Design Philosophy and Architecture

SQL Dog isn't just a query log analyzer—it acts as a watchdog for your database access, protecting both security and performance aspects. Below is a detailed explanation of the project's design philosophy.

### Core Problem and Solution Approach

**Problems Being Addressed**:
In many development environments, database access patterns degrade over time, leading to issues such as:

1. Queries missing essential WHERE clauses, causing full table scans
2. Queries retrieving deleted data by forgetting to check logical deletion flags (e.g., deleted_at)
3. Queries lacking proper permission conditions (tenant ID, user ID), creating security vulnerabilities

### Main Workflow

The overall program execution flow:

1. **Load Configuration**: Read database connection information and validation rules from YAML files
2. **Retrieve Query Logs**: Get executed queries from MySQL's general_log table
3. **Parse Queries**: Convert each query into an abstract syntax tree (AST) and extract structured information
4. **Validate Rules**: Compare extracted information with validation rules and detect violations
5. **Generate Reports**: Format and output validation results

```
+------------------+     +----------------+     +---------------+
| Load Configuration | --> | Retrieve Logs  | --> | Parse Queries |
+------------------+     +----------------+     +---------------+
                                                       |
                                                       v
                         +---------------+     +---------------+
                         | Output Reports | <-- | Validate Rules |
                         +---------------+     +---------------+
```

### Technology Choices

1. **PingCAP's SQL Parser**: Selected as a parser capable of accurately analyzing complex MySQL syntax
2. **GORM**: Chosen as a simple yet powerful ORM to simplify database access
3. **go-yaml**: Selected for parsing configuration files, considering readability and extensibility
4. **testify**: Adopted to write more concise and expressive test code

### Practical Usage Example

For example, in a multi-tenant application, you can set security rules like:

```yaml
tables:
  - name: users
    mustSelectColumns:
      - tenant_id # Tenant ID condition is required
    stmtTypePatterns:
      - select
      - update
      - delete
```

With this rule, the following query would trigger a warning:

```sql
-- Warning: no filtering by tenant_id
SELECT * FROM users WHERE name = 'John';
```

While this query would be allowed:

```sql
-- OK: filtered by tenant_id
SELECT * FROM users WHERE tenant_id = 123 AND name = 'John';
```

### Summary

SQL Dog is designed with the following philosophy:

1. **Defensive Programming**: Detect issues early to prevent production failures
2. **Configuration-Driven Approach**: Add and modify rules flexibly without changing code
3. **Domain-Driven Design**: View the technical domain of SQL queries from the business domain perspective of security and performance

By using this tool, you can maintain database access quality and prevent security and performance issues.

## Installation

### Prerequisites

- Go 1.18 or higher
- MySQL 5.7 or higher

### Installation Steps

```bash
# Clone the repository
git clone https://github.com/tkc/sql-dog.git
cd sql-dog

# Install dependencies
go mod download
```

## Configuration

### MySQL Configuration

Enable query logging in MySQL and configure it to record to the general_log table:

```sql
SET GLOBAL general_log = 'ON';
SET GLOBAL log_output = 'TABLE';

# Optional: Configure slow query log
SET GLOBAL slow_query_log = 'ON';
SET GLOBAL long_query_time = 0;
```

### Database Settings

1. Copy the sample configuration file:

   ```bash
   cp config.sample.yaml config.yaml
   ```

2. Edit `config.yaml` to configure your connection settings:
   ```yaml
   username: "root"
   password: "your_password"
   host: "localhost"
   port: 3306
   rootDatabase: "mysql"
   serviceDatabase: "your_database_name"
   ```

### Validation Rules

1. Copy the sample validation rules file:

   ```bash
   cp linter.sample.yaml linter.yaml
   ```

2. Edit `linter.yaml` to set up validation rules:

   ```yaml
   # Queries to exclude from validation
   ignores:
     - DELETE FROM temp_table

   # Tables and rules to validate
   tables:
     - name: users
       # Required column conditions for SELECT queries
       mustSelectColumns:
         - deleted_at
         - tenant_id
       # SQL statement types to target
       stmtTypePatterns:
         - select
         - update
         - delete
       # Columns requiring NOT NULL constraints
       notNullColumns:
         - deleted_at
   ```

## Usage

### Running Query Analysis

```bash
# Run query analysis and display report
go run ./cmd/lint/main.go

# Or use the compiled binary
./sql-dog-lint
```

### Clearing Log Table

```bash
# Clear records from the general_log table
go run ./cmd/clean/main.go

# Or use the compiled binary
./sql-dog-clean
```

## Features

- **WHERE Clause Checking**: Verifies that queries to specific tables include the required WHERE conditions
- **NOT NULL Constraint Checking**: Confirms that target tables have necessary NOT NULL constraints set
- **Multi-Table Support**: Configure different validation rules for multiple tables
- **Exclusion Rules**: Exempt specific queries from validation

## Developer Information

### Building

```bash
# Local build (outputs to ./bin/)
make build

# Install to $GOPATH/bin
make install
```

### Running Tests

```bash
# Run all tests
make test

# Run tests with race detection
make test-race

# Run tests with coverage reporting
make test-cover
```

### Formatting Code

```bash
# Format all Go code
make fmt
```

### Running Linter

```bash
# Run Go linter
make lint
```

## Troubleshooting

### Common Issues

1. **MySQL Connection Error**

   - Verify connection information in `config.yaml`
   - Check if MySQL server is running

2. **Empty general_log Table**

   - Verify MySQL logging is correctly enabled
   - Check if target queries have been executed

3. **Validation Rules Not Working**
   - Verify `linter.yaml` syntax is correct
   - Check that table and column names are accurately entered

## Roadmap

- [ ] Support for more query log formats (HTTP requests, text logs, etc.)
- [ ] Dashboard UI
- [ ] Real-time monitoring
- [ ] Automatic correction suggestions

## License

This project is released under the MIT License.
