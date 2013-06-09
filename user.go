package gobioweb

import (
	"code.google.com/p/go.crypto/bcrypt"
	dbi "database/sql"
)

type User struct {
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

func CreateUser(dbh *dbi.DB, user *User) error {
	pb, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		 return err
	}

	_, err = dbh.Exec(userInsertStmt,user.Email, user.FirstName, user.LastName, string(pb))
	if err != nil {
		 return err
	}
	return nil
}

func ValidateUser(dbh *dbi.DB, user *User) error {
	var hashedPass string
	err := dbh.QueryRow(userCheckStmt, user.Email).Scan(&hashedPass)
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword([]byte(hashedPass), []byte(user.Password))
	if err != nil {
		return err
	}
	return nil
}


