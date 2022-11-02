package db_creater //nolint:stylecheck

import (
	"database/sql"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5"
)

type DBCreator struct {
	user     string
	password string
	address  string
	port     int
	db       string
}

func New(user, password, address, db string, port int) *DBCreator {
	return &DBCreator{
		user:     user,
		password: password,
		address:  address,
		port:     port,
		db:       db,
	}
}

func (c *DBCreator) CreateWithMigration(migrationsDir string) (string, error) {
	dbName, err := c.generateRandomDBName()
	if err != nil {
		return "", err
	}

	err = c.create(dbName)
	if err != nil {
		return "", err
	}

	db, err := sql.Open("postgres", c.connectionString(dbName))
	if err != nil {
		return "", fmt.Errorf("failed to open connection %w", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return "", fmt.Errorf("failed to create driver for migrate %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		migrationsDir,
		dbName, driver)
	if err != nil {
		return "", fmt.Errorf("failed to create migrate db instance %w", err)
	}

	if err = m.Up(); err != nil {
		return "", fmt.Errorf("failed to migrate db %w", err)
	}

	sourceErr, databaseErr := m.Close()
	if sourceErr != nil {
		return "", fmt.Errorf("failed to close source migrator %w", sourceErr)
	}

	if databaseErr != nil {
		return "", fmt.Errorf("failed to close database migrator %w", databaseErr)
	}

	err = db.Close()
	if err != nil {
		return "", fmt.Errorf("failed to close db: %w", err)
	}

	return dbName, nil
}

func (c *DBCreator) Drop(dbName string) error {
	db, err := sql.Open("postgres", c.connectionString(c.db))
	if err != nil {
		return fmt.Errorf("failed to open connection %w", err)
	}
	defer db.Close()

	if _, err := db.Exec(fmt.Sprintf("DROP DATABASE %q", dbName)); err != nil {
		return fmt.Errorf("failed to drop database with name %s %w", dbName, err)
	}
	return nil
}

func (c *DBCreator) create(dbName string) error {
	db, err := sql.Open("postgres", c.connectionString(c.db))
	if err != nil {
		return fmt.Errorf("failed to open connection %w", err)
	}

	if _, err = db.Exec(fmt.Sprintf("CREATE DATABASE %q", dbName)); err != nil {
		return fmt.Errorf("failed to create database with name %s %w", dbName, err)
	}
	if err = db.Close(); err != nil {
		return fmt.Errorf("failed to close connection %w", err)
	}

	return nil
}

func (c *DBCreator) generateRandomDBName() (string, error) {
	dbNamePrefix, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("failed to create random uuid %w", err)
	}
	return "test_db_" + dbNamePrefix.String(), nil
}

func (c *DBCreator) connectionString(dbName string) string {
	return fmt.Sprintf(
		"host=%s port=%d dbname=%s user=%s password=%s sslmode=disable",
		c.address,
		c.port,
		dbName,
		c.user,
		c.password)
}
