package util

// getFebDayNumsByYear 通过年份获取该二月的天数
func getFebDayNumsByYear(year int) int {
	if (year%4 == 0 && year%100 != 0) || year%400 == 0 { // 闰年
		return 29
	}

	return 28
}

// GetDayNumsByMouth 通过月份获取天数
func GetDayNumsByMouth(year, mouth int) int {
	switch mouth {
	case 1, 3, 5, 7, 8, 10, 12:
		return 31
	case 4, 6, 9, 11:
		return 30
	case 2:
		return getFebDayNumsByYear(year)
	default:
		return 0
	}
}
