package uid

import (
	"bytes"
	"testing"
	"time"
)

func Test_base64charsToInt64(t *testing.T) {
	for i, b := range base64chars {
		got, err := base64charsToInt64(byte(b))
		if err != nil {
			t.Errorf("base64charsToInt64() error = %v", err)
		}
		if got != int64(i) {
			t.Errorf("base64charsToInt64(%s) = %v, want %v", string(b), got, i)
		}
	}
	_, err := base64charsToInt64(byte('^'))
	if err == nil {
		t.Errorf("base64charsToInt64(^), err = nil, want error")
	}
}

func TestParseUID(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    UID
		wantErr bool
	}{
		{
			name:    "normal",
			args:    args{"AABA"}, // AABA -> BAAA -> (1<<18)
			want:    UID(1 << 18),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseUID(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseUID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ParseUID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUID_ToBase64String(t *testing.T) {
	tests := []struct {
		name string
		u    UID
		want string
	}{
		{
			name: "normal",
			u:    UID(1 << 18),
			want: "AABA", // (1<<18) -> BAAA -> AABA
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.u.ToBase64String(); got != tt.want {
				t.Errorf("UID.ToBase64String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerator_NewUID(t *testing.T) {
	if zs := timeZero.Format("2006-01-02T15:04:05.000"); zs != "2018-03-01T00:00:00.000" {
		t.Errorf("time zero is wrong")
	}
	g := Generator{count: 1}
	prev := UID(0)
	getMilliseconds := func(t time.Time) int64 {
		return t.UnixNano() / int64(time.Millisecond)
	}
	for i := int64(1); i <= 10; i++ {
		now := time.Now()
		count := i
		uid := g.NewUID()
		if getMilliseconds(uid.GetTime()) < getMilliseconds(now) {
			t.Errorf("uid's timestamp part is error, should larger than %v, got %v", now, uid.GetTime())
		}
		if int64((uid>>9)%512) != count {
			t.Errorf("uid's count part is error, want=%v, got=%v", count, (uid>>9)%512)
		}
		if uid <= prev {
			t.Error("new uid must larger chan preview")
		}
		prev = uid
	}
}

func TestUID_UnmarshalGQL(t *testing.T) {
	uid := NewUID()
	str := uid.ToBase64String()
	type args struct {
		v interface{}
	}
	tuid := UID(0)
	tests := []struct {
		name    string
		u       *UID
		args    args
		wantErr bool
	}{
		{
			name:    "unmarshal",
			u:       &tuid,
			args:    args{str},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.u.UnmarshalGQL(tt.args.v); (err != nil) != tt.wantErr {
				t.Errorf("UID.UnmarshalGQL() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUID_MarshalGQL(t *testing.T) {
	uid := NewUID()
	tests := []struct {
		name  string
		u     UID
		wantW string
	}{
		{
			name:  "marshal",
			u:     uid,
			wantW: uid.ToBase64String(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			tt.u.MarshalGQL(w)
			if gotW := w.String(); gotW != tt.wantW {
				t.Errorf("UID.MarshalGQL() = %v, want %v", gotW, tt.wantW)
			}
		})
	}
}
