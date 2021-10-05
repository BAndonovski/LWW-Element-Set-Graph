package main

import (
	"flag"
	"time"
)

//used for tests
var mt time.Time

func MockTime(t time.Time) {
	mt = t
}

func GetMockTime() time.Time {
	return mt
}

func Now() time.Time {
	if flag.Lookup("test.v") == nil {
		return time.Now()
	} else {
		return mt
	}
}
