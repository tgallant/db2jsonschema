package db

type FieldType struct {
	Name   string
	Format string
}

type Field struct {
	Name string
	Type *FieldType
}

type Table struct {
	Name        string
	Fields      []*Field
	PrimaryKeys []string
}
