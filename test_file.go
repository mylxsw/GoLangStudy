package main

import (
	"log"
	"os"
	"path"
)

func main() {
	// 检查是否提供了合法参数
	arg_num := len(os.Args)
	if arg_num < 2 {
		log.Fatal("Usage: test_file source dest")
	}

	// 读取以源文件名及目标文件名
	filename := os.Args[1]
	var bak_filename string

	if arg_num >= 3 {
		bak_filename = os.Args[2]
	} else {
		// 如果没有提供目标文件名，则自动设置目标文件名为"源文件名_bak.源文件扩展名"
		ext := path.Ext(filename)
		bak_filename = filename[:len(filename)-len(ext)] + "_bak" + ext

		log.Printf("没有提供目标文件名，使用默认 %s\n", bak_filename)
	}

	// 打开源文件
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// 创建目标文件
	bak_file, err := os.Create(bak_filename)
	if err != nil {
		log.Fatal(err)
	}
	defer bak_file.Close()

	// 创建读取文件信息的缓冲区
	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("源文件大小: %.2fKB\n", float64(float64(fileInfo.Size())/1024.0))

	log.Printf("开始复制文件 %s --> %s", filename, bak_filename)
	buffer := make([]byte, 1024)
	counter := 0
	for {
		n, _ := file.Read(buffer)
		if 0 == n {
			break
		}

		bak_file.Write(buffer[:n])

		// 计数器，计算文件复制进度
		log.Printf("进度:%.2f%%", float64(float64(1024*counter)/float64(fileInfo.Size()))*100)
		counter++
	}

	log.Println("复制完成!")

}
