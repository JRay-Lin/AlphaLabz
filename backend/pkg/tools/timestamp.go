package tools

import "time"

func Timestamp() string {
	timestamp := time.Now().Unix()
	formattedTime := time.Unix(timestamp, 0).Format("2001-01-01 23:59:59")
	return formattedTime
}
