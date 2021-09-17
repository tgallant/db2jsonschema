package sqlite3

import (
	"database/sql"
	"fmt"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/alecthomas/participle/v2/lexer/stateful"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tgallant/db2jsonschema/internal/schema"
)

type Driver struct {
	DataSource string
}

type SQLiteTable struct {
	name string
	sql  string
}

type SQLiteCreateTable struct {
	TableName        string                   `parser:"'CREATE' 'TABLE' ( 'IF' 'NOT' 'EXISTS' )? @Ident"`
	FieldExpressions []*SQLiteFieldExpression `parser:"'(' @@ ( ',' @@ )* ( ',' )?"`
	PrimaryKeys      []string                 `parser:"( 'PRIMARY' 'KEY' '(' @Ident ( ',' @Ident )* ')' ( ',' )? )?"`
	ForeignKeys      []*SQLiteForeignKey      `parser:"( 'FOREIGN' 'KEY' @@ ( ',' 'FOREIGN' 'KEY' @@ )* ( ',' )? )?"`
	Constraints      []*SQLiteConstraint      `parser:"( 'CONSTRAINT' @@ ( ',' 'CONSTRAINT' @@ )* ( ',' )? )?"`
	Checks           []*SQLiteCheck           `parser:"( 'CHECK' '(' @@ ')' ( ',' 'CHECK' '(' @@ ')' )* )? ')'"`
}

type SQLiteFieldExpression struct {
	Name          string `parser:"@Ident"`
	Type          string `parser:"@Ident"`
	Limit         string `parser:"( '(' Number ')' )?"`
	NotNull       bool   `parser:"( @'NOT' 'NULL' | @'NOT_NULL'"`
	Default       string `parser:"| 'DEFAULT' '(' @Ident ')'"`
	AutoIncrement bool   `parser:"| @'AUTO_INCREMENT' )*"`
}

type SQLiteForeignKey struct {
	ForeignKey      string `parser:"'(' @Ident ')'"`
	ReferencedTable string `parser:"'REFERENCES' @Ident"`
	ReferencedField string `parser:"'(' @Ident ')'"`
}

type SQLiteCheck struct {
	Name   string `parser:"@Ident"`
	Values []int  `parser:"'IN' '(' @Number ',' @Number ')'"`
}

type SQLiteConstraint struct {
	Name            string `parser:"@Ident"`
	Key             string `parser:"('FOREIGN' | 'PRIMARY') 'KEY' '(' @Ident ')'"`
	ReferencedTable string `parser:"('REFERENCES' @Ident)?"`
	ReferencedField string `parser:"('(' @Ident ')')?"`
}

var (
	typesMap = map[string]*schema.FieldType{
		"int":      {Name: "number", Format: ""},
		"INTEGER":  {Name: "number", Format: ""},
		"integer":  {Name: "number", Format: ""},
		"varchar":  {Name: "string", Format: ""},
		"VARCHAR":  {Name: "string", Format: ""},
		"text":     {Name: "string", Format: ""},
		"boolean":  {Name: "boolean", Format: ""},
		"BOOLEAN":  {Name: "boolean", Format: ""},
		"datetime": {Name: "string", Format: "date-time"},
		"DATETIME": {Name: "string", Format: "date-time"},
	}

	sqlLexer = lexer.Must(stateful.NewSimple([]stateful.Rule{
		{Name: `Keyword`, Pattern: `(?i)\b(CREATE|TABLE|PRIMARY|FOREIGN|KEY|CONSTRAINT|REFERENCE|CHECK|IN)\b`, Action: nil},
		{Name: `Ident`, Pattern: `[a-zA-Z_][a-zA-Z0-9_]*`, Action: nil},
		{Name: `Number`, Pattern: `[-+]?\d*\.?\d+([eE][-+]?\d+)?`, Action: nil},
		{Name: `String`, Pattern: `'[^']*'`, Action: nil},
		{Name: `Operators`, Pattern: `<>|!=|<=|>=|[-+*/%,.()=<>]`, Action: nil},
		{Name: "whitespace", Pattern: `\s+`, Action: nil},
		{Name: "backtick", Pattern: "`", Action: nil},
		{Name: "quote", Pattern: `"`, Action: nil},
	}))

	parser = participle.MustBuild(
		&SQLiteCreateTable{},
		participle.Lexer(sqlLexer),
		participle.Unquote("String"),
	)
)

func MapSQLiteType(t string) (*schema.FieldType, error) {
	schemaType, exists := typesMap[t]
	if !exists {
		return &schema.FieldType{}, fmt.Errorf("Unknown data type: %s", t)
	}
	return schemaType, nil
}

func SelectTables(conn *sql.DB) ([]*SQLiteTable, error) {
	row, err := conn.Query(`select name, sql from sqlite_master where type = "table"`)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	var tables []*SQLiteTable
	for row.Next() {
		var name string
		var sql string
		err = row.Scan(&name, &sql)
		if err != nil {
			return nil, err
		}
		table := SQLiteTable{name, sql}
		tables = append(tables, &table)
	}
	return tables, nil
}

func ParseTableSQL(tableSQL string) (*schema.Table, error) {
	createTable := &SQLiteCreateTable{}
	err := parser.ParseString("", tableSQL, createTable)
	if err != nil {
		return &schema.Table{}, err
	}
	var fields []*schema.Field
	for _, fieldExpression := range createTable.FieldExpressions {
		schemaType, err := MapSQLiteType(fieldExpression.Type)
		if err != nil {
			return &schema.Table{}, err
		}
		field := &schema.Field{
			Name: fieldExpression.Name,
			Type: schemaType,
		}
		fields = append(fields, field)
	}
	table := &schema.Table{
		Name:        createTable.TableName,
		Fields:      fields,
		PrimaryKeys: createTable.PrimaryKeys,
	}
	return table, nil
}

func (d *Driver) ReadTables() ([]*schema.Table, error) {
	conn, err := sql.Open("sqlite3", d.DataSource)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	tables, err := SelectTables(conn)
	if err != nil {
		return nil, err
	}
	var parsedTables []*schema.Table
	for _, table := range tables {
		parsedTable, err := ParseTableSQL(table.sql)
		if err != nil {
			return nil, err
		}
		parsedTables = append(parsedTables, parsedTable)
	}
	return parsedTables, nil
}
