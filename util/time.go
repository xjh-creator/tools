package util

import "time"

var (
	timeLayoutMap = map[string]string{
		"y": "2006",
		"m": "2006-01",
		"d": "2006-01-02",
		"h": "2006-01-02 15",
		"i": "2006-01-02 15:04",
		"s": "2006-01-02 15:04:05",
	}

	weekDay = map[string]int{
		"Sunday":    0,
		"Monday":    1,
		"Tuesday":   2,
		"Wednesday": 3,
		"Thursday":  4,
		"Friday":    5,
		"Saturday":  6,
	}
	loc *time.Location
)

func StringToTime(datetime string) (time.Time, error) {
	curloc, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		panic(err.Error())
	}
	loc = curloc

	result, err := time.ParseInLocation("2006-01-02 15:04:05 -07:00", datetime, loc)
	if err != nil {
		return time.Time{}, err
	}
	return result, nil
}
