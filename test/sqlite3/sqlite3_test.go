package sqlite3_test

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tgallant/db2jsonschema"
	"github.com/tgallant/db2jsonschema/internal/generator"
	"github.com/tgallant/db2jsonschema/test"
	"gopkg.in/yaml.v2"
)

var tempDir string

var testDB = &test.TestDB{
	Driver:     "sqlite3",
	DataSource: "",
}

var expectedTables = map[string]bool{
	"genres":        true,
	"artists":       true,
	"albums":        true,
	"tracks":        true,
	"artist_tracks": true,
}

func TestMain(m *testing.M) {
	tempDir = os.TempDir()
	tempFile, err := os.CreateTemp(tempDir, "test.*.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	testDB.DataSource = tempFile.Name()
	err = testDB.Setup()
	if err != nil {
		fmt.Println(err)
		return
	}
	m.Run()
}

func TestYAMLOutput(t *testing.T) {
	schemaPath := filepath.Join(tempDir, "schemas_yaml")
	err := os.MkdirAll(schemaPath, os.ModePerm)
	assert.Nilf(t, err, "creating the dir %s should succeed", schemaPath)
	req := &db2jsonschema.Request{
		Driver:     testDB.Driver,
		DataSource: testDB.DataSource,
		Format:     "yaml",
		Outdir:     schemaPath,
	}
	err = req.Perform()
	assert.Nil(t, err, "performing the request should succeed")
	dir, err := os.ReadDir(schemaPath)
	assert.Nilf(t, err, "reading dir %s should succeed", schemaPath)
	assert.Equal(t, 5, len(dir), "there should be 5 schemas")
	for _, file := range dir {
		filename := file.Name()
		schemaName := strings.TrimSuffix(filename, filepath.Ext(filename))
		_, exists := expectedTables[schemaName]
		msg := fmt.Sprintf("The %s schema should exist", schemaName)
		assert.True(t, exists, msg)
		schema := &generator.JSONSchema{}
		fullPath := filepath.Join(schemaPath, filename)
		contents, err := os.ReadFile(fullPath)
		assert.Nilf(t, err, "reading file %s should succeed", filename)
		err = yaml.Unmarshal(contents, schema)
		assert.Nilf(t, err, "unmarshalling file %s should succeed", filename)
		assert.NotEmptyf(t, schema.Properties, "the schema %s should have properties")
	}
}

func TestJSONOutput(t *testing.T) {
	schemaPath := filepath.Join(tempDir, "schemas_json")
	err := os.MkdirAll(schemaPath, os.ModePerm)
	assert.Nilf(t, err, "creating the dir %s should succeed", schemaPath)
	req := &db2jsonschema.Request{
		Driver:     testDB.Driver,
		DataSource: testDB.DataSource,
		Format:     "json",
		Outdir:     schemaPath,
	}
	err = req.Perform()
	assert.Nil(t, err, "performing the request should succeed")
	dir, err := os.ReadDir(schemaPath)
	assert.Nilf(t, err, "reading dir %s should succeed", schemaPath)
	assert.Equal(t, 5, len(dir), "there should be 5 schemas")
	for _, file := range dir {
		filename := file.Name()
		schemaName := strings.TrimSuffix(filename, filepath.Ext(filename))
		_, exists := expectedTables[schemaName]
		msg := fmt.Sprintf("The %s schema should exist", schemaName)
		assert.True(t, exists, msg)
		schema := &generator.JSONSchema{}
		fullPath := filepath.Join(schemaPath, filename)
		contents, err := os.ReadFile(fullPath)
		assert.Nilf(t, err, "reading file %s should succeed", filename)
		err = json.Unmarshal(contents, schema)
		assert.Nilf(t, err, "unmarshalling file %s should succeed", filename)
		assert.NotEmptyf(t, schema.Properties, "the schema %s should have properties")
	}
}
