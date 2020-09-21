package services

import (
	"dt-services/errors"
	"fmt"
	"strconv"
	"strings"
	"time"
)

func getDate(timestamp time.Time) string {
	timestampS := strings.Split(timestamp.String(), "T")
	timestampSS := strings.Split(timestampS[0], " ")
	return timestampSS[0]
}

func getPeriod(date, period string) string {
	if period == PERIOD_DAILY {
		return date
	} else if period == PERIOD_WEEKLY {
		timeStrS := strings.Split(date, "-")
		year, err := strconv.Atoi(timeStrS[0])
		errors.Catch(err)
		month, err := strconv.Atoi(timeStrS[1])
		errors.Catch(err)
		day, err := strconv.Atoi(timeStrS[2])
		errors.Catch(err)

		var week int
		date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
		isoYear, isoWeek := date.ISOWeek()
		for date.Weekday() != time.Monday { // iterate back to Monday
			date = date.AddDate(0, 0, -1)
			isoYear, isoWeek = date.ISOWeek()
		}
		for isoYear < year { // iterate forward to the first day of the first week
			date = date.AddDate(0, 0, 1)
			isoYear, isoWeek = date.ISOWeek()
		}
		for isoWeek < week { // iterate forward to the first day of the given week
			date = date.AddDate(0, 0, 1)
			isoYear, isoWeek = date.ISOWeek()
		}
		data := strings.Split(date.String(), "-")
		return fmt.Sprintf("%s-%s-%s", data[0], data[1], data[2][:2])
	} else if period == PERIOD_MONTHLY {
		timeStrS := strings.Split(date, "-")
		return fmt.Sprintf("%s-%s", timeStrS[0], timeStrS[1])
	}

	return ""
}

func parseStrToTime(xtime string) (*time.Time, error) {
	timestampS := strings.Split(xtime, "T")
	timestampYMD := strings.Split(timestampS[0], "-")
	timestampHM := strings.Split(timestampS[1], ":")

	year, err := strconv.Atoi(timestampYMD[0])
	if err != nil {
		return nil, err
	}
	monthI, err := strconv.Atoi(timestampYMD[1])
	if err != nil {
		return nil, err
	}
	day, err := strconv.Atoi(timestampYMD[2])
	if err != nil {
		return nil, err
	}

	hour, err := strconv.Atoi(timestampHM[0])
	if err != nil {
		return nil, err
	}

	min, err := strconv.Atoi(timestampHM[1])
	if err != nil {
		return nil, err
	}

	loc, err := time.LoadLocation("UTC")
	if err != nil {
		return nil, err
	}
	timestamp := time.Date(year, time.Month(monthI), day, hour, min, 0, 0, loc)
	return &timestamp, nil
}
