package main

import "time"

// nowFunc is 現在時刻を返すfunction
var nowFunc = func() time.Time {
	return time.Now()
}

// Now is 現在時刻を返す
func Now() time.Time {
	return nowFunc()
}
