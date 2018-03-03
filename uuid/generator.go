package uuin

import (
	"crypto/rand"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

// Generator generate UUID by specified sections.
type Generator struct {
	Sections []Section
}

// New UUID
func (gen *Generator) New() (string, error) {
	var sections []string
	for _, section := range gen.Sections {
		sectionStr, err := section.New()
		if err != nil {
			return "", err
		}
		sections = append(sections, sectionStr)
	}
	return strings.Join(sections, ""), nil
}

// Section can only generate a single type ID
type Section interface {
	New() (string, error)
}

// TimestampSection generate id by timestamp, unit is 'ms'
type TimestampSection struct {
	Length uint
}

// 2018-1-1 00:00:00 UTC (micro seconds)
const timeZero = 1514764800000

// New timestamp section
func (ts *TimestampSection) New() (string, error) {
	ns := time.Now().UnixNano()
	ms := int(ns/1000/1000) - timeZero

	bs := make([]byte, ts.Length)
	for i := uint(0); i != ts.Length; i++ {
		// bits: 111111 => 0x3f
		bs[ts.Length-1-i] = base64EcodeByte(byte(ms & 0x3f))
		ms = ms >> 6
	}
	if ms != 0 {
		return "", fmt.Errorf("timestamp out of range")
	}
	return string(bs), nil
}

// RandomSection generate id in random
type RandomSection struct {
	Length uint
}

// New random section
func (rs *RandomSection) New() (string, error) {
	unitLen := int((rs.Length + 3) / 4)
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
	for i := uint(0); i != rs.Length; i++ {
		base64[i] = base64EcodeByte(base64[i])
	}
	return string(base64[0:rs.Length]), nil
}

// CounterSection generate id by counter.
type CounterSection struct {
	Length uint
	Last   uint
	Start  uint
	Recent uint
	Lock   sync.Mutex
}

// New counter section
func (cs *CounterSection) New() (string, error) {
	cs.Lock.Lock()
	count := cs.getCount()
	cs.Lock.Unlock()

	bytes := make([]byte, cs.Length)
	for i := uint(0); i != cs.Length; i++ {
		bytes[cs.Length-1-i] = base64EcodeByte(byte(count & 0x3f))
		count = count >> 6
	}
	if count > 0 {
		return "", fmt.Errorf("counter out of range")
	}
	return string(bytes), nil
}

func (cs *CounterSection) getCount() uint {
	now := uint(time.Now().UnixNano() / 1000000)
	next := cs.Last + 1
	limit := uint(1) << (cs.Length * 6)
	if next == limit {
		next = 0
	}

	if (now == cs.Recent) && (next == cs.Start) {
		oneMs, err := time.ParseDuration("1ms")
		if err != nil {
			log.Fatal(err)
		}
		// sleep 1 ms
		time.Sleep(oneMs)
		now = uint(time.Now().UnixNano() / 1000000)
	}
	cs.Last = next
	cs.Recent = now
	return next
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

func timeTest1() {
	gen := Generator{Sections: []Section{
		&TimestampSection{Length: 7},
		&CounterSection{Length: 1},
	}}
	start := time.Now().UnixNano()
	for i := 0; i != 1000; i++ {
		_, err := gen.New()
		if err != nil {
			log.Fatal(err)
		}
	}
	stop := time.Now().UnixNano()
	fmt.Println("1000 post id:", float64(stop-start)/1000000, "ms")
}

func timeTest2() {
	gen := Generator{Sections: []Section{
		&RandomSection{Length: 10},
		&CounterSection{Length: 2},
		&TimestampSection{Length: 7},
		&RandomSection{Length: 7},
	}}
	start := time.Now().UnixNano()
	for i := 0; i != 1000; i++ {
		_, err := gen.New()
		if err != nil {
			log.Fatal(err)
		}
	}
	stop := time.Now().UnixNano()
	fmt.Println("1000 token:", float64(stop-start)/1000000, "ms")
}
