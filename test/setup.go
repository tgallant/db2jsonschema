package test

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Genre struct {
	gorm.Model
	Name        string
	Description string
}

type Artist struct {
	gorm.Model
	Name        string
	Description string
	Tracks      []Track `gorm:"many2many:artist_tracks;"`
}

type Album struct {
	gorm.Model
	Title string
}

type Track struct {
	gorm.Model
	Title    string
	Duration uint
	AlbumID  uint
	GenreID  uint
	Album    Album
	Genre    Genre
}

func MigrateTables(db *gorm.DB) {
	db.AutoMigrate(
		&Genre{},
		&Artist{},
		&Album{},
		&Track{},
	)
}

func SetupSQLite(datasource string) (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open(datasource), &gorm.Config{})
	if err != nil {
		return &gorm.DB{}, err
	}
	return db, nil
}

type TestDB struct {
	Driver     string
	DataSource string
}

func (t *TestDB) NewConnection() (*gorm.DB, error) {
	switch t.Driver {
	case "sqlite3":
		return SetupSQLite(t.DataSource)
	default:
		return nil, fmt.Errorf("Unknown driver: %s", t.Driver)
	}
}

func (t *TestDB) Setup() error {
	db, err := t.NewConnection()
	if err != nil {
		return err
	}
	MigrateTables(db)
	return nil
}