package main

import (
	"bytes"
	markdown "github.com/russross/blackfriday"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Article struct {
	Title, Content, Author, Intro, PublishTime string
	ClickTimes                                 int
}

func ParseTemplate(templateName string, data interface{}) []byte {

	var result bytes.Buffer
	tmpl := template.Must(template.ParseFiles(templateName))
	tmpl.Execute(&result, data)

	return result.Bytes()
}

func CreateArticle(sourceFilename string) []byte {
	source, _ := os.Open(sourceFilename)
	defer source.Close()
	content, _ := ioutil.ReadAll(source)

	html := markdown.HtmlRenderer(markdown.HTML_TOC, "", "")
	content = markdown.Markdown(content, html, markdown.EXTENSION_TABLES)

	data := Article{
		Title:       sourceFilename[strings.LastIndex(sourceFilename, "/")+1 : len(sourceFilename)-3],
		Content:     string(content),
		Author:      "管宜尧",
		Intro:       "这里是摘要部分",
		PublishTime: "2015-05-27",
		ClickTimes:  1003,
	}

	return ParseTemplate("templates/article.tmpl", data)
}

const (
	BASE_DIR = "/Users/mylxsw/codes/github/mylxsw.github.io/_drafts/"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html; charset=utf-8")

		paths, _ := filepath.Glob(BASE_DIR + "*.md")
		paths2, _ := filepath.Glob(BASE_DIR + "**/*.md")

		var filenames []string = make([]string, len(paths)+len(paths2))

		for index, file := range paths {
			filenames[index] = file[len(BASE_DIR):]
		}

		for index, file := range paths2 {
			filenames[len(paths)+index] = file[len(BASE_DIR):]
		}

		res := ParseTemplate("templates/list.tmpl", filenames)
		w.Write(res)
	})

	http.HandleFunc("/zui/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "templates"+r.URL.Path)
	})

	http.HandleFunc("/show/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html; charset=utf-8")

		filename := r.URL.Path[len("/show/"):]
		log.Printf("request filename: %s\n", filename)

		w.Write(CreateArticle(BASE_DIR + filename))
	})

	http.ListenAndServe(":8080", nil)

}
