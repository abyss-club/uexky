package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

// RangeType is charactor type
type RangeType int

// charactor type
const (
	TimeStamp RangeType = iota
	Counter
	Random
)

// RangeConfig is a config of range charator
type RangeConfig struct {
	Type RangeType
	Len  uint
}

type rangeHandler func(*RangeConfig) (string, error)

var rangeHandlers = map[RangeType]rangeHandler{
	RangeType(Random):    randBase64,
	RangeType(TimeStamp): timeStampBase64,
	RangeType(Counter):   counterBase64,
}

// GenID generat ID by configs
func GenID(configs []RangeConfig) (string, error) {
	var rst []string
	for _, config := range configs {
		handler, ok := rangeHandlers[config.Type]
		if !ok {
			return "", fmt.Errorf("unknown type: %v", config.Type)
		}
		s, err := handler(&config)
		if err != nil {
			return "", err
		}
		rst = append(rst, s)
	}
	return strings.Join(rst, ""), nil
}

// 2018-1-1 00:00:00 UTC (micro seconds)
var timeZero = 1514764800000

func timeStampBase64(r *RangeConfig) (string, error) {
	ns := time.Now().UnixNano()
	ms := int(ns/1000/1000) - timeZero

	bs := make([]byte, r.Len)
	for i := uint(0); i != r.Len; i++ {
		// bits: 111111 => 0x3f
		bs[r.Len-1-i] = base64EcodeByte(byte(ms & 0x3f))
		ms = ms >> 6
	}
	if ms != 0 {
		return "", fmt.Errorf("timestamp out of range")
	}
	return string(bs), nil
}

func randBase64(r *RangeConfig) (string, error) {
	unitLen := int((r.Len + 3) / 4)
	byteLen := unitLen * 3
	raw := make([]byte, byteLen)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}

	base64Len := unitLen * 4
	base64 := make([]byte, base64Len)
	for i := 0; i != unitLen; i++ {
		base64[i*4] = (raw[i*3] & 0xfc) >> 2
		base64[i*4+1] = ((raw[i*3] & 0x03) << 4) | ((raw[i*3+1] & 0xf0) >> 4)
		base64[i*4+2] = ((raw[i*3+1] & 0x0f) << 2) | ((raw[i*3+2] & 0xc0) >> 6)
		base64[i*4+3] = raw[i*3+2] & 0x3f
	}
	for i := uint(0); i != r.Len; i++ {
		base64[i] = base64EcodeByte(base64[i])
	}
	return string(base64[0:r.Len]), nil
}

type counter struct {
	Max         uint
	Last        uint
	Start       uint
	MicroSecond uint
	Lock        sync.Mutex
}

func (c *counter) GetConut() uint {
	c.Lock.Lock()
	nMs := uint(time.Now().UnixNano() / 1000000)
	next := c.Last + 1
	if next == c.Max {
		next = 0
	}
	if (nMs == c.MicroSecond) && (next == c.Start) {
		oneMs, err := time.ParseDuration("1ms")
		if err != nil {
			log.Fatal(err)
		}
		time.Sleep(oneMs)
	}
	c.Last = next
	c.MicroSecond = nMs
	fmt.Println("ms:", nMs, " c:", next, "counter:", c)
	c.Lock.Unlock()
	return next
}

var counters = map[uint]*counter{}
var countersLock sync.Mutex

func counterBase64(r *RangeConfig) (string, error) {
	max := uint(1) << (r.Len * 6)
	c, ok := counters[r.Len]
	if !ok {
		countersLock.Lock()
		c = &counter{Max: max}
		counters[r.Len] = c
		countersLock.Unlock()
	}

	if r.Len == 0 {
		c.GetConut()
		return "", nil
	}

	count := c.GetConut()
	bs := make([]byte, r.Len)
	for i := uint(0); i != r.Len; i++ {
		bs[r.Len-1-i] = base64EcodeByte(byte(count & 0x3f))
		count = count >> 6
	}
	if count > 0 {
		return "", fmt.Errorf("counter out of range")
	}

	return string(bs), nil
}

func base64EcodeByte(i byte) byte {
	in := int(i)
	switch {
	case in < 26:
		return byte(in + 65)
	case in >= 26 && in < 52:
		return byte(in + 71)
	case in >= 52 && in < 62:
		return byte(in - 4)
	case in == 62:
		return byte(45)
	default:
		return byte(95)
	}
}

func main() {
	tokenConfig := []RangeConfig{
		RangeConfig{Random, 10},
		RangeConfig{TimeStamp, 7},
		RangeConfig{Random, 7},
	}

	start := time.Now().UnixNano()
	for i := 0; i != 100; i++ {
		token, err := GenID(tokenConfig)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(token)
	}
	stop := time.Now().UnixNano()
	fmt.Println(float64(stop-start) / 1000000)

	start = time.Now().UnixNano()
	for i := 0; i != 2000; i++ {
		GenID([]RangeConfig{RangeConfig{Counter, 0}})
	}
	stop = time.Now().UnixNano()
	fmt.Println(float64(stop-start) / 1000000)
}
