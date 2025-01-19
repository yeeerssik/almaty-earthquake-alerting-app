package helper

import "time"

// TimestampToDate timeConvert from unix epoch to timestamp
func TimestampToDate(timestamp int64) time.Time {
	location, err := time.LoadLocation("Asia/Almaty")
	if err != nil {
		panic(err)
	}
	return time.Unix(timestamp/1e3, 0).In(location)
}
