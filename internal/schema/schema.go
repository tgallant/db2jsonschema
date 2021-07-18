package schema

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

func MakeTableProperties(t *Table) *TableProperties {
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
