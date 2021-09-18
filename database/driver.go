package database

import (
	"fmt"

	"github.com/tgallant/db2jsonschema/database/mysql"
	"github.com/tgallant/db2jsonschema/database/sqlite3"
	"github.com/tgallant/db2jsonschema/internal/schema"
)

type Driver interface {
	ReadTables() ([]*schema.Table, error)
}

type ConnectionInfo struct {
	Driver     string
	DataSource string
}

func NewConnection(i *ConnectionInfo) (Driver, error) {
	switch i.Driver {
	case "sqlite3":
		driver := &sqlite3.Driver{
			DataSource: i.DataSource,
		}
		return driver, nil
	case "mysql":
		driver := &mysql.Driver{
			DataSource: i.DataSource,
		}
		return driver, nil
	default:
		return nil, fmt.Errorf("Unknown driver: %s", i.Driver)
	}
}
