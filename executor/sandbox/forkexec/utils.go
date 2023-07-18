package forkexec

import (
	"regexp"
	"runtime"
	"strconv"
)

func goVersionGEQ1dot17() bool {
	re, err := regexp.Compile("go(\\d+).(\\d+)(.*)")
	if err != nil {
		return false
	}
	matched := re.FindAllStringSubmatch(runtime.Version(), -1)
	if len(matched) < 1 || len(matched[0]) < 3 {
		return false
	}
	minor, err := strconv.ParseInt(matched[0][2], 10, 32)
	if err != nil {
		return false
	}
	return minor >= 17
}
