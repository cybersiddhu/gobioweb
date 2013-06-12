package gobioweb

import (
	"code.google.com/p/go.crypto/bcrypt"
	dbi "database/sql"
)

type User struct {
	Id        int64
	Email     string
	FirstName string
	LastName  string
	Password  string
}

const userCheckStmt = `
	SELECT users.password FROM users WHERE users.email = (:zaz)
`

const userInsertStmt = `
INSERT into users(email,firstname,lastname,password) 
values(:foo,:bar,:zas,:baz)
`

func (u *User) Create(dbh *dbi.DB) error {
	pb, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	rs, err := dbh.Exec(userInsertStmt, u.Email, u.FirstName, u.LastName, string(pb))
	if err != nil {
		return err
	}

	id, err := rs.LastInsertId()
	if err != nil {
		return err
	}
	u.Id = id
	return nil
}

func (u *User) Validate(dbh *dbi.DB) error {
	var hashedPass string
	err := dbh.QueryRow(userCheckStmt, u.Email).Scan(&hashedPass)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(u.Password))
	if err != nil {
		return err
	}
	return nil
}
