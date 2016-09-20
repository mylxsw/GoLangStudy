package http

import (
	"net/http"
	"strings"
	"io/ioutil"
	"net/http/cookiejar"
	"code.google.com/p/go.net/proxy"
	"net"
	"time"
	"net/url"
	"os"
	"io"
	"path"
	"crypto/md5"
	"encoding/hex"
)
//
type HttpClient struct {
	http.Client
	proxy       string
	cookie      string
	userAgent   string
	contentType string
}

func NewHttpClient(cookie, proxyUrl string) (*HttpClient, error) {
	cookieJar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}

	dialer, err := proxy.SOCKS5("tcp", proxyUrl, nil, &net.Dialer{
		Timeout: 30 * time.Second,
		KeepAlive: 30 * time.Second,
	})
	if err != nil {
		return nil, err
	}

	transport := &http.Transport{
		Proxy:               nil,
		Dial:                dialer.Dial,
		TLSHandshakeTimeout: 30 * time.Second,
	}

	client := &HttpClient{
		proxy: proxyUrl,
		cookie: cookie,
		userAgent: "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/36.0.1985.125 Safari/537.36",
		contentType: "application/x-www-form-urlencoded",
		Client: http.Client{
			Transport: transport,
			Jar: cookieJar,
		},
	}

	return client, nil
}

func (client *HttpClient) Request(method string, requestUrl string, params url.Values) (responseHtml string, err error) {
	if params == nil {
		params = url.Values{}
	}

	req, err := http.NewRequest(method, requestUrl, strings.NewReader(params.Encode()))
	if err != nil {
		return
	}

	req.Header.Set("User-Agent", client.userAgent)
	req.Header.Set("Content-Type", client.contentType)
	req.Header.Set("Cookie", client.cookie)

	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	responseHtml = string(body)

	return
}

func (client *HttpClient) Get(requestUrl string, params url.Values) (string, error) {
	return client.Request("GET", requestUrl, params)
}

func (client *HttpClient) Post(requestUrl string, params url.Values) (string, error) {
	return client.Request("POST", requestUrl, params)
}

func (client *HttpClient) Download(requestUrl string, savePath string) (string, error) {
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		return "", err
	}

	req.Header.Set("User-Agent", client.userAgent)
	req.Header.Set("Content-Type", client.contentType)

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	os.MkdirAll(savePath, 0777)

	absolutePath := path.Join(savePath, md5enc(requestUrl) + ".jpg")
	file, err := os.Create(absolutePath)
	if err != nil {
		return "", err
	}

	io.Copy(file, resp.Body)

	return absolutePath, nil
}

func md5enc(str string) string {
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(str))
	cipherStr := md5Ctx.Sum(nil)

	return hex.EncodeToString(cipherStr)
}