package db2jsonschema

import (
	log "github.com/sirupsen/logrus"
	"github.com/tgallant/db2jsonschema/database"
	"github.com/tgallant/db2jsonschema/internal/generator"
	"github.com/tgallant/db2jsonschema/internal/schema"
)

func MakeLookupMap(items []string) map[string]bool {
	var lookupMap = make(map[string]bool)
	for _, item := range items {
		lookupMap[item] = true
	}
	return lookupMap
}

type Request struct {
	Driver     string
	DataSource string
	Format     string
	Outdir     string
	SchemaType string
	IdTemplate string
	Includes   []string
	Excludes   []string
}

func (r *Request) FilterTables(tables []*schema.Table) []*schema.Table {
	hasIncludes := len(r.Includes) > 0
	hasExcludes := len(r.Excludes) > 0
	if !hasIncludes && !hasExcludes {
		return tables
	}
	var filteredTables []*schema.Table
	includesMap := MakeLookupMap(r.Includes)
	excludesMap := MakeLookupMap(r.Excludes)
	for _, t := range tables {
		_, matchesInclude := includesMap[t.Name]
		_, matchesExclude := excludesMap[t.Name]
		if hasIncludes && !matchesInclude {
			continue
		}
		if hasExcludes && matchesExclude {
			continue
		}
		filteredTables = append(filteredTables, t)
	}
	return filteredTables
}

func (r *Request) Perform() error {
	info := &database.ConnectionInfo{
		Driver:     r.Driver,
		DataSource: r.DataSource,
	}
	log.WithFields(log.Fields{
		"connectionInfo": info,
	}).Debug("Connecting to database")
	driver, err := database.NewConnection(info)
	if err != nil {
		return err
	}
	tables, err := driver.ReadTables()
	if err != nil {
		return err
	}
	filteredTables := r.FilterTables(tables)
	request := generator.Request{
		Tables:     filteredTables,
		Format:     r.Format,
		Outdir:     r.Outdir,
		SchemaType: r.SchemaType,
		IdTemplate: r.IdTemplate,
	}
	log.WithFields(log.Fields{
		"generatorRequest": request,
	}).Debug("Generating Schemas")
	err = request.Perform()
	if err != nil {
		return err
	}
	return nil
}
