package main

import (
	"log"
	"os"
	"text/template"
)

func CreateTestFile(tmpls map[string]string) {
	for name, value := range tmpls {
		file, err := os.Create(name)
		if err != nil {
			panic(err)
		}

		file.Write([]byte(value))
		file.Close()
	}
}

func main() {
	templates := map[string]string{
		"template/T0.tmpl": `T0 invokes T1: ({{template "T1"}}`,
		"template/T1.tmpl": `{{define "T1"}}T1 invokes T2: ({{template "T2"}}){{end}}`,
		"template/T2.tmpl": `{{define "T2"}}This is T2{{end}}`,
	}
	os.MkdirAll("template/", os.ModeDir|os.ModePerm)
	CreateTestFile(templates)

	defer os.RemoveAll("template/")

	tmpl := template.Must(template.ParseGlob("template/*.tmpl"))

	err := tmpl.Execute(os.Stdout, nil)
	if err != nil {
		log.Fatalf("template execution: %s\n", err)
	}

}
