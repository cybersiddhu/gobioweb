package gobioweb

import (
	dbi "database/sql"
	_ "github.com/mattn/go-sqlite3"
	"testing"
)

const tblStmt = `
CREATE TABLE users (
	 id INTEGER PRIMARY KEY autoincrement not null,
	 email TEXT UNIQUE,
	 password TEXT UNIQUE,
	 firstname TEXT,
	 lastname TEXT
)
`

func setUpTest() (dbh *dbi.DB, err error) {
	dbh, err = dbi.Open("sqlite3", ":memory:")
	if err != nil {
		return
	}
	if _, err = dbh.Exec(tblStmt); err != nil {
		return
	}
	return
}

func TestCreateUser(t *testing.T) {
	dbh, err := setUpTest()
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}

	u := &User{
		Email:     "tucker@tucker.com",
		FirstName: "tucker",
		LastName:  "jorn",
		Password:  "hashy",
	}
	if err := u.Create(dbh); err != nil {
		t.Errorf("error %s", err.Error())
	}
}

func TestValidateUser(t *testing.T) {
	dbh, err := setUpTest()
	if err != nil {
		t.Fatalf("error %s", err.Error())
	}

	u := &User{
		Email:     "tucker@tucker.com",
		FirstName: "tucker",
		LastName:  "jorn",
		Password:  "hashy",
	}

	if err := u.Create(dbh); err != nil {
		t.Errorf("error %s", err.Error())
	}

	if err := u.Validate(dbh); err != nil {
		t.Errorf("Error %s", err.Error())
	}
}
