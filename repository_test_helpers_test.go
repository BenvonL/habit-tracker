package main

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// newMockRepository returns a MySQLRepository backed by sqlmock. Queries are
// matched exactly (after whitespace normalisation), so a test fails when the
// SQL a repository method sends changes.
func newMockRepository(t *testing.T) (MySQLRepository, sqlmock.Sqlmock) {
	t.Helper()
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err != nil {
		t.Fatalf("creating sqlmock: %v", err)
	}
	t.Cleanup(func() {
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("unmet sql expectations: %v", err)
		}
		db.Close()
	})
	return MySQLRepository{DB: db}, mock
}
