package generator

import (
	"github.com/stretchr/testify/assert"
	"github.com/tgallant/db2jsonschema/internal/db"
	"testing"
)

func makeDbTable() *db.Table {
	var fields []*db.Field
	testField := db.Field{
		Name: "exampleField",
		Type: &db.FieldType{
			Name:   "string",
			Format: "",
		},
	}
	testField2 := db.Field{
		Name: "UserId",
		Type: &db.FieldType{
			Name:   "number",
			Format: "",
		},
	}
	fields = append(fields, &testField, &testField2)
	table := &db.Table{
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

func TestMakeDefinitionsDoc(t *testing.T) {
	var props []*TableProperties
	table := makeDbTable()
	p := MakeTableProperties(table)
	props = append(props, p)
	doc := MakeDefinitionsDoc(props)
	assert.Equal(t, 1, len(doc.Definitions), "should have 1 definition")
	def := doc.Definitions["Testing"]
	assert.Equal(t, 2, len(def), "should have 2 properties")
	property := def["UserId"]
	assert.Equal(t, "number", property.Type, "type should be `number`")
}

func TestPerform(t *testing.T) {
	var tables []*db.Table
	table := makeDbTable()
	tables = append(tables, table)
	request := Request{
		Tables: tables,
		Format: "json",
		Outdir: "",
	}
	err := request.Perform()
	assert.Nil(t, err, "performing the request should succeed")
}