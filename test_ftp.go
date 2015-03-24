package main

import (
	ftp "aicode.cc/tools"
	"flag"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"time"
)

var usage = `Usage: test_ftp [options] search_path`

var _addr = flag.String("h", "", "FTP 主机地址:端口号")
var _username = flag.String("u", "anonymous", "登录账户，默认为anonymous")
var _password = flag.String("p", "", "账户密码，默认为空")

// 最大上传大小限制
const BYTE_TO_GB = 1024.0 * 1024.0 * 1024.0
const MAX_UPLOAD_SIZE = BYTE_TO_GB * 150 // 150GB

func usageAndExit(message string) {
	if message != "" {
		fmt.Fprintf(os.Stderr, message)
		fmt.Fprintf(os.Stderr, "\n\n")
	}
	flag.Usage()
	fmt.Fprintf(os.Stderr, "\n")
	os.Exit(1)
}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage)
		fmt.Fprintf(os.Stderr, "\n\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if len(os.Args) != 2 {
		usageAndExit("")
	}

	path_rule := os.Args[1]

	addr := *_addr
	username := *_username
	password := *_password

	if addr == "" {
		usageAndExit("必须提供合法的主机地址!")
	}

	client, err := ftp.NewFtpClient(addr, username, password)
	if err != nil {
		log.Fatal(err)
	}

	// 创建目标目录
	current_date := time.Now().Format("2006/01/02")
	dirname := "/TransportStream/FLVLETV/TVM/" + current_date
	xmldirname := "/TransportStream/XMLLETV/TVM/" + current_date

	client.Mkdir(dirname)
	client.Mkdir(xmldirname)

	var uploaded_size int64 = 0

	files, _ := filepath.Glob(path_rule)
	for _, filename := range files {

		// 判断当前上传文件总大小是否超过本次允许总量，超过则退出
		if uploaded_size >= MAX_UPLOAD_SIZE {
			log.Printf("已经达到最大上传量: %.2f / %.2f (GB)\n", float64(uploaded_size)/BYTE_TO_GB, float64(MAX_UPLOAD_SIZE)/BYTE_TO_GB)
			break
		}

		log.Printf("当前已上传 %.2f GB，正在上传文件: %s\n", float64(uploaded_size)/BYTE_TO_GB, filename)

		// 上传文件信息
		client.ChangeWorkDir(dirname)
		err = client.UploadFile(filename)
		if err != nil {
			log.Printf("文件 %s 上传出错: %s\n", filename, err)
			continue
		}

		// 上传 XML 文件
		xmlfilename := path.Dir(path.Dir(filename)) + "/xml/" + path.Base(filename) + ".xml"

		client.ChangeWorkDir(xmldirname)
		err = client.UploadFile(xmlfilename)
		if err != nil {
			log.Printf("文件%s 上传XML文件失败: %s\n", filename, err)
			continue
		}

		// 增加总文件大小
		_file, _ := os.Open(filename)
		_fileInfo, _ := _file.Stat()
		uploaded_size += _fileInfo.Size()
		_file.Close()

		// 文件上传完成，移动文件到已上传文件夹
		tmp_filename := path.Dir(path.Dir(filename)) + "/uploaded/" + path.Base(filename)
		err = os.Rename(filename, tmp_filename)
		if err != nil {
			log.Fatalf("文件 %s 移动失败: %s\n", filename)
		}
		log.Printf("%s --> %s\n", filename, tmp_filename)
	}

	log.Printf("文件上传完成，本次共上传 %.2f GB.\n", float64(uploaded_size)/BYTE_TO_GB)

}
