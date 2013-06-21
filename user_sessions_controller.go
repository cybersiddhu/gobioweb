package gobioweb

var SessionName string = "user_session"

func (c *Controller) SaveUserSession(u *User) error {
	session, err := c.App.Session.Get(c.Request, SessionName)
	if err != nil {
		return err
	}
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

func (c *Controller) CurrentUser()(*User, error) {
	session, err := c.App.Session.Get(c.Request, SessionName)
	if err != nil {
		return nil, err
	}

	if id, ok := session.Values["user_id"]; ok {
		 u := &User{Id:id.(int64)}
		 u.FindById(c.App.Database)
		 return u,nil
	}
	return nil,nil
}

func (c *Controller) IsLoggedIn() bool {
	session, _ := c.App.Session.Get(c.Request, SessionName)
	if _, ok := session.Values["user_id"]; ok {
		 return true
	}
	return false
}
