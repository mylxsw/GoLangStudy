package main

import (
	"html/template"
	//"io"
	"log"
	"os"
)

type Context struct {
	Title, Content string
}

func main() {

	con := Context{Title: "Test title", Content: "Hello, Tom"}
	tmpl, err := template.ParseFiles("templates/test.tmpl")
	if err != nil {
		log.Fatalf("parse error: %s\n", err)
	}

	err = tmpl.Execute(os.Stdout, con)
	if err != nil {
		log.Fatalf("execute error: %s\n", err)
	}

	const tmpl_1 = `
---->
{{define "T"}}Hello, {{.}}!{{end}}
<----
`
	tmpl2, err := template.New("test").Parse(tmpl_1)
	if err != nil {
		log.Fatalf("parse error: %s\n", err)
	}

	tmpl2.Execute(os.Stdout, nil)
	tmpl2.ExecuteTemplate(os.Stdout, "T", "Tom")

}
