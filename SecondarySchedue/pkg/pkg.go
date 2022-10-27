package pkg

import (
	"fmt"
	"time"
	//"github.com/golangci/golangci-lint/pkg/result"
)

//只能 "20m" 格式
func TimeDurationFormat(str string) (time.Duration, error) {
	result, err := time.ParseDuration(str)
	if err != nil {
		return 0, err
	} else {
		return result, nil
	}
}

// time1 < time2 true
func TimeDurationCompare(time1 time.Duration, time2 time.Duration) bool {
	if int(time1) <= int(time2) {
		return true
	} else {
		//fmt.Println(int(time2))
		return false
	}

}

func TimeFormat(str string) (time.Time, error) {
	str = "2016-07-25 " + str + ":00"
	layout := "2006-01-02 15:04:05"
	t, err := time.Parse(layout, str)
	t.Format("2006.01.02 15:04:05")
	if err != nil {
		return time.Time{}, err
	} else {
		return t, nil
	}
}

// time1 < time2 true
func TimeCompare(time1 time.Time, time2 time.Time) bool {
	if time1.Before(time2) || time1.Equal(time2) {
		return true
	} else {
		return false
	}

}

func TimeAdd(t1 time.Time, t2 time.Duration) time.Time {
	result := t1.Add(t2)
	return result
}

func DurationAdd(t1 time.Duration, t2 time.Duration) time.Duration {
	now := time.Now()
	time1 := TimeAdd(now, t1)
	time2 := TimeAdd(time1, t2)
	result := TimeSub(time2, now)
	return result
}

func TimeSub(t1 time.Time, t2 time.Time) time.Duration {
	result := t1.Sub(t2)
	return result
}

func Add(t time.Time, durartion string) time.Time {
	d, err := TimeDurationFormat(durartion)
	if err != nil {
		fmt.Println(err)
		return time.Time{}
	} else {
		result := TimeAdd(t, d)
		return result
	}
}

// 剩余的执行时间
func Sub(t time.Duration, durartion string) time.Duration {
	now := time.Now()
	time1 := TimeAdd(now, t)
	d, err := TimeDurationFormat(durartion)
	if err != nil {
		fmt.Println(err)
	}
	time2 := TimeAdd(now, d)
	result := time1.Sub(time2)
	return result
}

func RemainTimeLessThanZero(t time.Duration) bool {
	time2, err := TimeDurationFormat("0m")
	if err != nil {
		fmt.Println(err)
	}
	if TimeDurationCompare(t, time2) {
		return true
	} else {
		return false
	}
}
