package main

import (
	"database/sql"
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/tealeg/xlsx"
)

// 重构思路
// 1. 遍历所有用户，找出所涉及到所有的日期，创建标题单元格
// 2. 遍历记录，获取每个用户每天的打卡记录
// 3. 遍历用户，遍历每天的打卡，根据涉及的日期，填写到相关的单元格

var sourceFile = flag.String("source", getCurrentDir()+"/kaoqin.xlsx", "考勤记录文件")
var destFile = flag.String("dest", time.Now().Format("2006-01-02-15-04-05.xlsx"), "目标文件")

func main() {
	flag.Parse()

	log.Print("程序处理中...")

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

	// 创建单元格样式
	dangerStyle := xlsx.NewStyle()
	dangerStyle.Fill = *xlsx.NewFill("solid", "00FF0000", "FF000000")

	warningStyle := xlsx.NewStyle()
	warningStyle.Fill = *xlsx.NewFill("solid", "00FFB61E", "FF000000")

	infoStyle := xlsx.NewStyle()
	infoStyle.Fill = *xlsx.NewFill("solid", "0044CEF6", "FF000000")

	// 标题
	row := sheet.AddRow()
	row.AddCell().Value = "部门"
	row.AddCell().Value = "编号"
	row.AddCell().Value = "姓名"

	for currentDate := minDate; currentDate.Before(maxDate.Add(24*time.Hour - 1*time.Nanosecond)); currentDate = currentDate.Add(24 * time.Hour) {
		cell := row.AddCell()
		weekStr := parseWeek(currentDate.Weekday())
		cell.Value = currentDate.Format("1月2日") + "(" + weekStr + ")"

		if weekStr == "周六" || weekStr == "周日" {
			cell.SetStyle(infoStyle)
		}
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
		row.AddCell().Value = depart
		row.AddCell().Value = userNumber
		row.AddCell().Value = username

		// 每天的打卡记录
		for currentDate := minDate; currentDate.Before(maxDate.Add(24*time.Hour - 1*time.Nanosecond)); currentDate = currentDate.Add(24 * time.Hour) {
			// 查询当前的打卡记录
			signResult := strings.Join(getSignResult(stmt, userNumber, depart, currentDate), " ")
			cell := row.AddCell()
			cell.Value = signResult
			if signResult == "O" {
				cell.SetStyle(dangerStyle)
			} else if !(signResult == "C" ||
				signResult == "W" ||
				signResult == "√" ||
				signResult == "A" ||
				signResult == "B" ||
				signResult == "A/B" ||
				signResult == "B/A") {
				cell.SetStyle(warningStyle)
			}
			log.Printf("用户 %s - %s [%s] %s 打卡记录: %s", depart, username, userNumber, currentDate.Format("2006-1-2"), signResult)
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

	if len(pickTimesInDay) == 0 && (depart != "维修部" && depart != "秩序维护部") {
		return []string{"O"}
	}

	if depart == "客户服务部" || depart == "综合管理部" {

		// 需要同时满足在11点前有打卡记录,同时,在17:30之后有打卡记录,同时上班总时间超过9小时
		// 如果缺勤,时间可能不准确,比如上午缺勤,无法根据下午下班时间判断是上的什么班
		morningOk := isSignedBefore("11:00", pickTimesInDay, currentDate)
		afternoonOk := isSignedAfter("17:30", pickTimesInDay, currentDate)
		isEnough := calTimeDiffInHours(pickTimesInDay) >= 9

		if !morningOk {
			results = append(results, "缺上")
		}

		if !afternoonOk {
			results = append(results, "缺下")
		}

		if morningOk && !isEnough && afternoonOk {
			results = append(results, "早退")
		}

		// if isSignedContainRange("11:00", "20:00", pickTimesInDay, currentDate) {
		// 	results = append(results, "W")
		// } else if isSignedContainRange("10:00", "19:00", pickTimesInDay, currentDate) {
		// 	results = append(results, "C")
		// }

	} else if depart == "维保修部/维保修大队（亳州恒大城）" || depart == "维保修大队" || depart == "维保修大队（亳州恒大城）" {

		if !isSignedBefore("8:30", pickTimesInDay, currentDate) {
			results = append(results, "缺上")
		}

		if !isSignedInTime("12:00", "12:59", pickTimesInDay, currentDate) {
			results = append(results, "缺上下")
		}

		if !isSignedInTime("13:00", "14:00", pickTimesInDay, currentDate) {
			results = append(results, "缺下上")
		}

		if !isSignedAfter("17:30", pickTimesInDay, currentDate) {
			results = append(results, "缺下")
		}

	} else if depart == "维修部" {
		pickTimesInNextDay := getPickTimesInDayForUser(stmt, userNumber, currentDate.Add(24*time.Hour))
		dateAtNextDay := currentDate.Add(24 * time.Hour)
		//if isSignedBefore("8:30", pickTimesInDay, currentDate) {
		//	// 白班
		//	if !isSignedAfter("17:30", pickTimesInDay, currentDate) {
		//		results = append(results, "异常")
		//	}
		//} else if isSignedInTime("12:00", "17:30", pickTimesInDay, currentDate) {
		//	// 晚班
		//	if !isSignedInTime("8:30", "10:30", pickTimesInNextDay, currentDate.Add(24*time.Hour)) {
		//		results = append(results, "异常")
		//	}
		//} else {
		//	results = append(results, "异常")
		//}

		if isSignedInTime("7:00", "8:30", pickTimesInDay, currentDate) {
			if isSignedInTime("17:00", "19:00", pickTimesInDay, currentDate) {
				results = append(results, "√")
			} else {
				results = append(results, "缺下")
			}
		} else if isSignedInTime("14:00", "17:30", pickTimesInDay, currentDate) {
			if isSignedInTime("8:30", "10:00", pickTimesInNextDay, dateAtNextDay) {
				results = append(results, "√")
			} else {
				results = append(results, "缺下")
			}
		} else if isSignedInTime("17:30", "19:00", pickTimesInDay, currentDate) {
			results = append(results, "缺上")
		} else if isSignedInTime("8:30", "10:00", pickTimesInNextDay, dateAtNextDay) {
			results = append(results, "缺上")
		} else {
			results = append(results, "O")
		}

	} else if depart == "秩序维护部" {
		pickTimesInNextDay := getPickTimesInDayForUser(stmt, userNumber, currentDate.Add(24*time.Hour))
		dateAtNextDay := currentDate.Add(24 * time.Hour)
		if isSignedInTime("6:00", "7:00", pickTimesInDay, currentDate) {
			if isSignedInTime("19:00", "23:00", pickTimesInDay, currentDate) {
				if isSignedInTime("0:00", "01:00", pickTimesInNextDay, dateAtNextDay) {
					if isSignedInTime("07:00", "08:00", pickTimesInNextDay, dateAtNextDay) {
						results = append(results, "A/B")
					} else {
						results = append(results, "A")
					}
				} else {
					results = append(results, "A")
				}
			} else {
				results = append(results, "缺下")
			}
		} else if isSignedInTime("15:00", "19:00", pickTimesInDay, currentDate) {
			if isSignedInTime("01:00", "03:00", pickTimesInNextDay, dateAtNextDay) {

				if isSignedInTime("7:00", "8:00", pickTimesInNextDay, dateAtNextDay) {
					results = append(results, "B")
				} else {
					results = append(results, "B/A")
				}

			} else if isSignedInTime("7:00", "12:00", pickTimesInNextDay, dateAtNextDay) {
				results = append(results, "B")
			} else {
				results = append(results, "缺下")
			}

		} else if isSignedInTime("19:00", "23:00", pickTimesInDay, currentDate) {
			results = append(results, "缺上")
		} else {
			if isSignedInTime("0:00", "01:00", pickTimesInNextDay, dateAtNextDay) {
				if isSignedInTime("7:00", "8:00", pickTimesInNextDay, dateAtNextDay) {
					results = append(results, "B/A")
				} else {
					results = append(results, "异常")
				}
			} else {
				if isSignedInTime("7:00", "12:00", pickTimesInNextDay, dateAtNextDay) {
					results = append(results, "缺上")
				} else {
					results = append(results, "O")
				}
			}
		}

	} else {
		morningOk := isSignedBefore("8:30", pickTimesInDay, currentDate)
		afternoonOk := isSignedAfter("17:30", pickTimesInDay, currentDate)

		if !morningOk && afternoonOk {
			results = append(results, "缺上")
		}

		if !afternoonOk && morningOk {
			results = append(results, "缺下")
		}

		if !morningOk && !afternoonOk {
			results = append(results, "早退")
		}
	}

	if len(results) == 0 {
		results = append(results, "√")
	}

	// 标记0-5点的打卡
	//if isSignedInTime("00:00", "05:00", pickTimesInDay, currentDate) {
	//    results = append(results, "☻")
	//}

	//results = append(results, "(" + strconv.Itoa(len(pickTimesInDay)) + ")")

	return results
}

// 是否是含有指定签到时间
func isSignedInTime(checkTimeStartStr, checkTimeEndStr string, signTimesInDay []time.Time, currentDate time.Time) bool {
	checkTimeStart := parseDate(currentDate.Format("2006-1-2 ") + checkTimeStartStr + ":00+00:00")
	checkTimeEnd := parseDate(currentDate.Format("2006-1-2 ") + checkTimeEndStr + ":00+00:00")
	for _, signTime := range signTimesInDay {
		if (signTime.Before(checkTimeEnd) && signTime.After(checkTimeStart)) ||
			signTime.Equal(checkTimeEnd) || signTime.Equal(checkTimeStart) {
			return true
		}
	}

	return false
}

// 检查在指定时间之前是否有打卡记录
func isSignedBefore(checkTimeStr string, signTimesInDay []time.Time, currentDate time.Time) bool {
	checkTime := parseDate(currentDate.Format("2006-1-2 ") + checkTimeStr + ":00+00:00")
	for _, signTime := range signTimesInDay {
		if signTime.Before(checkTime) || signTime.Equal(checkTime) {
			return true
		}
	}

	return false
}

// 检查在指定时间之后是否有打卡记录
func isSignedAfter(checkTimeStr string, signTimesInDay []time.Time, currentDate time.Time) bool {
	checkTime := parseDate(currentDate.Format("2006-1-2 ") + checkTimeStr + ":00+00:00")
	for _, signTime := range signTimesInDay {
		if signTime.After(checkTime) || signTime.Equal(checkTime) {
			return true
		}
	}

	return false
}

// 判断打卡时间是否包含指定的时间范围
func isSignedContainRange(startTimeStr, endTimeStr string, signTimesInDay []time.Time, currentDate time.Time) bool {
	return isSignedBefore(startTimeStr, signTimesInDay, currentDate) && isSignedAfter(endTimeStr, signTimesInDay, currentDate)
}

// 计算一天中最早打卡和最晚打卡的时间差
func calTimeDiffInHours(signTimesInDay []time.Time) float64 {

	if len(signTimesInDay) == 0 {
		return 0
	}

	signTimeFirst := getFirstSignTime(signTimesInDay)
	signTimeLast := getLastSignTime(signTimesInDay)

	timeDiff := signTimeLast.Sub(signTimeFirst)

	return timeDiff.Hours()
}

// 获取一天中最早的打卡时间
func getFirstSignTime(signTimesInDay []time.Time) time.Time {
	signTimeFirst := signTimesInDay[0]
	for _, signTime := range signTimesInDay[1:] {
		if signTimeFirst.After(signTime) {
			signTimeFirst = signTime
		}
	}

	return signTimeFirst
}

// 获取一天中最晚的打卡时间
func getLastSignTime(signTimesInDay []time.Time) time.Time {
	signTimeLast := signTimesInDay[0]
	for _, signTime := range signTimesInDay[1:] {
		if signTimeLast.Before(signTime) {
			signTimeLast = signTime
		}
	}

	return signTimeLast
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

	log.Print("打卡数据导入临时数据库...")

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

			depart, _ := row.Cells[0].String()
			number, _ := row.Cells[1].String()
			username, _ := row.Cells[2].String()

			if depart == "" {
				break
			}

			timeDate, _ := row.Cells[3].String()
			timeTime, _ := row.Cells[4].String()

			log.Printf("Raw: %s %s %s %s %s", depart, number, username, timeDate, timeTime)

			if strings.Count(timeTime, ":") == 2 {
				timeTime = timeTime[:len(timeTime)-3]
			}

			signTime, err := time.Parse("2006/1/2 15:04:05", strings.Split(timeDate, " ")[0]+" "+timeTime+":00")
			if err != nil {
				signTime, err = time.Parse("1-2-06 15:04:05", strings.Split(timeDate, " ")[0]+" "+timeTime+":00")
				if err != nil {
					signTime, err = time.Parse("1/2/06 15:04:05", strings.Split(timeDate, " ")[0]+" "+timeTime+":00")
					if err != nil {
						signTime, err = time.Parse("2006-1-2 15:04:05", strings.Split(timeDate, " ")[0]+" "+timeTime+":00")
						if err != nil {
							log.Fatal(err)
							os.Exit(2)
						}
					}
				}
			}

			log.Printf("导入用户 %s - %s [%s] at %s", depart, username, number, signTime.Format("2006-1-2 15:04:05"))

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
	log.Printf("初始化临时数据库: %s", tableCreateSql)

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
