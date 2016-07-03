package main

import (
	"database/sql"
	"flag"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tealeg/xlsx"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// 重构思路
// 1. 遍历所有用户，找出所涉及到所有的日期，创建标题单元格
// 2. 遍历记录，获取每个用户每天的打卡记录
// 3. 遍历用户，遍历每天的打卡，根据涉及的日期，填写到相关的单元格

var sourceFile = flag.String("source", getCurrentDir()+"/kaoqin.xlsx", "考勤记录文件")
var destFile = flag.String("dest", time.Now().Format("2006-01-02-15-04-05.xlsx"), "目标文件")

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
	row.AddCell().Value = "部门"

	for currentDate := minDate; currentDate.Before(maxDate); currentDate = currentDate.Add(24 * time.Hour) {
		row.AddCell().Value = currentDate.Format("1月2日") + "(" + parseWeek(currentDate.Weekday()) + ")"
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
	for numbers.Next() {
		var userNumber string
		var id int
		var username, depart string
		numbers.Scan(&userNumber)

		// 用户信息
		user := db.QueryRow("select id, username, depart from records where number = ?", userNumber)
		user.Scan(&id, &username, &depart)

		row = sheet.AddRow()
		row.AddCell().Value = userNumber
		row.AddCell().Value = username
		row.AddCell().Value = depart

		// 每天的打卡记录
		for currentDate := minDate; currentDate.Before(maxDate); currentDate = currentDate.Add(24 * time.Hour) {
			// 查询当前的打卡记录
			row.AddCell().Value = strings.Join(getSignResult(stmt, userNumber, depart, currentDate), "\n")
		}
	}

	// for id, checkType := range dataToUpdate {
	//     _, err = db.Exec("update records set pick_type = ? where id = ?", transformKqForHuman(checkType), id)
	//     if err != nil {
	//         panic(err)
	//     }
	// }

	err = destExcelFile.Save(*destFile)
	if err != nil {
		panic(err)
	}

}

// 转换每周第几天为字符串展示
func parseWeek(dayInWeek time.Weekday) string {
	var weekStr string
	switch dayInWeek {
	case time.Sunday:
		weekStr = "周日"
	case time.Monday:
		weekStr = "周一"
	case time.Tuesday:
		weekStr = "周二"
	case time.Wednesday:
		weekStr = "周三"
	case time.Thursday:
		weekStr = "周四"
	case time.Friday:
		weekStr = "周五"
	case time.Saturday:
		weekStr = "周六"
	}

	return weekStr
}

// 获取签到结果
func getSignResult(stmt *sql.Stmt, userNumber string, depart string, currentDate time.Time) []string {
	results := []string{}

	pickTimesInDay := getPickTimesInDayForUser(stmt, userNumber, currentDate)

	if len(pickTimesInDay) == 0 {
		return []string{"×"}
	}

	if depart == "客户服务部" || depart == "综合管理部" {

		// 早8:30-18:30
		if isSignedInTime("7:30", "8:30", pickTimesInDay, currentDate) || isSignedInTime("17:30", "18:30", pickTimesInDay, currentDate) {
			if !isSignedInTime("7:30", "8:30", pickTimesInDay, currentDate) {
				results = append(results, "早上缺")
			}

			if !isSignedInTime("17:30", "18:30", pickTimesInDay, currentDate) {
				results = append(results, "下午下缺")
			}

		} else if isSignedInTime("9:00", "10:00", pickTimesInDay, currentDate) || isSignedInTime("18:00", "19:00", pickTimesInDay, currentDate) {

			if !isSignedInTime("9:00", "10:00", pickTimesInDay, currentDate) {
				results = append(results, "早上缺")
			}

			if !isSignedInTime("18:00", "19:00", pickTimesInDay, currentDate) {
				results = append(results, "下午下缺")
			}

		} else if isSignedInTime("10:00", "11:00", pickTimesInDay, currentDate) || isSignedInTime("19:00", "20:00", pickTimesInDay, currentDate) {

			if !isSignedInTime("10:00", "11:00", pickTimesInDay, currentDate) {
				results = append(results, "早上缺")
			}

			if !isSignedInTime("19:00", "20:00", pickTimesInDay, currentDate) {
				results = append(results, "下午下缺")
			}

		} else {
			results = append(results, "上下班无，其余有")
		}

	} else if depart == "维保修大队" {

		if !isSignedInTime("7:30", "8:30", pickTimesInDay, currentDate) {
			results = append(results, "早上缺")
		}

		if !isSignedInTime("12:00", "12:30", pickTimesInDay, currentDate) {
			results = append(results, "上午下缺")
		}

		if !isSignedInTime("13:30", "14:00", pickTimesInDay, currentDate) {
			results = append(results, "下午上缺")
		}

		if !isSignedInTime("17:30", "18:30", pickTimesInDay, currentDate) {
			results = append(results, "下午下缺")
		}

	} else {

		if !isSignedInTime("7:30", "8:30", pickTimesInDay, currentDate) {
			results = append(results, "早上缺")
		}

		if !isSignedInTime("17:30", "18:30", pickTimesInDay, currentDate) {
			results = append(results, "下午下缺")
		}
	}

	if len(results) == 0 {
		results = append(results, "√")
	}

	// 标记0-5点的打卡
	if isSignedInTime("00:00", "05:00", pickTimesInDay, currentDate) {
		results = append(results, "☻")
	}

    results = append(results, "(" + strconv.Itoa(len(pickTimesInDay)) + ")")

	return results
}

// 是否是含有指定签到时间
func isSignedInTime(checkTimeStartStr, checkTimeEndStr string, signTimesInDay []time.Time, currentDate time.Time) bool {
	checkTimeStart := parseDate(currentDate.Format("2006-1-2 ") + checkTimeStartStr + ":00+00:00")
	checkTimeEnd := parseDate(currentDate.Format("2006-1-2 ") + checkTimeEndStr + ":00+00:00")
	for _, signTime := range signTimesInDay {
		if signTime.Before(checkTimeEnd) && signTime.After(checkTimeStart) {
			return true
		}
	}

	return false
}

// 数据库中的日期字符串转为time.Time类型
func parseDate(date string) time.Time {
	dateTime, err := time.Parse("2006-1-2 15:04:05+00:00", date)
	if err != nil {
		panic(err)
	}

	return dateTime
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

// 获取某用户某一天的所有考勤时间
func getPickTimesInDayForUser(stmt *sql.Stmt, userNumber string, currentDate time.Time) []time.Time {
	records, err := stmt.Query(currentDate, currentDate.Add(24*time.Hour), userNumber)
	if err != nil {
		panic(err)
	}
	defer records.Close()

	pickTimesInDay := []time.Time{}

	for records.Next() {
		var id int
		var pickTime string
		records.Scan(&id, &pickTime)

		pickTimesInDay = append(pickTimesInDay, parseDate(pickTime))
	}

	return pickTimesInDay
}
