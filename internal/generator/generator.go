package generator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"text/template"

	log "github.com/sirupsen/logrus"
	"github.com/tgallant/db2jsonschema/internal/schema"
	"gopkg.in/yaml.v2"
)

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

const (
	defaultFormat     = "json"
	defaultSchemaType = "https://json-schema.org/draft/2020-12/schema"
	defaultIdTemplate = "{{ .Name }}.{{ .Format }}"
)

type Request struct {
	Tables     []*schema.Table
	Format     string
	Outdir     string
	SchemaType string
	IdTemplate string
}

func (r *Request) GetFormat() string {
	if len(r.Format) > 0 {
		return r.Format
	}
	return defaultFormat
}

func (r *Request) GetSchemaType() string {
	if len(r.SchemaType) > 0 {
		return r.SchemaType
	}
	return defaultSchemaType
}

func (r *Request) GetIdTemplate() string {
	if len(r.IdTemplate) > 0 {
		return r.IdTemplate
	}
	return defaultIdTemplate
}

type IdTemplateOptions struct {
	Name   string
	Format string
}

func (r *Request) FormatIdTemplate(name string) (string, error) {
	opts := &IdTemplateOptions{
		Name:   name,
		Format: r.GetFormat(),
	}
	idTemplate, err := template.New("idTemplate").Parse(r.GetIdTemplate())
	if err != nil {
		return "", err
	}
	var idValue bytes.Buffer
	err = idTemplate.Execute(&idValue, opts)
	if err != nil {
		return "", err
	}
	return idValue.String(), nil
}

func (r *Request) MakeDefinitionsDoc(tables []*schema.TableProperties) (*schema.DefinitionsDocument, error) {
	var definitions = make(map[string]map[string]*schema.JSONProperty)
	for _, t := range tables {
		props := schema.MakePropertiesMap(t.Properties)
		definitions[t.Name] = props
	}
	schemaId, err := r.FormatIdTemplate("definitions")
	if err != nil {
		return &schema.DefinitionsDocument{}, err
	}
	doc := &schema.DefinitionsDocument{
		Schema:      r.GetSchemaType(),
		Id:          schemaId,
		Title:       "Definitions",
		Definitions: definitions,
	}
	return doc, nil
}

func (r *Request) MakeSchema(tables []*schema.TableProperties) ([]*schema.JSONSchema, error) {
	var jsonSchemas []*schema.JSONSchema
	for _, t := range tables {
		properties := schema.MakePropertiesMap(t.Properties)
		schemaId, err := r.FormatIdTemplate(t.Name)
		if err != nil {
			return []*schema.JSONSchema{}, err
		}
		jsonSchema := &schema.JSONSchema{
			Schema:     r.GetSchemaType(),
			Id:         schemaId,
			Title:      t.Name,
			Type:       "object",
			Properties: properties,
		}
		jsonSchemas = append(jsonSchemas, jsonSchema)
	}
	return jsonSchemas, nil
}

func (r *Request) FormatSchema(schema interface{}) ([]byte, error) {
	format := r.GetFormat()
	switch format {
	case "json":
		return FormatJSON(schema)
	case "yaml":
		return FormatYAML(schema)
	default:
		return nil, fmt.Errorf("Unknown format: %s", format)
	}
}

func (r *Request) HandleStandardOutput(tables []*schema.TableProperties) error {
	doc, err := r.MakeDefinitionsDoc(tables)
	if err != nil {
		return err
	}
	res, err := r.FormatSchema(doc)
	if err != nil {
		return err
	}
	fmt.Println(string(res))
	return nil
}

func (r *Request) HandleDirectoryOutput(tables []*schema.TableProperties) error {
	err := os.MkdirAll(r.Outdir, os.ModePerm)
	if err != nil {
		return err
	}
	schema, err := r.MakeSchema(tables)
	if err != nil {
		return err
	}
	for _, s := range schema {
		res, err := r.FormatSchema(s)
		if err != nil {
			return err
		}
		filename := fmt.Sprintf("%s.%s", s.Title, r.GetFormat())
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
	var tables []*schema.TableProperties
	for _, table := range r.Tables {
		properties := schema.MakeTableProperties(table)
		tables = append(tables, properties)
	}
	if len(r.Outdir) > 0 {
		return r.HandleDirectoryOutput(tables)
	}
	return r.HandleStandardOutput(tables)
}
