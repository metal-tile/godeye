package main

import (
	"testing"
	"time"
)

func TestNow(t *testing.T) {
	count := -1
	now1 := time.Date(2017, 12, 30, 10, 0, 0, 0, time.UTC)
	now2 := time.Date(2017, 12, 31, 10, 0, 0, 0, time.UTC)
	nowArray := []time.Time{
		now1, now2,
	}
	nowFunc = func() time.Time {
		count++
		return nowArray[count]
	}

	time1 := Now()
	if time1.Equal(now1) == false {
		t.Fatalf("unexpected now 1")
	}
	time2 := Now()
	if time2.Equal(now2) == false {
		t.Fatalf("unexpected now 2")
	}
}
