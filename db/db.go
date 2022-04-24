package db

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type MissingEnvVarsError struct {
	Vars []string
}

func (e *MissingEnvVarsError) Error() string {
	errString := ""
	for i, v := range e.Vars {
		if i != len(e.Vars)-1 {
			errString += fmt.Sprintf("%s, ", v)
		} else {
			errString += v
		}
	}
	return errString
}

func NewConn() (*gorm.DB, error) {
	dsn, err := getDSN()
	if err != nil {
		return nil, err
	}
	conn, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func getDSN() (string, error) {
	var missingEnvVars []string
	host, ok := os.LookupEnv("POSTGRES_HOST")
	if !ok {
		missingEnvVars = append(missingEnvVars, "POSTGRES_HOST")
	}
	port, ok := os.LookupEnv("POSTGRES_PORT")
	if !ok {
		missingEnvVars = append(missingEnvVars, "POSTGRES_PORT")
	}
	user, ok := os.LookupEnv("POSTGRES_USER")
	if !ok {
		missingEnvVars = append(missingEnvVars, "POSTGRES_USER")
	}
	database, ok := os.LookupEnv("POSTGRES_DB")
	if !ok {
		missingEnvVars = append(missingEnvVars, "POSTGRES_DB")
	}
	if len(missingEnvVars) != 0 {
		return "", &MissingEnvVarsError{Vars: missingEnvVars}
	}
	url := fmt.Sprintf(
		"host=%s user=%s dbname=%s port=%s sslmode=disable",
		host, user, database, port,
	)
	password, passwordOk := os.LookupEnv("POSTGRES_PASSWORD")
	if passwordOk {
		url += fmt.Sprintf(" password=%s", password)
	}
	return url, nil
}
