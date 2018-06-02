package uuid64

import (
	"crypto/rand"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
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
	Length    int
	Unit      time.Duration
	NoPadding bool
}

// 2018-1-1 00:00:00 UTC (nano seconds)
const timeZero int64 = 1514764800 * 1000 * 1000 * 1000

// New timestamp section
func (ts *TimestampSection) New() (string, error) {
	tick := (time.Now().UnixNano() - timeZero) / int64(ts.Unit)
	bs := base64EcodeInt64(tick)
	return ensureLength(bs, ts.Length, !ts.NoPadding)
}

// RandomSection generate id in random
type RandomSection struct {
	Length int
}

// New random section
func (rs *RandomSection) New() (string, error) {
	bytesLen := ((rs.Length*6)-1)/8 + 1
	raw := make([]byte, bytesLen)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}
	bs := base64EcodeBytes(raw)
	return bs[:rs.Length], nil
}

// CounterSection generate id by counter.
type CounterSection struct {
	Length int
	Unit   time.Duration

	next   int64 // number to return
	recent int64 // last time returned number
	count  int64 // count in uint
	lock   sync.Mutex
}

// New counter section
func (cs *CounterSection) New() (string, error) {
	cs.lock.Lock()
	count := cs.getCount()
	cs.lock.Unlock()
	bs := base64EcodeInt64(count)
	return ensureLength(bs, cs.Length, true)
}

func (cs *CounterSection) getCount() int64 {
	now := time.Now().UnixNano() / int64(cs.Unit)
	limit := int64(1) << uint(cs.Length*6)
	n := cs.next
	cs.next++
	if cs.next >= limit {
		cs.next = 0
	}

	if now != cs.recent {
		cs.count = 1
		cs.recent = now
		return n
	}

	cs.count++ // ensure next number is safe in any section order
	if cs.count < limit {
		return n
	}
	time.Sleep(time.Duration(
		((now+1)*int64(cs.Unit) - time.Now().UnixNano()),
	))
	cs.recent = time.Now().UnixNano() / int64(cs.Unit)
	cs.count = 1
	return n
}

const chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"

// CompareUUID lh > rh return >0, lh < rh return <0, lh == rh return 0
func CompareUUID(lh, rh string) int {
	if len(lh) != len(rh) {
		return len(lh) - len(rh)
	}
	for i := 0; i < len(lh); i++ {
		li := strings.Index(chars, string(lh[i]))
		ri := strings.Index(chars, string(rh[i]))
		if li != ri {
			return li - ri
		}
	}
	return 0
}

func base64EcodeByte(b byte) byte {
	in := int(b)
	switch {
	case in < 26:
		return byte(in + 65) // 'A-Z'
	case in >= 26 && in < 52:
		return byte(in + 71) // 'a-z'
	case in >= 52 && in < 62:
		return byte(in - 4)
	case in == 62:
		return byte(45) // '-'
	default:
		return byte(95) // '_'
	}
}

func base64EcodeInt64(i int64) string {
	bytes := []byte{base64EcodeByte(byte(i & 0x3f))}
	for i = i >> 6; i != 0; i = i >> 6 {
		bytes = append(bytes, base64EcodeByte(byte(i&0x3f)))
	}
	for i := 0; i < len(bytes)/2; i++ {
		j := len(bytes) - i - 1
		bytes[i], bytes[j] = bytes[j], bytes[i]
	}
	return string(bytes)
}

func paddingToLength(s string, l int) string {
	if len(s) >= l {
		return s
	}
	bs := []string{}
	for i := 0; i < (l - len(s)); i++ {
		bs = append(bs, "A")
	}
	bs = append(bs, s)
	return strings.Join(bs, "")
}

func ensureLength(s string, l int, padding bool) (string, error) {
	switch {
	case len(s) < l:
		if padding {
			return paddingToLength(s, l), nil
		}
		return s, nil
	case len(s) == l:
		return s, nil
	default:
		return "", errors.New("Out of range")
	}
}

func base64EcodeBytes(bs []byte) string {
	if len(bs) == 0 {
		return ""
	}

	unitLen := (len(bs)-1)/3 + 1 // 3 bytes to 1 unit
	for i := len(bs); i < unitLen*3; i++ {
		bs = append(bs, byte(0))
	}

	s := []string{}
	for i := 0; i < unitLen; i++ {
		n := int64(bs[i*3])<<16 + int64(bs[i*3+1])<<8 + int64(bs[i*3+2])
		s = append(s, paddingToLength(base64EcodeInt64(n), 4))
	}
	return strings.Join(s, "")
}
