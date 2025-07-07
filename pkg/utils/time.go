package utils

import "time"

// FormatDate .
func FormatDate(tm time.Time) string {
	return tm.Format("2006-01-02")
}

// FormatDateMonth .
func FormatDateMonth(tm time.Time) string {
	return tm.Format("2006-01")
}

// FormatTime .
func FormatTime(tm time.Time) string {
	return tm.Format("2006-01-02 15:04:05")
}

// FormatLocalTime .
func FormatLocalTime(tm time.Time) string {
	loc, err := time.LoadLocation("Local")
	if err != nil {
		return ""
	}
	return tm.In(loc).Format("2006-01-02 15:04:05")
}

// TimestampToTime .
func TimestampToTime(timestamp int64) string {
	return time.Unix(timestamp, 0).Format("2006-01-02 15:04:05")
}

// FormatStr2Time .
func FormatStr2Time(strTime string) (time.Time, error) {
	formatTime, err := time.Parse("2006-01-02 15:04:05", strTime)
	if err != nil {
		return time.Now(), err
	}
	return formatTime.Add(-8 * time.Hour), nil
}

// FormatStr2Timestamps .
func FormatStr2Timestamps(strTime string) int64 {
	formatTime, err := FormatStr2Time(strTime)
	if err != nil {
		return time.Now().Unix()
	}
	return formatTime.Unix()
}

// FormatStr2LocalTime .
func FormatStr2LocalTime(strTime string) (time.Time, error) {
	formatTime, err := time.Parse("2006-01-02 15:04:05", strTime)
	if err != nil {
		return time.Now(), err
	}
	return formatTime, nil
}

// FormatNowTime .
func FormatNowTime() string {
	return time.Now().Format("20060102150405")
}
