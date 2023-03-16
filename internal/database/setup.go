package database

import (
	_ "embed"
	"errors"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

//go:embed settings.yaml
var settingsYaml []byte

//go:embed ddl.sql
var dbSchemaSetup string

type dbParams struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DBName   string `yaml:"db_name"`
	Params   string `yaml:"params"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
}

func SetUpDB(env string) (*sqlx.DB, error) {
	if env == "" {
		return nil, errors.New("empty environment passed to SetUpDB")
	}

	var envsSettings map[string]dbParams
	err := yaml.Unmarshal(settingsYaml, &envsSettings)
	if err != nil {
		return nil, fmt.Errorf("can't read db settings: %v", err)
	}

	settings, ok := envsSettings[env]
	if !ok {
		return nil, fmt.Errorf("no db settings is found for environment '%s'", env)
	}

	dbAddress := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?%s",
		settings.Username, settings.Password, settings.Host, settings.Port, settings.DBName, settings.Params,
	)
	log.Info().Msgf("Trying to connect mysql address: %s", dbAddress)
	db, err := sqlx.Connect("mysql", dbAddress)
	if err != nil {
		return nil, err
	}
	db.SetConnMaxLifetime(3 * time.Second)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	if env == "test" {
		_, err = db.Exec(fmt.Sprintf("drop database %s", settings.DBName))
		_, err = db.Exec(fmt.Sprintf("create database %s", settings.DBName))
		_, err = db.Exec(fmt.Sprintf("use %s", settings.DBName))
		_, err = db.Exec(dbSchemaSetup)
		if err != nil {
			return nil, fmt.Errorf("failed to drop test db: %w", err)
		}
	}
	return db, nil
}
