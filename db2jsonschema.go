package db2jsonschema

import (
	log "github.com/sirupsen/logrus"
	"github.com/tgallant/db2jsonschema/database"
	"github.com/tgallant/db2jsonschema/internal/generator"
)

type Request struct {
	Driver     string
	DataSource string
	Format     string
	Outdir     string
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
	request := generator.Request{
		Tables: tables,
		Format: r.Format,
		Outdir: r.Outdir,
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
