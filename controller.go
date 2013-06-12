package gobioweb

import (
	"github.com/gorilla/sessions"
	"html/template"
	"net/http"
)

type AppError struct {
	Error   error
	Message string
	Code    int
}

type App struct {
	Template *template.Template
	Session  *sessions.CookieStore
}

type handlerFunc func(http.ResponseWriter, *http.Request, *App) *AppError

type Controller struct {
	Handler handlerFunc
	App     *App
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
				   The requested page <strong> {{.URL.Path}} </strong> cannot be served 
				  </p>
			</div>
	</body>
</html>
`

func (c *Controller) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t := c.App.Template
	if e := c.Handler(w, r, c.App); e != nil {
		if e.Code == 500 {
			if errt := t.Lookup("500.tmpl"); errt != nil {
				err := t.ExecuteTemplate(w, "500.tmpl", r)
				if err != nil {
					http.Error(w, "cannot display 500.tmpl template", 500)
				}
			} else {
				tmpT, err := template.New("Internal error template").Parse(srvErrTmpl)
				if err != nil {
					http.Error(w, "cannot parse internal error template", 500)
				}
				err = tmpT.Execute(w, r)
				if err != nil {
					http.Error(w, "cannot execute internal error template", 500)
				}
			}
		} else {
			http.Error(w, e.Message, e.Code)
		}
	}
}

func NewController(h handlerFunc, a *App) *Controller {
	return &Controller{App: a, Handler: h}
}

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

func NotFound(w http.ResponseWriter, r *http.Request, a *App) *AppError {
	t := a.Template
	if errort := t.Lookup("404.tmpl"); errort != nil {
		err := t.ExecuteTemplate(w, "404.tmpl", r)
		if err != nil {
			return &AppError{Error: err, Code: 500, Message: "Cannot display template"}
		}
	} else {
		tmpT, err := template.New("error template").Parse(errTmpl)
		if err != nil {
			return &AppError{Error: err, Code: 500, Message: "Cannot display template"}
		}
		err = tmpT.Execute(w, r)
		if err != nil {
			return &AppError{Error: err, Code: 500, Message: "Cannot display template"}
		}

	}
	return nil
}

func InternalError(w http.ResponseWriter, r *http.Request, a *App) *AppError {
	t := a.Template
	if errort := t.Lookup("500.tmpl"); errort != nil {
		err := t.ExecuteTemplate(w, "500.tmpl", r)
		if err != nil {
			return &AppError{Error: err, Code: 500, Message: "Cannot display template"}
		}
	} else {
		tmpT, err := template.New("error template").Parse(errTmpl)
		if err != nil {
			return &AppError{Error: err, Code: 500, Message: "Cannot display template"}
		}
		err = tmpT.Execute(w, r)
		if err != nil {
			return &AppError{Error: err, Code: 500, Message: "Cannot display template"}
		}

	}
	return nil
}
