package querier

import (
	"testing"
)

func TestQuerier(t *testing.T) {
	q := NewQuerier()

	query := q.Select("column1", "column2").From("table1").Build()

	expectedQuery := "SELECT column1, column2 FROM table1"

	if query != expectedQuery {
		t.Errorf("Expected '%s', but got '%s'", expectedQuery, query)
	}
}
