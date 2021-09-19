package mysql

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/tgallant/db2jsonschema/internal/schema"
)

type Driver struct {
	DataSource string
}

var (
	typesMap = map[string]*schema.FieldType{
		"bigint unsigned": {Name: "number", Format: ""},
		"longtext":        {Name: "string", Format: ""},
		"datetime(3)":     {Name: "string", Format: "date-time"},
		"tinyint(1)":      {Name: "boolean", Format: ""},
	}
)

func MapMySQLType(t string) (*schema.FieldType, error) {
	schemaType, exists := typesMap[t]
	if !exists {
		return &schema.FieldType{}, fmt.Errorf("Unknown data type: %s", t)
	}
	return schemaType, nil
}

func SelectTables(conn *sql.DB) ([]string, error) {
	row, err := conn.Query(`show tables`)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	var tables []string
	for row.Next() {
		var table string
		err = row.Scan(&table)
		if err != nil {
			return nil, err
		}
		tables = append(tables, table)
	}
	return tables, nil
}

func DescribeTable(conn *sql.DB, tableName string) (*schema.Table, error) {
	query := fmt.Sprintf("describe %s", tableName)
	row, err := conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer row.Close()
	var fields []*schema.Field
	for row.Next() {
		var name string
		var datatype string
		var nullable sql.NullString
		var key sql.NullString
		var defaultValue sql.NullString
		var extra sql.NullString
		err := row.Scan(
			&name,
			&datatype,
			&nullable,
			&key,
			&defaultValue,
			&extra,
		)
		if err != nil {
			return nil, err
		}
		fieldType, err := MapMySQLType(datatype)
		if err != nil {
			return nil, err
		}
		field := &schema.Field{
			Name: name,
			Type: fieldType,
		}
		fields = append(fields, field)
	}
	table := &schema.Table{
		Name:   tableName,
		Fields: fields,
	}
	return table, nil
}

func (d *Driver) ReadTables() ([]*schema.Table, error) {
	conn, err := sql.Open("mysql", d.DataSource)
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
		parsedTable, err := DescribeTable(conn, table)
		if err != nil {
			return nil, err
		}
		parsedTables = append(parsedTables, parsedTable)
	}
	return parsedTables, nil
}
