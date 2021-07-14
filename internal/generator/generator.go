package generator

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	log "github.com/sirupsen/logrus"
	"github.com/tgallant/db2jsonschema/internal/db"
	"gopkg.in/yaml.v2"
)

type JSONProperty struct {
	Name   string `json:"name" yaml:"name"`
	Type   string `json:"type" yaml:"type"`
	Format string `json:"format,omitempty" yaml:"format,omitempty"`
}

type JSONSchema struct {
	Schema     string                   `json:"$schema" yaml:"$schema"`
	Id         string                   `json:"$id" yaml:"$id"`
	Title      string                   `json:"title" yaml:"title"`
	Type       string                   `json:"type" yaml:"type"`
	Properties map[string]*JSONProperty `json:"properties" yaml:"properties"`
}

type DefinitionsDocument struct {
	Schema      string                              `json:"$schema" yaml:"$schema"`
	Id          string                              `json:"$id" yaml:"$id"`
	Title       string                              `json:"title" yaml:"title"`
	Definitions map[string]map[string]*JSONProperty `json:"definitions" yaml:"definitions"`
}

type TableProperties struct {
	Name       string
	Properties []*JSONProperty
}

func MakeTableProperties(t *db.Table) *TableProperties {
	var properties []*JSONProperty
	for _, field := range t.Fields {
		prop := &JSONProperty{
			Name:   field.Name,
			Type:   field.Type.Name,
			Format: field.Type.Format,
		}
		properties = append(properties, prop)
	}
	tableProperties := &TableProperties{
		Name:       t.Name,
		Properties: properties,
	}
	return tableProperties
}

func MakePropertiesMap(props []*JSONProperty) map[string]*JSONProperty {
	var propsMap = make(map[string]*JSONProperty)
	for _, p := range props {
		propsMap[p.Name] = p
	}
	return propsMap
}

func MakeDefinitionsDoc(tables []*TableProperties) *DefinitionsDocument {
	var definitions = make(map[string]map[string]*JSONProperty)
	for _, t := range tables {
		props := MakePropertiesMap(t.Properties)
		definitions[t.Name] = props
	}
	doc := &DefinitionsDocument{
		Schema:      "https://json-schema.org/draft/2020-12/schema",
		Id:          "https://example.com/product.schema.json",
		Title:       "Example",
		Definitions: definitions,
	}
	return doc
}

func MakeSchemas(tables []*TableProperties) []*JSONSchema {
	var schemas []*JSONSchema
	for _, t := range tables {
		properties := MakePropertiesMap(t.Properties)
		schema := &JSONSchema{
			Schema:     "https://json-schema.org/draft/2020-12/schema",
			Id:         "https://example.com/product.schema.json",
			Title:      t.Name,
			Type:       "object",
			Properties: properties,
		}
		schemas = append(schemas, schema)
	}
	return schemas
}

func FormatJSON(schema interface{}) ([]byte, error) {
	res, err := json.MarshalIndent(schema, "", "  ")
	if err != nil {
		return nil, err
	}
	return res, nil
}

func FormatYAML(schema interface{}) ([]byte, error) {
	res, err := yaml.Marshal(schema)
	if err != nil {
		return nil, err
	}
	return res, nil
}

type Request struct {
	Tables []*db.Table
	Format string
	Outdir string
}

func (r *Request) FormatSchema(schema interface{}) ([]byte, error) {
	switch r.Format {
	case "json":
		return FormatJSON(schema)
	case "yaml":
		return FormatYAML(schema)
	default:
		return nil, fmt.Errorf("Unknown format: %s", r.Format)
	}
}

func (r *Request) HandleStandardOutput(tables []*TableProperties) error {
	doc := MakeDefinitionsDoc(tables)
	res, err := r.FormatSchema(doc)
	if err != nil {
		return err
	}
	fmt.Println(string(res))
	return nil
}

func (r *Request) HandleDirectoryOutput(tables []*TableProperties) error {
	err := os.MkdirAll(r.Outdir, os.ModePerm)
	if err != nil {
		return err
	}
	schemas := MakeSchemas(tables)
	for _, s := range schemas {
		res, err := r.FormatSchema(s)
		if err != nil {
			return err
		}
		filename := fmt.Sprintf("%s.%s", s.Title, r.Format)
		outputPath := filepath.Join(r.Outdir, filename)
		log.Infof("Writing to %s", outputPath)
		err = os.WriteFile(outputPath, res, 0666)
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *Request) Perform() error {
	var tables []*TableProperties
	for _, table := range r.Tables {
		properties := MakeTableProperties(table)
		tables = append(tables, properties)
	}
	if len(r.Outdir) > 0 {
		return r.HandleDirectoryOutput(tables)
	}
	return r.HandleStandardOutput(tables)
}
