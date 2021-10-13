package pkg

import (
	"regexp"
	"strconv"
	"strings"
)

// ParseFavAndDay
// "6人关注 / 5天以前发布"
// "6人关注 / 5个月以前发布"
// "6人关注 / 五年以前发布"
func ParseFavAndDay(text string) (int, int) {
	numMap := map[string]int{
		"一": 1,
		"二": 2,
		"三": 3,
		"四": 4,
		"五": 5,
		"六": 6,
		"七": 7,
		"八": 8,
		"九": 9,
	}
	numRe, _ := regexp.Compile(`\d+`)
	dayRe, _ := regexp.Compile(`\d+天`)
	monthRe, _ := regexp.Compile(`\d+个月`)
	yearRe, _ := regexp.Compile(`^.+年`)

	favString := string(numRe.Find([]byte(strings.Split(text, "/")[0]))) //正则表达式用来匹配数字
	fav, _ := strconv.Atoi(favString)
	date := strings.Split(text, "/")[1]

	dayString := numRe.FindString(dayRe.FindString(date))
	monthString := numRe.FindString(monthRe.FindString(date))
	yearString := strings.ReplaceAll(strings.ReplaceAll((yearRe.FindString(date)), "年", ""), " ", "")

	if dayString != "" {
		day, _ := strconv.Atoi(dayString)
		return fav, day

	}

	if monthString != "" {
		month, _ := strconv.Atoi(monthString)
		return fav, month * 30
	}

	return fav, numMap[yearString] * 360
}
