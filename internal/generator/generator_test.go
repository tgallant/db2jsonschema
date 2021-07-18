package generator

import (
	"github.com/stretchr/testify/assert"
	"github.com/tgallant/db2jsonschema/internal/schema"
	"testing"
)

func makeDbTable() *schema.Table {
	var fields []*schema.Field
	testField := schema.Field{
		Name: "exampleField",
		Type: &schema.FieldType{
			Name:   "string",
			Format: "",
		},
	}
	testField2 := schema.Field{
		Name: "UserId",
		Type: &schema.FieldType{
			Name:   "number",
			Format: "",
		},
	}
	fields = append(fields, &testField, &testField2)
	table := &schema.Table{
		Name:   "Testing",
		Fields: fields,
	}
	return table
}

func TestMakeDefinitionsDoc(t *testing.T) {
	var props []*schema.TableProperties
	table := makeDbTable()
	p := schema.MakeTableProperties(table)
	props = append(props, p)
	r := &Request{}
	doc, err := r.MakeDefinitionsDoc(props)
	assert.Nilf(t, err, "creating the definitions doc for %s should succeed", table.Name)
	assert.Equal(t, 1, len(doc.Definitions), "should have 1 definition")
	def := doc.Definitions["Testing"]
	assert.Equal(t, 2, len(def), "should have 2 properties")
	property := def["UserId"]
	assert.Equal(t, "number", property.Type, "type should be `number`")
}

func TestMakeDefinitionsDocWithIdTemplate(t *testing.T) {
	var props []*schema.TableProperties
	table := makeDbTable()
	p := schema.MakeTableProperties(table)
	props = append(props, p)
	idValue := "https://example.com/schemas/test.json"
	r := &Request{
		IdTemplate: idValue,
	}
	doc, err := r.MakeDefinitionsDoc(props)
	assert.Nilf(t, err, "creating the definitions doc for %s should succeed", table.Name)
	assert.Equalf(t, idValue, doc.Id, "the $id value should be %s", idValue)
	assert.Equal(t, 1, len(doc.Definitions), "should have 1 definition")
	def := doc.Definitions["Testing"]
	assert.Equal(t, 2, len(def), "should have 2 properties")
	property := def["UserId"]
	assert.Equal(t, "number", property.Type, "type should be `number`")
}

func TestMakeDefinitionsDocWithSchemaType(t *testing.T) {
	var props []*schema.TableProperties
	table := makeDbTable()
	p := schema.MakeTableProperties(table)
	props = append(props, p)
	schemaValue := "https://example.com/schema"
	r := &Request{
		SchemaType: schemaValue,
	}
	doc, err := r.MakeDefinitionsDoc(props)
	assert.Nilf(t, err, "creating the definitions doc for %s should succeed", table.Name)
	assert.Equalf(t, schemaValue, doc.Schema, "the $schema value should be %s", schemaValue)
	assert.Equal(t, 1, len(doc.Definitions), "should have 1 definition")
	def := doc.Definitions["Testing"]
	assert.Equal(t, 2, len(def), "should have 2 properties")
	property := def["UserId"]
	assert.Equal(t, "number", property.Type, "type should be `number`")
}

func TestPerform(t *testing.T) {
	var tables []*schema.Table
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
