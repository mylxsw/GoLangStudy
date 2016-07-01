package main

import (
	"database/sql"
	"flag"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tealeg/xlsx"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
    "fmt"
)

// 重构思路
// 1. 遍历所有用户，找出所涉及到所有的日期，创建标题单元格
// 2. 遍历记录，获取每个用户每天的打卡记录
// 3. 遍历用户，遍历每天的打卡，根据涉及的日期，填写到相关的单元格

var sourceFile = flag.String("source", getCurrentDir()+"/kaoqin.xlsx", "考勤记录文件")
var destFile = flag.String("dest", time.Now().Format("2006-01-02-15-04-05.xlsx"), "目标文件")

type kqUser struct {
	Depart        string
	Number        string
	Username      string
	KqTimes       []time.Time
	KqInEveryDays map[string][]uint
}

const (
	checkTypeNone = iota
	checkTypeMorningWork
	checkTypeMorningOffWork
	checkTypeAfternoonWork
	checkTypeAfternoonOffWork
)

func main() {
	flag.Parse()

    fmt.Println("Start")

	db, err := sql.Open("sqlite3", "./sqlite")
	if err != nil {
		log.Fatalf("连接数据库失败: %s", err)
		os.Exit(2)
	}
	defer db.Close()

	// 初始化records表
	initRecordTable(db)
	// 将Excel中的内容导入到数据库
	excelToDatabase(*sourceFile, db)

	// 获取Excel中日期的边界
	minDate, maxDate := getDateBorder(db)

    destExcelFile := xlsx.NewFile()
	sheet, err := destExcelFile.AddSheet("Sheet1")
	if err != nil {
		panic(err)
	}

	// 标题
	row := sheet.AddRow()
	row.AddCell().Value = "编号"
	row.AddCell().Value = "姓名"

    for currentDate := minDate; currentDate.Before(maxDate); currentDate = currentDate.Add(24 * time.Hour) {
        row.AddCell().Value = currentDate.Format("1月2日")
    }

    // 查询用户编号列表
    numbers, err := db.Query("select distinct number from records")
    if err != nil {
        panic(err)
    }
    defer numbers.Close()

    stmt, err := db.Prepare("select id, pick_time from records where pick_time > ? and pick_time < ? and number = ?")
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    // 遍历用户编号，处理每个用户的打卡记录
    var pickCountInDay int
    var records *sql.Rows
    for numbers.Next() {
        var userNumber string
        var id int
        var username, depart string
        numbers.Scan(&userNumber)

        row = sheet.AddRow()
        row.AddCell().Value = userNumber

        // 用户信息
        user := db.QueryRow("select id, username, depart from records where number = ?", userNumber)
        user.Scan(&id, &username, &depart)

        row.AddCell().Value = username


        // 每天的打卡记录
        for currentDate := minDate; currentDate.Before(maxDate); currentDate = currentDate.Add(24 * time.Hour) {
            //row.AddCell().Value = currentDate.Format("2006-1-2")
            records, err = stmt.Query(currentDate, currentDate.Add(24 * time.Hour), userNumber)
            if err != nil {
                panic(err)
            }

            pickCountInDay = 0
            for records.Next() {
                var id int
                var pickTime string
                records.Scan(&id, &pickTime)
                fmt.Println(id, pickTime)

                pickCountInDay = pickCountInDay + 1
            }

            row.AddCell().Value = strconv.Itoa(pickCountInDay)

            records.Close()
        }

        // records, err := db.Query("select id, number, username, depart, pick_time from records where number = ?", userNumber)
        // if err != nil {
        //     panic(err)
        // }
        //
        // for records.Next() {
        //     var id int
        //     var number, username, depart, pickTimeStr string
        //     var pickTime time.Time
        //
        //     records.Scan(&id, &number, &username, &depart, &pickTimeStr)
        //     pickTime = parseDate(pickTimeStr)
        //
        //     fmt.Println(id, number, username, depart, pickTime.Format("2006-1-2 15:04"))
        //
        //
        //
        // }
    }



    err = destExcelFile.Save(*destFile)
	if err != nil {
        panic(err)
	}

}

// 数据库中的日期字符串转为time.Time类型
func parseDate(date string) time.Time {
    dateTime, err := time.Parse("2006-1-2 15:04:05+00:00", date)
    if err != nil {
        panic(err)
    }

    return dateTime
}

// 判断整数是否在数组里
func inArray(needle uint, haystack []uint) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}

	return false
}

// Excel文件导入到sqlite数据库
func excelToDatabase(sourceFileName string, db *sql.DB) {
	xlFile, err := xlsx.OpenFile(sourceFileName)
	if err != nil {
		log.Fatalf("文件打开失败: %s", err)
		os.Exit(2)
	}

	stmt, err := db.Prepare("INSERT INTO records (depart, number, username, pick_time) values(?, ?, ?, ?)")
	if err != nil {
		log.Fatalf("准备插入语句出错: %s", err)
		os.Exit(2)
	}
	defer stmt.Close()

	for _, sheet := range xlFile.Sheets {
		for index, row := range sheet.Rows {
			if index == 0 {
				continue
			}

			depart, _ := row.Cells[1].String()
			number, _ := row.Cells[2].String()
			username, _ := row.Cells[3].String()
			timeDate, _ := row.Cells[4].String()
			timeTime, _ := row.Cells[5].String()

			signTime, err := time.Parse("2006/1/2 15:04:05", strings.Split(timeDate, " ")[0]+" "+timeTime+":00")
			if err != nil {
				log.Fatal(err)
				os.Exit(2)
			}

			_, err = stmt.Exec(depart, number, username, signTime)
			if err != nil {
				log.Fatalf("添加打卡记录失败： %s", err)
				os.Exit(2)
			}
		}
	}
}

// 检查用户打卡时间是否为合法的考勤
func cardChecked(kqTime time.Time) (kqDate string, checkType uint) {

	kqDate = kqTime.Format("2006-01-02")

	// 早上上班打卡，打卡时间 7:30-8:30
	kqTime1, _ := time.Parse("2006-01-02 15:04", kqDate+" 7:30")
	kqTime1End, _ := time.Parse("2006-01-02 15:04", kqDate+" 8:30")

	// 中午下班打卡，打卡时间 12:00-12:30
	kqTime2, _ := time.Parse("2006-01-02 15:04", kqDate+" 12:00")
	kqTime2End, _ := time.Parse("2006-01-02 15:04", kqDate+" 12:30")

	// 下午上班打卡，打卡时间 13:30-14:00
	kqTime3, _ := time.Parse("2006-01-02 15:04", kqDate+" 13:30")
	kqTime3End, _ := time.Parse("2006-01-02 15:04", kqDate+" 14:00")

	// 下午下班打卡，打卡时间 17:30-18:30
	kqTime4, _ := time.Parse("2006-01-02 15:04", kqDate+" 17:30")
	kqTime4End, _ := time.Parse("2006-01-02 15:04", kqDate+" 18:30")

	switch {
	case kqTime.After(kqTime1) && kqTime.Before(kqTime1End):
		checkType = checkTypeMorningWork
	case kqTime.After(kqTime2) && kqTime.Before(kqTime2End):
		checkType = checkTypeMorningOffWork
	case kqTime.After(kqTime3) && kqTime.Before(kqTime3End):
		checkType = checkTypeAfternoonWork
	case kqTime.After(kqTime4) && kqTime.Before(kqTime4End):
		checkType = checkTypeAfternoonOffWork
	default:
		checkType = checkTypeNone
	}

	return
}

// 转换考勤标识符为可读的考勤时间
func transformKqForHuman(checkType uint) string {
	var kqExpress string

	switch checkType {
	case checkTypeNone:
		kqExpress = "无效打卡"
	case checkTypeMorningWork:
		kqExpress = "上午上班打卡"
	case checkTypeMorningOffWork:
		kqExpress = "上午下班打卡"
	case checkTypeAfternoonWork:
		kqExpress = "下午上班打卡"
	case checkTypeAfternoonOffWork:
		kqExpress = "下午下班打卡"
	default:
		kqExpress = "无效打卡"
	}

	return kqExpress
}

// 获取当前工作目录
func getCurrentDir() string {
	file, _ := exec.LookPath(os.Args[0])
	path := filepath.Dir(file)

	return path
}

// 初始化records表
func initRecordTable(db *sql.DB) {
	tableCreateSql := `
    DROP TABLE IF EXISTS records;
    CREATE TABLE records (
        id INTEGER PRIMARY KEY,
        depart TEXT,
        number TEXT,
        username TEXT,
        pick_time TEXT,
        pick_type TEXT
    )
    `
	_, err := db.Exec(tableCreateSql)
	if err != nil {
		log.Fatalf("创建表失败: %s", err)
		os.Exit(2)
	}
}

// 获取日期起止边界
func getDateBorder(db *sql.DB) (minDate, maxDate time.Time) {
	rows, err := db.Query("select min(pick_time), max(pick_time)  from records")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	rows.Next()
    var min, max string
	err = rows.Scan(&min, &max)
	if err != nil {
		panic(err)
	}

    minDate, err = time.Parse("2006-1-2", strings.Split(min, " ")[0])
    if err != nil {
        panic(err)
    }

    maxDate, err = time.Parse("2006-1-2", strings.Split(max, " ")[0])
    if err != nil {
        panic(err)
    }

	return
}
