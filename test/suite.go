package test

import (
	"context"
	"jackbot/db/migrate"
	"time"

	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"gorm.io/gorm"
)

const gormDbReadyMsg = "database system is ready to accept connections"

func StartTestDb(s *suite.Suite) testcontainers.Container {
	env := map[string]string{
		"POSTGRES_DB":       "test",
		"POSTGRES_USER":     "test",
		"POSTGRES_PASSWORD": "test",
	}
	port := "5432/tcp"
	ctx := context.Background()
	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "postgres:latest",
			ExposedPorts: []string{port},
			Env:          env,
			WaitingFor:   wait.ForLog(gormDbReadyMsg).WithPollInterval(100 * time.Millisecond).WithOccurrence(2),
		},
		Started: true,
	}
	testDbContainer, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		panic(err)
	}
	return testDbContainer
}

func RunMigrations(db *gorm.DB) error {
	err := migrate.Migrate(db)
	return err
}

func CleanTestDb(db *gorm.DB) error {
	res := db.Exec("DELETE FROM rows;DELETE FROM raffles;DELETE FROM users;DELETE FROM games;")
	return res.Error
}
