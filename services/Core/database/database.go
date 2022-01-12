package database

import (
	"database/sql"
	"fmt"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jalexanderII/zero_fintech/services/Core/config"
	"github.com/jalexanderII/zero_fintech/services/Core/database/genDB"
	_ "github.com/lib/pq"
)

type CoreDB struct {
	Queries *genDB.Queries
	DB      *sql.DB
}

func NewCoreDB(db *sql.DB) *CoreDB {
	return &CoreDB{DB: db, Queries: genDB.New(db)}
}

func ConnectToDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("dbname=%s password=%s user=postgres sslmode=disable", config.GetEnv("COREDB_NAME"), config.GetEnv("COREDB_PASSWORD")))
	if err != nil {
		panic(err)
	}
	return db, err
}
