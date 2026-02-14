package db

import (
	"context"
	"fmt"

	"dvarapala/ent"
	"dvarapala/ent/migrate"

	_ "github.com/mattn/go-sqlite3"
)

// NewSQLiteClient creates a new ent.Client for SQLite.
func NewSQLiteClient(path string) (*ent.Client, error) {
	client, err := ent.Open("sqlite3", fmt.Sprintf("file:%s?cache=shared&_fk=1", path))
	if err != nil {
		return nil, fmt.Errorf("failed opening connection to sqlite: %v", err)
	}
	
	// Run the auto migration tool if you want to keep it simple, 
	// OR use the versioned migrations.
	// For versioned migrations, we typically use the migrate package.
	if err := client.Schema.Create(context.Background(), migrate.WithGlobalUniqueID(true)); err != nil {
		client.Close()
		return nil, fmt.Errorf("failed creating schema resources: %v", err)
	}
	
	return client, nil
}
