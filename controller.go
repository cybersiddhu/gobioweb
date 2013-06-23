package gobioweb

import (
	"github.com/gorilla/context"
	"html/template"
	"net/http"
	"path/filepath"
)

type handlerFunc func(*Controller) *AppError

type Controller struct {
	Handler  handlerFunc
	App      *App
	Response http.ResponseWriter
	Request  *http.Request
	Flash    []interface{}
}

func NewController(h handlerFunc, a *App) *Controller {
	return &Controller{App: a, Handler: h}
}

func (c *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.Request = r
	c.Response = w

	//	form errors if any
	flash, err := c.FormErrors()
	if err != nil {
		c.ServerError(err)
		return
	}

	if flash != nil {
		c.Flash = flash
	} else {
		flash, err := c.Notices()
		if err != nil {
			c.ServerError(err)
		}
		if flash != nil {
			c.Flash = flash
		}
	}

	if e := c.Handler(c); e != nil {
		e.Path = r.URL.Path
		if e.Code == http.StatusInternalServerError {
			c.ServerError(e)
		} else {
			http.Error(c.Response, e.Message, e.Code)
		}
	}
}

func (c *Controller) ServerError(e *AppError) {
	w := c.Response
	t := c.App.Template
	tmpl, err := template.New("500").ParseFiles(filepath.Join(t.Path, "500.tmpl"))
	if err != nil {
		tmpT, err := template.New("Internal error template").Parse(srvErrTmpl)
		if err != nil {
			http.Error(w, "cannot parse internal error template", 500)
		}
		err = tmpT.Execute(w, e)
		if err != nil {
			http.Error(w, "cannot execute internal error template", 500)
		}
	} else {
		err := tmpl.Execute(w, e)
		if err != nil {
			http.Error(w, "cannot display 500.tmpl template", 500)
		}
	}
}

func NotFound(c *Controller) *AppError {
	t := c.App.Template
	tmpl, err := template.New("400").ParseFiles(filepath.Join(t.Path, "400.tmpl"))
	if err != nil {
		tmpT, err := template.New("error template").Parse(errTmpl)
		if err != nil {
			return &AppError{Error: err, Code: 500, Message: "Cannot display template"}
		}
		err = tmpT.Execute(c.Response, c.Request)
		if err != nil {
			return &AppError{Error: err, Code: 500, Message: "Cannot display template"}
		}

	} else {
		err := tmpl.Execute(c.Response, c.Request)
		if err != nil {
			return &AppError{Error: err, Code: 500, Message: "Cannot display template"}
		}
	}
	return nil
}

func (c *Controller) GetFromStash(key interface{}) interface{} {
	if rv := context.Get(c.Request, key); rv != nil {
		return rv
	}
	return nil
}

func (c *Controller) Stash(key interface{}, value interface{}) {
	context.Set(c.Request, key, value)
}

func (c *Controller) SetFormErrors(msg string) error {
	session, err := c.App.Session.Get(c.Request, "flashes")
	if err != nil {
		return err
	}
	session.AddFlash("form-errors", msg)
	session.Save(c.Request, c.Response)
	return nil

}

func (c *Controller) FormErrors() ([]interface{}, *AppError) {
	session, err := c.App.Session.Get(c.Request, "flashes")
	if err != nil {
		return nil, &AppError{
			Code:    http.StatusInternalServerError,
			Error:   err,
			Message: err.Error(),
		}
	}
	if flashes := session.Flashes("form-errors"); len(flashes) > 0 {
		session.Save(c.Request, c.Response)
		return flashes, nil
	}
	return nil, nil
}

func (c *Controller) Notices(values ...string) ([]interface{}, *AppError) {
	session, err := c.App.Session.Get(c.Request, "flashes")
	if err != nil {
		return nil, &AppError{
			Code:    http.StatusInternalServerError,
			Error:   err,
			Message: err.Error(),
		}
	}

	//set function
	if len(values) > 0 {
		for _, val := range values {
			session.AddFlash("notices", val)
			session.Save(c.Request, c.Response)
		}
		return nil, nil
	}

	//get function
	if flashes := session.Flashes("notices"); len(flashes) > 0 {
		session.Save(c.Request, c.Response)
		return flashes, nil
	}
	return nil, nil
}

const srvErrTmpl = `
<!DOCTYPE html>
<html>
	<head>
		<title> Internal server error 500 </title>
		<link
		href="//netdna.bootstrapcdn.com/twitter-bootstrap/2.3.2/css/bootstrap-combined.min.css"
		rel="stylesheet">
		<style type="text/css">
		  .center {
		  					text-align: center; 
		  					margin-left: auto; 
		  					margin-right: auto; 
		  					margin-bottom: auto; 
		  					margin-top: auto;
		  				}
		  h2.red {
				 color: red;
		  }
		 </style>
	</head>
	<body>
			<div class="hero-unit center">
					<h1> Internal Server Error </h1>
					<h2 class="red"> Error 500 </h2>
				  <p class="lead"> 
				   The requested page <strong> {{.Path}} </strong> cannot be served 
				  </p>
				  <p class="lead"> 
						Reason: {{.Message}}  
				  </p>
			</div>
	</body>
</html>
`

const errTmpl = `
<!DOCTYPE html>
<htm>
	<head>
		<title> Page not found 404 </title>
		<link
		href="//netdna.bootstrapcdn.com/twitter-bootstrap/2.3.2/css/bootstrap-combined.min.css"
		rel="stylesheet">
		<style type="text/css">
		  .center {
		  					text-align: center; 
		  					margin-left: auto; 
		  					margin-right: auto; 
		  					margin-bottom: auto; 
		  					margin-top: auto;
		  				}
		  h2.red {
				 color: red;
		  }
		 </style>
	</head>
	<body>
			<div class="hero-unit center">
					<h1> Page Not Found </h1>
					<h2 class="red"> Error 404 </h2>
				  <p class="lead"> 
				   The requested page <strong> {{.URL.Path}} </strong> could not be found 
				  </p>
			</div>
	</body>
</html>
`
