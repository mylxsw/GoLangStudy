package main

import (
	"github.com/tealeg/xlsx"
	"log"
	"strings"
	"time"
)

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
	excelFileName := "/Users/mylxsw/Downloads/222.xlsx"
	xlFile, err := xlsx.OpenFile(excelFileName)
	if err != nil {
		log.Fatalf("文件打开失败: %s", err)
	}

	users := make(map[string]kqUser)
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

	for username, user := range users {
		for kqDate, checkTypes := range user.KqInEveryDays {
			log.Printf("用户 %s 于 %s 考勤 %d 次", username, kqDate, len(checkTypes))

			kqOk := true

			// TODO 考勤检查
			if user.Depart == "客户服务部" {

				if !inArray(checkTypeMorningWork, checkTypes) {
					kqOk = false
					log.Printf("用户 %s 于 %s 早上缺勤", username, kqDate)
				}

				if !inArray(checkTypeAfternoonOffWork, checkTypes) {
					kqOk = false
					log.Printf("用户 %s 于 %s 下午下班缺勤", username, kqDate)
				}

			} else if user.Depart == "秩序维护部" {

			} else if user.Depart == "综合管理部" {

			} else if user.Depart == "维修部" {

				if !inArray(checkTypeMorningWork, checkTypes) {
					kqOk = false
					log.Printf("用户 %s 于 %s 早上缺勤", username, kqDate)
				}

				if !inArray(checkTypeMorningOffWork, checkTypes) {
					kqOk = false
					log.Printf("用户 %s 于 %s 上午下班未打卡", username, kqDate)
				}

				if !inArray(checkTypeAfternoonWork, checkTypes) {
					kqOk = false
					log.Printf("用户 %s 于 %s 下午上班未打卡", username, kqDate)
				}

				if !inArray(checkTypeAfternoonOffWork, checkTypes) {
					kqOk = false
					log.Printf("用户 %s 于 %s 下午下班未打卡", username, kqDate)
				}

			} else if user.Depart == "亳州恒大城物业服务中心" {

			} else {
				log.Printf("用户 %s 所属部门 %s 不再考勤范围内", user.Username, user.Depart)
				continue
			}

			if kqOk {
				log.Printf("用户 %s 与 %s 考勤完整", username, kqDate)
			}
		}
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
