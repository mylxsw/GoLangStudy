package main

import (
	"github.com/tealeg/xlsx"
	"log"
	"strings"
	"time"
	"flag"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

// 重构思路
// 1. 遍历所有用户，找出所涉及到所有的日期，创建标题单元格
// 2. 遍历记录，获取每个用户每天的打卡记录
// 3. 遍历用户，遍历每天的打卡，根据涉及的日期，填写到相关的单元格

var sourceFile = flag.String("source", getCurrentDir() + "/kaoqin.xlsx", "考勤记录文件")
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

	xlFile, err := xlsx.OpenFile(*sourceFile)
	if err != nil {
		log.Fatalf("文件打开失败: %s", err)
	}

	users := make(map[string]kqUser)
	// 第一次遍历，读取excel中的考勤记录，以用户为单位切分
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
			}

			if _, ok := users[username]; !ok {
				users[username] = kqUser{depart, number, username, []time.Time{}, map[string][]uint{}}
			}
			user, _ := users[username]
			user.KqTimes = append(user.KqTimes, signTime)
			users[username] = user
		}
	}

	log.Printf("本次考勤用户共 %d 人", len(users))

	// 遍历用户，计算用户考勤
	for username, user := range users {
		log.Printf("用户 %s 本月共 %d 次考勤", username, len(user.KqTimes))

		kqInEveryDays := map[string][]uint{}
		for _, kqTime := range user.KqTimes {
			log.Printf("用户 %s 打卡: %s", user.Username, kqTime.Format("2006-01-02 15:04"))

			kqDate, checkType := cardChecked(kqTime)

			if _, ok := kqInEveryDays[kqDate]; !ok {
				kqInEveryDays[kqDate] = []uint{}
			}

			kqInEveryDays[kqDate] = append(kqInEveryDays[kqDate], checkType)
		}

		user.KqInEveryDays = kqInEveryDays
		users[username] = user
	}

	// 遍历用户，统计考勤结果
	destExcelFile := xlsx.NewFile()
	sheet, err := destExcelFile.AddSheet("Sheet1")
	if err != nil {
		log.Fatalf("创建Excel文件失败: %s", err)
	}

	// 标题
	row := sheet.AddRow()
	row.AddCell().Value = "编号"
	row.AddCell().Value = "姓名"

	for day  := 0; day < 31; day ++ {
		row.AddCell().Value = strconv.Itoa(day + 1)
	}

	for username, user := range users {
		row := sheet.AddRow()

		cell := row.AddCell()
		cell.Value = user.Number

		cell2 :=row.AddCell()
		cell2.Value = username

		for day := 0; day < 31; day ++ {
			row.AddCell()
		}

		for kqDate, checkTypes := range user.KqInEveryDays {
			log.Printf("用户 %s 于 %s 考勤 %d 次", username, kqDate, len(checkTypes))
			kqOk := true

			var kqStatus []string

			// TODO 考勤检查
			if user.Depart == "客户服务部" {

				if !inArray(checkTypeMorningWork, checkTypes) {
					kqOk = false
					log.Printf("用户 %s 于 %s 早上缺勤", username, kqDate)
					kqStatus = append(kqStatus, "早上缺勤")
				}

				if !inArray(checkTypeAfternoonOffWork, checkTypes) {
					kqOk = false
					log.Printf("用户 %s 于 %s 下午下班缺勤", username, kqDate)
					kqStatus = append(kqStatus, "下午下班缺勤")
				}

			} else if user.Depart == "秩序维护部" {

			} else if user.Depart == "综合管理部" {

			} else if user.Depart == "维修部" {

				if !inArray(checkTypeMorningWork, checkTypes) {
					kqOk = false
					log.Printf("用户 %s 于 %s 早上缺勤", username, kqDate)
					kqStatus = append(kqStatus, "早上缺勤")
				}

				if !inArray(checkTypeMorningOffWork, checkTypes) {
					kqOk = false
					log.Printf("用户 %s 于 %s 上午下班未打卡", username, kqDate)
					kqStatus = append(kqStatus, "上午下班未打卡")
				}

				if !inArray(checkTypeAfternoonWork, checkTypes) {
					kqOk = false
					log.Printf("用户 %s 于 %s 下午上班未打卡", username, kqDate)
					kqStatus = append(kqStatus, "下午上班未打卡")
				}

				if !inArray(checkTypeAfternoonOffWork, checkTypes) {
					kqOk = false
					log.Printf("用户 %s 于 %s 下午下班未打卡", username, kqDate)
					kqStatus = append(kqStatus, "下午下班未打卡")
				}

			} else if user.Depart == "亳州恒大城物业服务中心" {

			} else {
				log.Printf("用户 %s 所属部门 %s 不再考勤范围内", user.Username, user.Depart)
				continue
			}

			currentDate, _ := time.Parse("2006-01-02", kqDate)
			if kqOk {
				log.Printf("用户 %s 与 %s 考勤完整", username, kqDate)
				row.Cells[currentDate.Day() + 1].Value = "正常"
			} else {
				row.Cells[currentDate.Day() + 1].Value = strings.Join(kqStatus, ",")
			}

		}
	}

	err = destExcelFile.Save(*destFile)
	if err != nil {
		log.Fatalf("保存Excel失败: %s", err)
	}

}

func inArray(needle uint, haystack []uint) bool {
	for _, item := range haystack {
		if item == needle {
			return true
		}
	}

	return false
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

func getCurrentDir() string {
	file, _ := exec.LookPath(os.Args[0])
	path := filepath.Dir(file)

	return path
}
