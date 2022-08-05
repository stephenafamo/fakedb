package fakedb_test

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/stephenafamo/fakedb"
)

func TestQuery(t *testing.T) {
	ctx := context.Background()

	db, err := sql.Open("test", "identifier")
	if err != nil {
		t.Fatalf("Error opening testdb %v", err)
	}

	exec(t, db, "CREATE|users|id=int64,name=string")
	exec(t, db, "INSERT|users|id=?,name=?", 1, "foo")
	exec(t, db, "INSERT|users|id=?,name=?", 2, "bar")

	rows, err := db.QueryContext(ctx, "SELECT|users|id,name|")
	if err != nil {
		t.Fatal(err)
	}

	users := []string{"foo", "bar"}
	for rows.Next() {
		var id int
		var name string
		rows.Scan(&id, &name)

		expectedName := users[id-1]
		if name != users[id-1] {
			t.Fatalf("User %d should have name %q but had %q", id, expectedName, name)
		}
	}
}

func exec(tb testing.TB, exec *sql.DB, query string, args ...interface{}) sql.Result {
	tb.Helper()
	result, err := exec.ExecContext(context.Background(), query, args...)
	if err != nil {
		tb.Fatalf("Exec of %q: %v", query, err)
	}

	return result
}
