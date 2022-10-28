package pkg

import (
	"errors"
	"fmt"
	"testing"
)

func TestTimeFormate(t *testing.T) {
	result, err := TimeFormat("10:00")
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println(result)
	}
}

func TestTimeDurationFormat(t *testing.T) {
	result, err := TimeDurationFormat("20m")
	if err != nil {
		t.Error(err)
	} else {
		fmt.Println(result)
	}
}

func TestTimeDurationCompare(t *testing.T) {

	time1, err := TimeDurationFormat("-60m")
	if err != nil {
		t.Error(err)
	}

	time2, err := TimeDurationFormat("0m")
	if err != nil {
		t.Error(err)
	}
	if !TimeDurationCompare(time1, time2) {
		t.Error(errors.New("compare err"))
	} else {
		fmt.Println()
	}
}

func TestTimeCompare(t *testing.T) {
	time1, err := TimeFormat("10:20")
	if err != nil {
		t.Error(err)
	}
	time2, err := TimeFormat("10:20")
	if err != nil {
		t.Error(err)
	}
	if TimeCompare(time1, time2) == true {
		fmt.Println("true")
	} else {
		t.Error(errors.New("compare err"))
	}
}

func TestTimeAdd(t *testing.T) {
	time1, err := TimeFormat("10:00")
	if err != nil {
		t.Error(err)
	}
	time2, err := TimeFormat("10:20")
	if err != nil {
		t.Error(err)
	}
	duration := TimeSub(time2, time1)

	result := TimeAdd(time1, duration)
	fmt.Println(result)
}

func TestDurationAdd(t *testing.T) {
	time1, err := TimeDurationFormat("40m")
	if err != nil {
		t.Error(err)
	}
	time2, err := TimeDurationFormat("40m")
	if err != nil {
		t.Error(err)
	}
	duration := DurationAdd(time2, time1)
	fmt.Println(duration)
}

func TestTimeSub(t *testing.T) {
	time1, err := TimeFormat("10:00")
	if err != nil {
		t.Error(err)
	}
	time2, err := TimeFormat("10:20")
	if err != nil {
		t.Error(err)
	}
	result := TimeSub(time2, time1)
	fmt.Println(result)

}

func TestAdd(t *testing.T) {
	time1, err := TimeFormat("10:00")
	if err != nil {
		t.Error(err)
	}
	result := Add(time1, "120m")
	fmt.Println(result)
}

func TestSub(t *testing.T) {
	time1, err := TimeDurationFormat("40m")
	if err != nil {
		t.Error(err)
	}
	result := Sub(time1, "40m")
	fmt.Println(result)
}
func TestRemainTimeLessThanZero(t *testing.T) {
	time1, err := TimeDurationFormat("-20m")
	if err != nil {
		t.Error(err)
	}
	// 小于等于0
	if RemainTimeLessThanZero(time1) == true {
	} else {
		t.Error("40m less than 0m")
	}
}

func TestTimeDurationAdd(t *testing.T) {
	d, err := TimeDurationFormat("0m")
	if err != nil {
		t.Error(err)
	}
	result := TimeDurationAdd(d, "20m")
	fmt.Println(result)

}
