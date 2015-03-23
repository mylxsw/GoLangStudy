package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

func main() {
	log.New(os.Stderr, "aicode.cc", log.Ldate)

	mux := http.NewServeMux()
	//mux.Handle("/api/", apiHandler{})
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/" {
			http.NotFound(w, req)
			log.Printf("Request 404: %s", req.URL.Path)
			return
		}
		fmt.Fprintf(w, "Welcome to the home page!")
		log.Printf("Request Ok: %s", req.URL.Path)
	})

	log.Printf(" Server 启动中，监听8000 端口")
	http.ListenAndServe(":8000", mux)

}
