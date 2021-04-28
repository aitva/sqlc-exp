package main

import (
	"fmt"
	"os"
	"strconv"
)

type dbInfo struct {
	Host string
	Port int
	Name string
	User string
	Pass string
}

func (info *dbInfo) String() string {
	return fmt.Sprintf("host=%s port=%d dbname=%s user=%s password=%s",
		info.Host, info.Port, info.Name, info.User, info.Pass)
}

func loadDBInfo() (*dbInfo, error) {
	info := &dbInfo{
		Host: os.Getenv("DB_HOST"),
		Port: 5432,
		Name: os.Getenv("DB_NAME"),
		User: os.Getenv("DB_USER"),
		Pass: os.Getenv("DB_PASS"),
	}
	if info.Host == "" {
		return nil, fmt.Errorf("DB_HOST is required")
	}
	if info.User == "" {
		return nil, fmt.Errorf("DB_USER is required")
	}
	if info.Name == "" {
		info.Name = info.User
	}
	if tmp := os.Getenv("DB_PORT"); tmp != "" {
		var err error
		info.Port, err = strconv.Atoi(tmp)
		if err != nil {
			return nil, fmt.Errorf("fail to parse DB_PORT: %v", err)
		}
	}
	return info, nil
}
