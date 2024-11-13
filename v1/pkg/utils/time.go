package utils

import (
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CalendarDetails struct {
	Time        string `json:"time,omitempty"`
	Day         string `json:"day,omitempty"`
	Date        string `json:"date,omitempty"`
	Month       string `json:"month,omitempty"`
	Year        string `json:"year,omitempty"`
	DayWithDate string `json:"day_with_date,omitempty"`
}

func GetCalenderDetail(epoch int64) *CalendarDetails {
	pbTime := GetTimeStampFromEpochTime(epoch)
	// Extract day, date, and month from the startTime
	// We need to set Zone here
	// GMT time zoe in case of pbst
	var locationAsiaKolkata, err = time.LoadLocation("Asia/Kolkata")
	if err != nil {
		return nil
	}

	startTime := pbTime.AsTime().In(locationAsiaKolkata)
	day := startTime.Format("Mon")
	date := startTime.Format("02")
	month := startTime.Format("Jan")
	year := startTime.Format("2006")
	dayWithDate := startTime.Format("Monday, January 2")

	return &CalendarDetails{
		Time:        startTime.Format("3:04 PM"),
		Day:         day,
		Date:        date,
		Month:       month,
		Year:        year,
		DayWithDate: dayWithDate,
	}
}

func GetEpochTimeFromTimeStamp(t *timestamp.Timestamp) int64 {
	goTime := GetTimeFromTimeStamp(t)
	return GetEpochTimeFromTime(goTime)
}

func GetTimeFromTimeStamp(t *timestamp.Timestamp) time.Time {
	if t == nil {
		return time.Now()
	}

	return t.AsTime()
}

func GetEpochTimeFromTime(t time.Time) int64 {
	return t.Unix()
}

func GetTimeStampFromEpochTime(epoch int64) *timestamp.Timestamp {
	return timestamppb.New(time.Unix(epoch, 0))
}
