package tools

import (
	ftp "github.com/jlaffaye/ftp"
	"log"
	"os"
	"path"
)

// FtpClient 客户端结构体
type FtpClient struct {
	addr       string
	username   string
	password   string
	serverConn *ftp.ServerConn
}

// 创建 FtpClient 实例
func NewFtpClient(addr, username, password string) (*FtpClient, error) {
	// 连接 Ftp 服务器
	serverConn, err := ftp.Connect(addr)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// 登录服务器
	err = serverConn.Login(username, password)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return &FtpClient{addr: addr, username: username, password: password, serverConn: serverConn}, nil
}

// 创建目录
func (client *FtpClient) Mkdir(dirname string) error {
	log.Printf("创建目录 %s .", dirname)

	err := client.serverConn.MakeDir(dirname)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// 改变当前工作目录
func (client *FtpClient) ChangeWorkDir(dirname string) error {
	// 改变当前工作目录
	err := client.serverConn.ChangeDir(dirname)
	if err != nil {
		log.Println(err)
		return err
	}

	currentDir, _ := client.serverConn.CurrentDir()
	log.Println(currentDir)

	return nil
}

// 上传文件
func (client *FtpClient) UploadFile(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		log.Println(err)
		return err
	}
	defer file.Close()

	err = client.serverConn.Stor(path.Base(filename), file)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}
