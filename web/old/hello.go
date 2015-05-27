package main

import (
	"aicode.cc/web/session"
	_ "aicode.cc/web/session/memory"
	"crypto/md5"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var globalSessions *session.Manager

func init() {
	globalSessions, _ = session.NewManager("memory", "gosessionid", 3600)
	go globalSessions.GC()
}

func sayHelloName(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	fmt.Println(req.Form)
	fmt.Println("path", req.URL.Path)
	fmt.Println("scheme", req.URL.Scheme)
	fmt.Println(req.Form["url_long"])

	for k, v := range req.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello Browser")
}

func login(w http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	fmt.Println("method:", req.Method)

	session := globalSessions.SessionStart(w, req)

	if req.Method == "GET" {
		currentime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(currentime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))
		session.Set("token", token)

		t, _ := template.ParseFiles("views/login.gtpl")
		t.Execute(w, token)
	} else {
		token := req.Form.Get("token")
		fmt.Println("token:", token)
		saved_token := session.Get("token")
		if saved_token == token {
			fmt.Println("TOKEN VALID")
		} else {
			fmt.Println("TOKEN INVALID")
		}

		fmt.Println("username:", req.Form["username"])
		fmt.Println("password:", req.Form["password"])

		template.HTMLEscape(w, []byte(req.Form.Get("username")))
	}
}

func main() {
	//http.HandleFunc("/", sayHelloName)
	http.HandleFunc("/login", login)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
