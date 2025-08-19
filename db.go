package exp

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"
)

// DB stores the connection information of a database.
type DB struct {
	Host    string
	Port    int
	Name    string
	User    string
	Pass    string
	SSLMode string

	Schema string

	*sql.DB
}

// LoadDB loads database information from the environment.
func LoadDB() (*DB, error) {
	db := &DB{
		Host:    "127.0.0.1",
		Port:    5432,
		Name:    os.Getenv("DB_NAME"),
		User:    os.Getenv("DB_USER"),
		Pass:    os.Getenv("DB_PASS"),
		SSLMode: os.Getenv("DB_SSLMODE"),
		Schema:  "schema.sql",
	}
	if tmp := os.Getenv("DB_HOST"); tmp != "" {
		db.Host = tmp
	}
	if tmp := os.Getenv("DB_PORT"); tmp != "" {
		var err error
		db.Port, err = strconv.Atoi(tmp)
		if err != nil {
			return nil, fmt.Errorf("fail to parse DB_PORT: %v", err)
		}
	}
	if db.User == "" {
		return nil, fmt.Errorf("DB_USER is required")
	}
	if db.Name == "" {
		db.Name = db.User
	}
	if tmp := os.Getenv("DB_SCHEMA"); tmp != "" {
		db.Schema = tmp
	}
	if db.SSLMode == "" {
		db.SSLMode = "required"
	}

	var err error
	db.DB, err = sql.Open("postgres", db.connString())
	if err != nil {
		return nil, fmt.Errorf("fail to open connection: %v", err)
	}

	return db, nil
}

func (db *DB) Up() error {
	data, err := os.ReadFile(db.Schema)
	if err != nil {
		return fmt.Errorf("read %v: %v", db.Schema, err)
	}

	_, err = db.Exec(string(data))
	if err != nil {
		return fmt.Errorf("exec %v: %v", db.Schema, err)
	}

	return nil
}

func (db *DB) Drop() error {
	_, err := db.Exec(fmt.Sprintf("DROP OWNED BY %q", db.User))
	if err != nil {
		return fmt.Errorf("exec: %v", err)
	}
	return nil
}

func (db *DB) connString() string {
	return fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s sslmode=%s",
		db.Host, db.Port, db.Name, db.User, db.Pass, db.SSLMode)
}
