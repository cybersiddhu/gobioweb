package gobioweb

import "fmt"

var SessionName string = "user_session"

func (c *Controller) SaveUserSession(u *User) error {
	session, err := c.App.Session.Get(c.Request, SessionName)
	if err != nil {
		return err
	}

	fmt.Printf("got id to save %d\n",u.Id)
	session.Values["user_id"] = u.Id
	session.Save(c.Request, c.Response)
	return nil
}

func (c *Controller) DeleteUserSession() error {
	session, err := c.App.Session.Get(c.Request, SessionName)
	if err != nil {
		return err
	}
	if _, ok := session.Values["user_id"]; ok {
		delete(session.Values, "user_id")
		session.Save(c.Request, c.Response)
	}
	return nil
}

func (c *Controller) CurrentUser() (*User, error) {
	session, err := c.App.Session.Get(c.Request, SessionName)
	if err != nil {
		return nil, err
	}

	if id, ok := session.Values["user_id"]; ok {
		fmt.Printf("got id %d\n",id)
		u := &User{Id: id.(int64)}
		newuser,_ := u.FindById(c.App.Database)
		fmt.Printf("name:%s email:%s\n",newuser.FirstName,newuser.Email)
		return newuser, nil
	}
	return nil, nil
}

func (c *Controller) IsLoggedIn() bool {
	session, _ := c.App.Session.Get(c.Request, SessionName)
	if _, ok := session.Values["user_id"]; ok {
		return true
	}
	return false
}
