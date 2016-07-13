package main

import (
    "fmt"
    "log"
    "net/http"
    "io/ioutil"
    "encoding/json"
    "os"
)

type message struct {
    CreateTime   string `json:"CreateTime"`
    ValidateCode string `json:"ValidateCode"`
    SMSContent   string `json:"SMSContent"`
}

type response struct {
    Status  string `json:"status"`
    Message string `json:"msg"`
    Data    []message `json:"data"`
}

func main() {
    if len(os.Args) == 1 {
        fmt.Println("缺少要查询的手机号码")
        os.Exit(0)
    }

    tel := os.Args[1]

    url := "http://10.1.20.54/Service.ashx?action=SmsQuery&mobile=" + tel

    client := &http.Client{}
    req, _ := http.NewRequest(http.MethodGet, url, nil)
    req.Host = "t3devhelp.frontpay.cn"

    res, _ := client.Do(req)

    content, err := ioutil.ReadAll(res.Body)
    if err != nil {
        log.Fatal(err)
    }

    var contentJSON response
    err = json.Unmarshal(content[1:len(content) - 1], &contentJSON)
    if err != nil {
        log.Fatal(err)
    }

    if len(contentJSON.Data) == 0 {
        fmt.Println("没有查询到相关短信")
    }

    for _, data := range contentJSON.Data {
        fmt.Printf("%s    %s\n", data.CreateTime, data.SMSContent)
    }
}
