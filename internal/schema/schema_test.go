package schema

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func makeDbTable() *Table {
	var fields []*Field
	testField := Field{
		Name: "exampleField",
		Type: &FieldType{
			Name:   "string",
			Format: "",
		},
	}
	testField2 := Field{
		Name: "UserId",
		Type: &FieldType{
			Name:   "number",
			Format: "",
		},
	}
	fields = append(fields, &testField, &testField2)
	table := &Table{
		Name:   "Testing",
		Fields: fields,
	}
	return table
}

func TestMakeTableProperties(t *testing.T) {
	table := makeDbTable()
	p := MakeTableProperties(table)
	assert.Equal(t, "Testing", p.Name, "name should be `Testing`")
	assert.Equal(t, 2, len(p.Properties), "should have 2 properties")
}

func TestMakePropertiesMap(t *testing.T) {
	var props []*JSONProperty
	prop := &JSONProperty{
		Name: "email",
		Type: "string",
	}
	props = append(props, prop)
	m := MakePropertiesMap(props)
	assert.Equal(t, 1, len(m), "should have 1 property")
	p := m["email"]
	assert.Equal(t, "string", p.Type, "type should be `string`")
}
