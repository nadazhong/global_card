package misc

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"
)

func CheckError(err error) {
	if err != nil {
		fmt.Printf("Fatal error: %v", err)
		os.Exit(-1)
	}
}

// 返回毫秒时间戳。
func ms(t time.Time) int32 {
	return (int32)((t.UnixNano() / 1000000) % 1000000000)
}

// 返回unix时间戳。
func Current() int64 {
	return int64(time.Now().UnixNano() / 1000000)
}

// date format: "2006-01-02 13:04:00"
func S2UnixTime(value string) int64 {
	re := regexp.MustCompile(`([\d]+)-([\d]+)-([\d]+) ([\d]+):([\d]+):([\d]+)`)
	slices := re.FindStringSubmatch(value)
	if slices == nil || len(slices) != 7 {
		fmt.Printf("time[%s] format error, expect format: 2006-01-02 13:04:00...", value)
		return 0
	}
	year, _ := strconv.Atoi(slices[1])
	month, _ := strconv.Atoi(slices[2])
	day, _ := strconv.Atoi(slices[3])
	hour, _ := strconv.Atoi(slices[4])
	min, _ := strconv.Atoi(slices[5])
	sec, _ := strconv.Atoi(slices[6])
	loc, _ := time.LoadLocation("UTC") // use UTC instend of Local
	t := time.Date(year, time.Month(month), day, hour, min, sec, 0, loc)
	return int64(t.UnixNano() / 1000000)
}

/*
  a simple 32 bit checksum that can be upadted from either end
  (inspired by Mark Adler's Adler-32 checksum)
*/
func Adler32(data []byte) uint32 {
	size := len(data)
	s1 := uint32(0)
	s2 := uint32(0)
	i := 0
	for i < size-4 {
		s2 += 4*(s1+uint32(data[i])) + 3*uint32(data[i+1]) + 2*uint32(data[i+2]) + uint32(data[i+3])
		s1 += uint32(data[i+0]) + uint32(data[i+1]) + uint32(data[i+2]) + uint32(data[i+3])
		i += 4
	}
	for i < size {
		s1 += uint32(data[i])
		s2 += s1
		i++
	}
	return (s1 & 0xffff) + (s2 << 16)
}
