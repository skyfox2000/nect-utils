package jsmodule

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

// dayjs模块
var dayjsModule = map[string]interface{}{
	"format": func(date interface{}, format string) (string, error) {
		goTime, _, err := convertToGoTime(date)
		if err != nil {
			return "", err
		}

		targetFormat := customFormat(format)
		return goTime.Format(targetFormat), nil
	},
	"date": func() string {
		return time.Now().Format("2006-01-02")
	},
	"now": func() string {
		return time.Now().Format("2006-01-02 15:04:05")
	},
	"unix": func() int64 {
		return time.Now().Unix()
	},
	"unixmicro": func() int64 {
		return time.Now().UnixMicro()
	},
	"unixnano": func() int64 {
		return time.Now().UnixNano()
	},
	"add": func(datetime interface{}, duration int, period string) (string, error) {
		return calculateTime(datetime, duration, period, 1)
	},
	"subtract": func(datetime interface{}, duration int, period string) (string, error) {
		return calculateTime(datetime, duration, period, -1)
	},
	"diff": func(datetime1 interface{}, datetime2 interface{}, period string) (int, error) {
		goTime1, _, err := convertToGoTime(datetime1)
		if err != nil {
			return 0, err
		}

		goTime2, _, err := convertToGoTime(datetime2)
		if err != nil {
			return 0, err
		}

		switch period {
		case "seconds":
			return int(goTime1.Sub(goTime2).Seconds()), nil
		case "minutes":
			return int(goTime1.Sub(goTime2).Minutes()), nil
		case "hours":
			return int(goTime1.Sub(goTime2).Hours()), nil
		case "days":
			return int(goTime1.Sub(goTime2).Hours() / 24), nil
		}
		return 0, errors.New("unsupported period")
	},
}

func init() {
	JSModules["dayjs"] = dayjsModule
}

func calculateTime(datetime interface{}, duration int, period string, multiplier int) (string, error) {
	goTime, format, err := convertToGoTime(datetime)
	if err != nil {
		return "", err
	}

	var result time.Time
	switch period {
	case "seconds":
		result = goTime.Add(time.Duration(multiplier*duration) * time.Second)
	case "minutes":
		result = goTime.Add(time.Duration(multiplier*duration) * time.Minute)
	case "hours":
		result = goTime.Add(time.Duration(multiplier*duration) * time.Hour)
	case "days":
		result = goTime.AddDate(0, 0, multiplier*duration)
	case "weeks":
		result = goTime.AddDate(0, 0, multiplier*7*duration)
	case "months":
		result = goTime.AddDate(0, multiplier*duration, 0)
	case "years":
		result = goTime.AddDate(multiplier*duration, 0, 0)
	default:
		return "", errors.New("unsupported period")
	}

	switch datetime.(type) {
	case string:
		return result.Format(format), nil
	default:
		return "", errors.New("unsupported date format")
	}
}

// 将日期数据转换为Go时间格式
func convertToGoTime(date interface{}) (time.Time, string, error) {
	switch date := date.(type) {
	case string:
		// 尝试解析日期字符串
		formats := []string{"2006-01-02", "2006-01-02 15:04:05", "2006-01-02T15:04:05", "2006-01-02T15:04:05.999-07:00"}
		var goTime time.Time
		var parseError error

		var format string
		for _, f := range formats {
			goTime, parseError = time.Parse(f, date)
			if parseError == nil {
				format = f
				break
			}
		}

		if parseError != nil {
			return time.Time{}, "", errors.New("failed to parse date: " + parseError.Error())
		}

		return goTime, format, nil
	case float64:
		// 尝试解析JavaScript时间戳
		seconds := int64(date / 1000)
		nanos := int64((date - float64(seconds*1000)) * 1e6)
		return time.Unix(seconds, nanos), "Unix", nil
	default:
		return time.Time{}, "", errors.New("unsupported date format")
	}
}

func customFormat(inputFormat string) string {
	targetFormat := inputFormat
	re := regexp.MustCompile(`(YYYY|MM|DD|HH|mm|ss)`)
	matches := re.FindAllString(targetFormat, -1)

	for _, match := range matches {
		switch match {
		case "YYYY":
			targetFormat = strings.Replace(targetFormat, "YYYY", "2006", 1)
		case "MM":
			targetFormat = strings.Replace(targetFormat, "MM", "01", 1)
		case "DD":
			targetFormat = strings.Replace(targetFormat, "DD", "02", 1)
		case "HH":
			targetFormat = strings.Replace(targetFormat, "HH", "15", 1)
		case "mm":
			targetFormat = strings.Replace(targetFormat, "mm", "04", 1)
		case "ss":
			targetFormat = strings.Replace(targetFormat, "ss", "05", 1)
		}
	}

	return targetFormat
}
