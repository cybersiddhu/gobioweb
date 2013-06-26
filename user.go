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

const userFindStmt = `
SELECT users.id,users.firstname,users.lastname FROM users
WHERE users.email = ?
`

const userFindStmt2 = `
SELECT user.email,users.firstname,users.lastname FROM users
WHERE users.id = ?
`

//Functions
func FindUserByEmail(dbh *dbi.DB, email string) (*User, error) {
	var first, last string
	var id int64
	err := dbh.QueryRow(userFindStmt, email).Scan(&id, &first, &last)
	if err != nil {
		return nil, err
	}
	u := &User{
		Id:        id,
		FirstName: first,
		LastName:  last,
		Email:     email,
	}
	return u, nil
}

func FindUserById(dbh *dbi.DB, id int64) (*User, error) {
	var first, last, email string
	err := dbh.QueryRow(userFindStmt2, id).Scan(&email, &first, &last)
	if err != nil {
		return nil, err
	}

	u := &User{
		Id:        id,
		FirstName: first,
		LastName:  last,
		Email:     email,
	}
	return u, nil
}

func (u *User) FindById(dbh *dbi.DB) (*User, error) {
	var first, last, email string
	err := dbh.QueryRow(userFindStmt2, u.Id).Scan(&email, &first, &last)
	if err != nil {
		return nil, err
	}
	u.FirstName = first
	u.LastName = last
	u.Email = email
	return u, nil
}

func CreateUser(dbh *dbi.DB, u *User) (*User, error) {
	pb, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	rs, err := dbh.Exec(userInsertStmt, u.Email, u.FirstName, u.LastName, string(pb))
	if err != nil {
		return nil, err
	}

	id, err := rs.LastInsertId()
	if err != nil {
		return nil, err
	}
	newuser := &User{
		Id:        id,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}
	return newuser, nil
}

func ValidateUser(dbh *dbi.DB, u *User) error {
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

//Methods
func (u *User) Find(dbh *dbi.DB) (*User, error) {
	var first, last string
	var id int64
	err := dbh.QueryRow(userFindStmt, u.Email).Scan(&id, &first, &last)
	if err != nil {
		return nil, err
	}
	u.Id = id
	u.FirstName = first
	u.LastName = last
	return u, nil
}

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
