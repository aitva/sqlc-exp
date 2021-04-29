package exp

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

// DB stores the connection information of a database.
type DB struct {
	Host string
	Port int
	Name string
	User string
	Pass string

	Schema string

	*sql.DB
}

// LoadDB loads database information from the environment.
func LoadDB() (*DB, error) {
	db := &DB{
		Host:   os.Getenv("DB_HOST"),
		Port:   5432,
		Name:   os.Getenv("DB_NAME"),
		User:   os.Getenv("DB_USER"),
		Pass:   os.Getenv("DB_PASS"),
		Schema: os.Getenv("DB_SCHEMA"),
	}
	if db.Host == "" {
		return nil, fmt.Errorf("DB_HOST is required")
	}
	if db.User == "" {
		return nil, fmt.Errorf("DB_USER is required")
	}
	if db.Name == "" {
		db.Name = db.User
	}
	if tmp := os.Getenv("DB_PORT"); tmp != "" {
		var err error
		db.Port, err = strconv.Atoi(tmp)
		if err != nil {
			return nil, fmt.Errorf("fail to parse DB_PORT: %v", err)
		}
	}
	if db.Schema == "" {
		return nil, fmt.Errorf("DB_SCHEMA is required")
	}

	var err error
	db.DB, err = sql.Open("postgres", db.connString())
	if err != nil {
		return nil, fmt.Errorf("fail to open connection: %v", err)
	}

	return db, nil
}

func (db *DB) Up() error {
	data, err := ioutil.ReadFile(db.Schema)
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
	_, err := db.Exec("DROP OWNED BY " + db.User)
	if err != nil {
		return fmt.Errorf("exec: %v", err)
	}
	return nil
}

func (db *DB) connString() string {
	return fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s",
		db.Host, db.Port, db.Name, db.User, db.Pass)
}
