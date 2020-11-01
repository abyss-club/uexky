package uid

import (
	"fmt"
	"io"
	"math/rand"
	"time"

	"gitlab.com/abyss.club/uexky/lib/errors"
)

type UID int64

const base64chars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"
const minStrLen = 4
const minIntVal = 1 << 18

func base64charsToInt64(b byte) (int64, error) {
	// ascii           base64
	// A: 65, Z: 90    A: 0, Z: 25
	// a: 97, z: 122   a: 26, z: 51
	// 0: 48, 9:57     0: 52, z: 61
	// -: 45, _: 95    -: 62, _: 63
	switch {
	case b >= 65 && b <= 90:
		return int64(b - (65 - 0)), nil
	case b >= 97 && b <= 122:
		return int64(b - (97 - 26)), nil
	case b >= 48 && b <= 57:
		return int64(b + (52 - 48)), nil
	case b == 45:
		return int64(62), nil
	case b == 95:
		return int64(63), nil
	default:
		return 0, errors.BadParams.Errorf("invalid char: %v", b)
	}
}

func ParseUID(s string) (UID, error) {
	chars := []byte(s)
	length := len(chars)
	if len(chars) < minStrLen {
		return 0, errors.BadParams.Errorf("invalid uid base64 string: %s", s)
	}
	chars[0], chars[1] = chars[1], chars[0]
	for i := 2; i < length/2+1; i++ {
		chars[i], chars[length-i+1] = chars[length-i+1], chars[i]
	}
	uid := UID(0)
	// ascii           base64
	// A: 65, Z: 90    A: 0, Z: 25
	// a: 97, z: 122   a: 26, z: 51
	// 0: 48, 9:57     0: 52, z: 61
	// -: 45, _: 95    -: 62, _: 63
	for i := length - 1; i >= 0; i-- {
		i, err := base64charsToInt64(chars[i])
		if err != nil {
			return UID(0), errors.Wrapf(err, "ParseUID(s=%s)", s)
		}
		uid = uid*64 + UID(i)
	}
	return uid, nil
}

func (u UID) ToBase64String() string {
	if u < minIntVal {
		panic(fmt.Sprintf("invalid uid: %+v", u))
	}
	var chars []byte
	for i := u; i > 0; i /= 64 {
		c := i % 64
		chars = append(chars, base64chars[c])
	}
	length := len(chars)
	chars[0], chars[1] = chars[1], chars[0]
	for i := 2; i < length/2+1; i++ {
		chars[i], chars[length-i+1] = chars[length-i+1], chars[i]
	}
	return string(chars)
}

func (u UID) GetTime() time.Time {
	timestamp := int64(u >> 18)
	return timeZero.Add(time.Duration(timestamp) * time.Millisecond)
}

type Generator struct {
	count int64
}

var timeZero, _ = time.Parse("2006-01-02", "2018-03-01")

func (g *Generator) NewUID() UID {
	timestamp := (time.Since(timeZero) / time.Millisecond) << 18
	count := g.count << 9
	randnum := rand.Int63n(512)
	g.count = (g.count + 1) % 512
	return UID(int64(timestamp) + count + randnum)
}

var gGenerator = Generator{count: 1}

func NewUID() UID {
	return gGenerator.NewUID()
}

// UnmarshalGQL for UID scalar type in graphql
func (u *UID) UnmarshalGQL(v interface{}) error {
	uidStr, ok := v.(string)
	if !ok {
		return errors.BadParams.New("uid in graphql must be strings")
	}
	var err error
	*u, err = ParseUID(uidStr)
	if err != nil {
		return errors.Wrapf(err, "UnmarshalGQL(v=%+v)", v)
	}
	return nil
}

// MarshalGQL for UID scalar type in graphql
func (u UID) MarshalGQL(w io.Writer) {
	rawStr := fmt.Sprintf(`"%s"`, u.ToBase64String())
	_, err := w.Write([]byte(rawStr))
	if err != nil {
		panic(err)
	}
}

func RandomBase64Str(length int) string {
	bytes := make([]byte, 0, length)
	for i := 0; i < length; i++ {
		i := rand.Intn(64)
		bytes = append(bytes, base64chars[i])
	}
	return string(bytes)
}
