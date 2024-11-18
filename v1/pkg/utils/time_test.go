package utils

import (
	"testing"
	"time"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestGetTimeStampFromEpochTime(t *testing.T) {
	tests := []struct {
		name        string
		epoch       int64
		expected    *timestamp.Timestamp
		expectedErr bool
	}{
		{
			name:        "Valid epoch time",
			epoch:       1619193600,
			expected:    &timestamppb.Timestamp{Seconds: 1619193600},
			expectedErr: false,
		},
		{
			name:        "Zero epoch time",
			epoch:       0,
			expected:    &timestamppb.Timestamp{Seconds: 0},
			expectedErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetTimeStampFromEpochTime(tt.epoch)
			assert.Equal(t, tt.expected, result, "Result timestamp does not match expected")
		})
	}
}

func TestGetEpochTimeFromTime(t *testing.T) {
	tests := []struct {
		name        string
		inputTime   time.Time
		expected    int64
		expectedErr bool
	}{
		{
			name:        "Valid time",
			inputTime:   time.Date(2024, time.April, 26, 12, 0, 0, 0, time.UTC),
			expected:    1714132800, // April 26, 2024, 12:00:00 PM UTC
			expectedErr: false,
		},
		{
			name:        "Zero time",
			inputTime:   time.Time{},
			expected:    -62135596800,
			expectedErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetEpochTimeFromTime(tt.inputTime)
			assert.Equal(t, tt.expected, result, "Result epoch time does not match expected")
		})
	}
}

func TestGetEpochTimeFromTimeStamp(t *testing.T) {
	tests := []struct {
		name     string
		input    *timestamp.Timestamp
		expected int64
	}{
		{
			name:     "Valid timestamp",
			input:    timestamppb.New(time.Date(2022, time.January, 1, 0, 0, 0, 0, time.UTC)),
			expected: 1640995200, // January 1, 2022, 00:00:00 UTC
		},
		{
			name:     "Nil timestamp",
			input:    nil,
			expected: time.Now().Unix(), // Expecting current time since input is nil
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetEpochTimeFromTimeStamp(tt.input)
			assert.Equal(t, tt.expected, result, "Result epoch time does not match expected")
		})
	}
}

func TestGetCalenderDetail(t *testing.T) {
	tests := []struct {
		name     string
		epoch    int64
		expected *CalendarDetails
	}{
		{
			name:  "Valid epoch time",
			epoch: 1619193600, // April 23, 2021, 12:00:00 PM UTC
			expected: &CalendarDetails{
				Time:        "9:30 PM", // Time in Asia/Kolkata timezone
				Day:         "Fri",
				Date:        "23",
				Month:       "Apr",
				Year:        "2021",
				DayWithDate: "Friday, April 23",
			},
		},
		{
			name:  "Zero epoch time",
			epoch: 0, // January 1, 1970, 00:00:00 UTC
			expected: &CalendarDetails{
				Time:        "5:30 AM", // Time in Asia/Kolkata timezone
				Day:         "Thu",
				Date:        "01",
				Month:       "Jan",
				Year:        "1970",
				DayWithDate: "Thursday, January 1",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetCalenderDetail(tt.epoch)
			assert.Equal(t, tt.expected, result, "Result calendar details do not match expected")
		})
	}
}
