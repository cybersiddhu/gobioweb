package gobioweb

import (
	 "net/http"
)

type controllerError struct {
	Error   error
	Message string
	Code    int
}

type Handler struct {
	Template *template.Template
	Handler  func(http.ResponseWriter, *http.Request, *template.Template) *controllerError
}

func (c *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := c.Handler(w, r, c.Template); e != nil {
		http.Error(w, e.Message, e.Code)
	}
}
