package gobioweb

import (
	"github.com/gorilla/sessions"
	"html/template"
	"net/http"
	dbi "database/sql"
	"github.com/gorilla/context"

)

type AppError struct {
	Error   error
	Message string
	Code    int
	Path string
}

type App struct {
	Template *template.Template
	Session  sessions.Store
	Database *dbi.DB
}

type handlerFunc func(*Controller) *AppError

type Controller struct {
	Handler handlerFunc
	App     *App
	Response http.ResponseWriter
	Request *http.Request
	Flash []interface{}
}


func NewController(h handlerFunc, a *App) *Controller {
	return &Controller{App: a, Handler: h}
}

func (c *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c.Response = w
	c.Request = r
	t := c.App.Template
	//form errors if any
	c.Flash = c.FormErrors()

	if e := c.Handler(c); e != nil {
		e.Path = r.URL.Path
		if e.Code == 500 {
			if errt := t.Lookup("500.tmpl"); errt != nil {
				err := t.ExecuteTemplate(w, "500.tmpl", e)
				if err != nil {
					http.Error(w, "cannot display 500.tmpl template", 500)
				}
			} else {
				tmpT, err := template.New("Internal error template").Parse(srvErrTmpl)
				if err != nil {
					http.Error(w, "cannot parse internal error template", 500)
				}
				err = tmpT.Execute(w, e)
				if err != nil {
					http.Error(w, "cannot execute internal error template", 500)
				}
			}
		} else {
			http.Error(w, e.Message, e.Code)
		}
	}
}

func (c *Controller) GetFromStash(key interface{}) interface{} {
	 if rv := context.Get(c.Request,key); rv != nil {
	 		return rv
	 }
	 return nil
}

func (c *Controller) Stash(key interface{}, value interface{}) {
	 context.Set(c.Request,key,value)
}

func (c *Controller) SetFormErrors(msg string) error {
	 session,err := c.App.Session.Get(c.Request,"form-errors")
	 if err != nil {
	 		return err
	 }
	 session.AddFlash(msg)
	 session.Save(c.Request,c.Response)
	 return nil

}

func (c *Controller) FormErrors() []interface{} {
	 session,_ := c.App.Session.Get(c.Request,"form-errors")
	 if flashes := session.Flashes(); len(flashes) > 0 {
	 		session.Save(c.Request,c.Response)
	 		return flashes
	 }
	 return nil
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

func NotFound(c *Controller) *AppError {
	t := c.App.Template
	if errort := t.Lookup("404.tmpl"); errort != nil {
		err := t.ExecuteTemplate(c.Response, "404.tmpl", c.Request)
		if err != nil {
			return &AppError{Error: err, Code: 500, Message: "Cannot display template"}
		}
	} else {
		tmpT, err := template.New("error template").Parse(errTmpl)
		if err != nil {
			return &AppError{Error: err, Code: 500, Message: "Cannot display template"}
		}
		err = tmpT.Execute(c.Response, c.Request)
		if err != nil {
			return &AppError{Error: err, Code: 500, Message: "Cannot display template"}
		}

	}
	return nil
}

func InternalError( c *Controller) *AppError {
	t := c.App.Template
	if errort := t.Lookup("500.tmpl"); errort != nil {
		err := t.ExecuteTemplate(c.Response, "500.tmpl", c.Request)
		if err != nil {
			return &AppError{Error: err, Code: 500, Message: "Cannot display template"}
		}
	} else {
		tmpT, err := template.New("error template").Parse(errTmpl)
		if err != nil {
			return &AppError{Error: err, Code: 500, Message: "Cannot display template"}
		}
		err = tmpT.Execute(c.Response, c.Request)
		if err != nil {
			return &AppError{Error: err, Code: 500, Message: "Cannot display template"}
		}

	}
	return nil
}
