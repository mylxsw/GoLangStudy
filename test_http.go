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

// 礼包内容结构体
type Libao struct {
	title string // 标题
	code  string // 礼包码
	fsize string // 文件大小
}

// 抓取指定页面的所有礼包
func FetchLibaos(page int) (libaos []Libao) {
	// 请求页面信息
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

	// 初始化匹配正则表达式
	re := regexp.MustCompile("<tr>([\\s\\S]*?)</tr>")

	re_title := regexp.MustCompile("<a target=\"_blank\" href=\"(.*?)\">(.*?)</a>")
	re_fsize := regexp.MustCompile("<td style=\"text-align:center;\">(.*?)</td>")
	re_code := regexp.MustCompile("<td style=\"color:red;\">(.*?)</td>")

	// 匹配内容
	tag_trs := re.FindAllString(body_html, -1)
	for i := 0; i < len(tag_trs); i++ {
		// 礼包标题
		tag_title := re_title.FindStringSubmatch(tag_trs[i])

		if len(tag_title) > 2 && tag_title[2] != "" {
			// 礼包码
			tag_code := re_code.FindStringSubmatch(tag_trs[i])

			if len(tag_code) > 1 && tag_code[1] != "" {
				// 礼包大小，如果没有查询到，则使用 unknown
				fsize := "unknown"
				tag_fsize := re_fsize.FindStringSubmatch(tag_trs[i])
				if len(tag_fsize) > 1 && tag_fsize[1] != "" {
					fsize = tag_fsize[1]
				}

				// 追加到返回结果
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

// 将礼包信息写入到文件
func WriteHtml(filename string, libaos []Libao) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	for i := 0; i < len(libaos); i++ {
		record := fmt.Sprintf("%-20s :::: %s :::: %s\n", libaos[i].title, libaos[i].code, libaos[i].fsize)
		_, err := file.WriteString(record)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func main() {
	for i := 1; i < 3; i++ {
		WriteHtml("result.txt", FetchLibaos(i))
	}
}
