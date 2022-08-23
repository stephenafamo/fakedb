# FakeDB

This is copied from <http://golang.org/src/pkg/database/sql/fakedb_test.go>

It registers a fake database driver named `test`, just for testing.

It speaks a query language that's semantically similar to but
syntactically different and simpler than SQL.

## Query Methods

The query syntax is as follows

        WIPE
        CREATE|<tablename>|<col>=<type>,<col>=<type>,...
        INSERT|<tablename>|col=val,col2=val2,col3=?
        SELECT|<tablename>|projectcol1,projectcol2|filtercol=?,filtercol2=?
        SELECT|<tablename>|projectcol1,projectcol2|filtercol=?param1,filtercol2=?param2

* WIPE: Wipes all data in the database, including table information

        WIPE

* CREATE: Creates a new table

        CREATE|tableName|columnName1=columnType1,columnName2=columnType2

* DROP: Drops a table

        DROP|tableName

* INSERT: Inserts a row into a table

        INSERT|tableName|column1=?,column2=?

* SELECT: Queries the table

        SELECT|tableName|column1,column2|where=?,where2=?named

Any of these can be preceded by `PANIC|<method>|`, to cause the
named method on fakeStmt to panic.

Any of these can be proceeded by `WAIT|<duration>|`, to cause the
named method on fakeStmt to sleep for the specified duration.

Multiple of these can be combined when separated with a semicolon.

When opening a fakeDriver's database, it starts empty with no
tables. All tables and data are stored in memory only.

## Placeholders

In any query, `?` is used to denote a positional placeholder, and `?name` for a named placeholder

## Allowed types

* bool
* string
* byte
* int16
* int32
* in64
* float64
* datetime
* any

> **NOTE:** Every type can be nullable using nulltype. E.g nullstring for a nullable string

## Usage

As seen on [fakedb_test.go](fakedb_test.go)

```go
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
```
