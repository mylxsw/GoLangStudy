package main

import (
	"fmt"
	"net/http"
	"path/filepath"
)

type HttpHandler struct{}

func (handler *HttpHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Header().Add("Content-Type", "text/html; charset=utf-8")
	if req.URL.RequestURI() == "/" {
		files, err := filepath.Glob("/Users/mylxsw/codes/github/mylxsw.github.io/_drafts/**/*.md")
		if err != nil {
			fmt.Fprintln(w, err.Error())
			return
		}

		for _, file := range files {
			fmt.Fprintf(w, "%s<br />", file)
		}
	}
}

func main() {
	http.ListenAndServe(":8080", &HttpHandler{})
}
