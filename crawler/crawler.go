package main

import (
	"fmt"
	"os"

	"github.com/crufter/goquery"
	"aicode.cc/crawler/http"
	"log"
	"sync"
)

type imageLink struct {
	title string
	url string
}

func producer(httpClient *http.HttpClient, imageQueue chan imageLink) {
	html, err := httpClient.Get("http://ck101.com/forum-80-1.html", nil)
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
	nodes, _ := goquery.ParseString(html)
	linkNodes := nodes.Find(".blockTitle a.xst")
	if linkNodes.Length() > 0 {
		var wg sync.WaitGroup
		wg.Add(linkNodes.Length())
		for i := 0; i < linkNodes.Length(); i ++ {
			node := linkNodes.Eq(i)
			href := node.Attr("href")
			title := node.Attr("title")
			log.Printf("找到链接 %s : %s", title, href)
			go func(link string) {
				defer wg.Done()
				resp, err := httpClient.Get(link, nil)
				if err != nil {
					log.Printf("请求出错: %s", err.Error())
					return
				}
				nodes, _ := goquery.ParseString(resp)
				imageNodes := nodes.Find("img.zoom")
				for i := 0; i < imageNodes.Length(); i ++ {
					imageUrl := imageNodes.Eq(i).Attr("file")
					log.Printf("%s -> %s", title, imageUrl)

					imageQueue <- imageLink {
						title: title,
						url: imageUrl,
					}
				}
			}(href)
		}
		wg.Wait()
	}

	log.Println("链接抓取完毕")
}

func consumers(httpClient *http.HttpClient, imageQueue chan imageLink) {
	for image := range imageQueue {
		log.Printf("Downloading %s ...", image)

		savePath, err := httpClient.Download(image.url, "/Users/mylxsw/Downloads/" + image.title)
		if err != nil {
			log.Printf("文件 %s 下载失败: %s", savePath, err.Error())
			continue
		}

		log.Printf("文件 %s 下载完成", savePath)
	}


}

func main() {
	cookie := "datr=YGSgV2uQVaK3f_y2nfCKNmAU; sb=E3jbVxSdIkM-1ZBsqnhr6Ac0; c_user=100003957548908; xs=40%3AOAi_u2h4Vh1yow%3A2%3A1474000915%3A11280; fr=0Ufsnv4yUWNXaO95L.AWUyLOyuErUJuqLdlv7QVBGcHww.BXmB_q.Oz.AAA.0.0.BX23gT.AWVsfaMM; csm=2; s=Aa7-Enb1NEt3Z7Yv.BX23gT; pl=y; lu=ghh-ZJIdHjmUiQet1stnKABg"

	httpClient, err := http.NewHttpClient(cookie, "127.0.0.1:1080")
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	var done chan struct{}
	var imageQueue chan imageLink = make(chan imageLink, 10)

	go producer(httpClient, imageQueue)
	go consumers(httpClient, imageQueue)

	<-done

	//links := nodes.Find(".blockTitle a.xst").Attrs("href")
	//for _, link := range links {
	//
	//	subHtml, _ := httpClient.Get(link, nil)
	//	subNodes, _ := goquery.ParseString(subHtml)
	//	fmt.Println("thanks " + strconv.Itoa(subNodes.Find("a.lockThankBtn").Length()))
	//	if subNodes.Find("a.lockThankBtn").Length() > 0 {
	//		thankUrl := "http://ck101.com/" + subNodes.Find("a.lockThankBtn").Attr("href") + "&infloat=yes&extra=&thankssubmit=yes"
	//		_, err = httpClient.Post(thankUrl, nil)
	//		if err != nil {
	//			panic(err)
	//		}
	//	}
	//
	//	imageLinks := subNodes.Find("img.zoom").Attrs("file")
	//
	//	for _, imageLink := range imageLinks {
	//		fmt.Println(imageLink)
	//	}
	//}
}
