package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"

	"time"

	"strings"

	"encoding/json"

	mgo "gopkg.in/mgo.v2"
)

type CacheLog struct {
	Time   time.Time
	Type   string
	Key    string
	Finger string
}

func main() {
	outputs := make(chan string, 100)

	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	c := session.DB("test").C("logs")

	cmd := exec.Command("/bin/sh", "-x", "-c", "tail -n 0 -f /data/logs/scm-site.log | grep --line-buffered CACHE_")
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		reader := bufio.NewReader(stdoutPipe)
		for {
			line, err := reader.ReadString('\n')
			if err != nil || io.EOF == err {
				break
			}

			outputs <- fmt.Sprintf("%s", line)
		}
	}()

	go func() {
		location, _ := time.LoadLocation("Local")
		for output := range outputs {
			log.Printf("-> %s", output)

			logTime, err := time.ParseInLocation("2006-01-02 15:04:05", output[1:20], location)
			if err != nil {
				log.Printf("Error: %s", err)
				continue
			}

			logRec := strings.Split(output, " ")
			cacheLog := CacheLog{
				Time: logTime,
				Type: logRec[3],
			}
			if err = json.Unmarshal([]byte(logRec[4]), &cacheLog); err != nil {
				log.Printf("Error: %s", err)
				continue
			}

			if err = c.Insert(&cacheLog); err != nil {
				log.Printf("Error: %s", err)
				continue
			}
		}
	}()

	if err = cmd.Start(); err != nil {
		log.Fatal(err)
	}

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
}
