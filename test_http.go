package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
)

type Libao struct {
	title string
	code  string
	fsize string
}

func FetchLibaos(page int) (libaos []Libao) {
	url := "http://www.libao.so/index.php?page=" + strconv.Itoa(page)
	log.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	//fmt.Println(string(body[:len(body)]))
	body_html := string(body[:len(body)])
	log.Println(body_html)

	re := regexp.MustCompile("<tr>([\\s\\S]*?)</tr>")

	re_title := regexp.MustCompile("<a target=\"_blank\" href=\"(.*?)\">(.*?)</a>")
	re_fsize := regexp.MustCompile("<td style=\"text-align:center;\">(.*?)</td>")
	re_code := regexp.MustCompile("<td style=\"color:red;\">(.*?)</td>")

	tag_trs := re.FindAllString(body_html, -1)
	for i := 0; i < len(tag_trs); i++ {
		tag_title := re_title.FindStringSubmatch(tag_trs[i])

		if len(tag_title) > 2 && tag_title[2] != "" {
			tag_code := re_code.FindStringSubmatch(tag_trs[i])

			if len(tag_code) > 1 && tag_code[1] != "" {

				fsize := "unknown"
				tag_fsize := re_fsize.FindStringSubmatch(tag_trs[i])
				if len(tag_fsize) > 1 && tag_fsize[1] != "" {
					fsize = tag_fsize[1]
				}

				libaos = append(libaos, Libao{
					title: tag_title[2],
					code:  tag_code[1],
					fsize: fsize,
				})
			}
		}
	}
	return
}

func WriteHtml(filename string, libaos []Libao) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	for i := 0; i < len(libaos); i++ {
		record := fmt.Sprintf("%-20s :::: %s :::: %s\n", libaos[i].title, libaos[i].code, libaos[i].fsize)
		file.WriteString(record)
	}
}

func main() {
	for i := 1; i < 3; i++ {
		WriteHtml("result.txt", FetchLibaos(i))
	}
}
