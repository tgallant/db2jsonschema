package sqlite3

import (
	"database/sql"
	"fmt"

	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
	"github.com/alecthomas/participle/v2/lexer/stateful"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tgallant/db2jsonschema/internal/db"
)

type Driver struct {
	DataSource string
}

type SQLiteTable struct {
	name string
	sql  string
}

type SQLiteCreateTable struct {
	TableName        string                   `"CREATE" "TABLE" @Ident`
	FieldExpressions []*SQLiteFieldExpression `( "(" @@ ( "," @@ )* ( "," )?`
	PrimaryKeys      []string                 `( "PRIMARY" "KEY" "(" @Ident ( "," @Ident )* ")" )* ( "," )? )?`
	Constraints      []*SQLiteConstraint      `( "CONSTRAINT" @@ ( "," "CONSTRAINT" @@ )* )? ")"`
}

type SQLiteFieldExpression struct {
	Name          string `@Ident`
	Type          string `@Ident ( "(" Number ")" )*`
	NotNull       bool   `( @"NOT_NULL"`
	AutoIncrement bool   `| @"AUTO_INCREMENT" )*`
}

type SQLiteConstraint struct {
	Name            string `@Ident`
	ForeignKey      string `"FOREIGN" "KEY" "(" @Ident ")"`
	ReferencedTable string `"REFERENCES" @Ident`
	ReferencedField string `"(" @Ident ")"`
}

var (
	typesMap = map[string]*db.FieldType{
		"int": &db.FieldType{
			Name:   "number",
			Format: "",
		},
		"integer": &db.FieldType{
			Name:   "number",
			Format: "",
		},
		"varchar": &db.FieldType{
			Name:   "string",
			Format: "",
		},
		"text": &db.FieldType{
			Name:   "string",
			Format: "",
		},
		"datetime": &db.FieldType{
			Name:   "string",
			Format: "date-time",
		},
	}

	sqlLexer = lexer.Must(stateful.NewSimple([]stateful.Rule{
		{`Keyword`, `(?i)\b(CREATE|TABLE|PRIMARY|FOREIGN|KEY|CONSTRAINT|REFERENCE)\b`, nil},
		{`Ident`, `[a-zA-Z_][a-zA-Z0-9_]*`, nil},
		{`Number`, `[-+]?\d*\.?\d+([eE][-+]?\d+)?`, nil},
		{`String`, `'[^']*'|"[^"]*"`, nil},
		{`Operators`, `<>|!=|<=|>=|[-+*/%,.()=<>]`, nil},
		{"whitespace", `\s+`, nil},
		{"backtick", "`", nil},
	}))

	parser = participle.MustBuild(
		&SQLiteCreateTable{},
		participle.Lexer(sqlLexer),
		participle.Unquote("String"),
	)
)

func MapSQLiteType(t string) (*db.FieldType, error) {
	schemaType, exists := typesMap[t]
	if !exists {
		return &db.FieldType{}, fmt.Errorf("Unknown data type: %s", t)
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
		row.Scan(&name, &sql)
		table := SQLiteTable{name, sql}
		tables = append(tables, &table)
	}
	return tables, nil
}

func ParseTableSQL(tableSQL string) (*db.Table, error) {
	createTable := &SQLiteCreateTable{}
	err := parser.ParseString("", tableSQL, createTable)
	if err != nil {
		return &db.Table{}, err
	}
	var fields []*db.Field
	for _, fieldExpression := range createTable.FieldExpressions {
		schemaType, err := MapSQLiteType(fieldExpression.Type)
		if err != nil {
			return &db.Table{}, err
		}
		field := &db.Field{
			Name: fieldExpression.Name,
			Type: schemaType,
		}
		fields = append(fields, field)
	}
	table := &db.Table{
		Name:        createTable.TableName,
		Fields:      fields,
		PrimaryKeys: createTable.PrimaryKeys,
	}
	return table, nil
}

func (d *Driver) ReadTables() ([]*db.Table, error) {
	conn, err := sql.Open("sqlite3", d.DataSource)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	tables, err := SelectTables(conn)
	if err != nil {
		return nil, err
	}
	var parsedTables []*db.Table
	for _, table := range tables {
		parsedTable, err := ParseTableSQL(table.sql)
		if err != nil {
			return nil, err
		}
		parsedTables = append(parsedTables, parsedTable)
	}
	return parsedTables, nil
}
