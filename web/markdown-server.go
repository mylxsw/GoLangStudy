package main

import (
	"bytes"
	//"container/list"
	markdown "github.com/russross/blackfriday"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	//"path/filepath"
	"bufio"
	"strings"
	"text/template"
)

type Article struct {
	Title, Content, Author, Intro, PublishTime string
	ClickTimes                                 int
	Categories                                 []string
	Tags                                       []string
}

type File struct {
	Name  string
	IsDir bool
	Path  string
	Date  string
}

func ParseTemplate(templateName string, data interface{}) []byte {

	var result bytes.Buffer
	tmpl := template.Must(template.ParseFiles(templateName))
	tmpl.Execute(&result, data)

	return result.Bytes()
}

func CreateArticle(sourceFilename string) []byte {
	content_str, meta := ReadMarkdownSource(sourceFilename)

	html := markdown.HtmlRenderer(markdown.HTML_TOC, "", "")
	content := markdown.Markdown([]byte(content_str), html, markdown.EXTENSION_TABLES)

	data := Article{
		Title:       sourceFilename[strings.LastIndex(sourceFilename, "/")+1 : len(sourceFilename)-3],
		Content:     string(content),
		Author:      meta["author"],
		Intro:       meta["intro"],
		PublishTime: meta["publishDate"],
		ClickTimes:  0,
		Categories:  strings.Split(strings.Trim(strings.Trim(meta["categories"], "["), "]"), ","),
		Tags:        strings.Split(strings.Trim(strings.Trim(meta["tags"], "["), "]"), ","),
	}

	return ParseTemplate("templates/article.tmpl", data)
}

func ReadDirectory(directory string) []File {
	fileInfos, _ := ioutil.ReadDir(directory)
	var filelist []File = *new([]File)

	for _, info := range fileInfos {
		if strings.HasPrefix(info.Name(), ".") {
			continue
		}
		if !info.IsDir() && !strings.HasSuffix(info.Name(), ".md") {
			continue
		}
		filelist = append(filelist, File{
			Name:  info.Name(),
			IsDir: info.IsDir(),
			Path:  strings.TrimLeft(strings.TrimRight(directory[len(BASE_DIR):], "/")+"/"+info.Name(), "/"),
			Date:  info.ModTime().Format("2006-01-02"),
		})
	}

	return filelist
}

func ReadMarkdownSource(sourceFileName string) (string, map[string]string) {
	file, _ := os.Open(sourceFileName)
	defer file.Close()

	fileInfo, _ := file.Stat()

	var metas map[string]string = make(map[string]string, 10)
	metas["author"] = "mylxsw"
	metas["publishDate"] = fileInfo.ModTime().Format("2006-01-02")
	metas["intro"] = "..."

	scaner := bufio.NewScanner(file)
	scaner.Scan()
	if scaner.Text() == "---" {
		for scaner.Scan() {
			if scaner.Text() == "---" {
				break
			}

			items := strings.SplitN(scaner.Text(), ":", 2)
			metas[items[0]] = strings.Trim(items[1], " ")
		}
	}

	content := ""
	for scaner.Scan() {
		content += scaner.Text() + "\n"
	}

	return content, metas

}

const (
	BASE_DIR = "/Users/mylxsw/codes/github/mylxsw.github.io/"
)

func main() {

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html; charset=utf-8")

		res := ParseTemplate("templates/list.tmpl", ReadDirectory(BASE_DIR))
		w.Write(res)
	})

	http.HandleFunc("/list/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "text/html; charset=utf-8")

		res := ParseTemplate("templates/list.tmpl", ReadDirectory(BASE_DIR+r.URL.Path[len("/list/"):]))
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
