package uuid64

import (
	"testing"
	"time"
)

func TestBase64Encode(t *testing.T) {
	chars := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"
	for i, c := range chars {
		if string(base64EcodeByte(byte(i))) != string(c) {
			t.Fatalf("%v must be encode to %s", i, string(c))
		}
	}
}

func TestBase64EncodeInt64(t *testing.T) {
	cases := map[int64]string{
		0:       "A",
		1:       "B",
		1 << 6:  "BA",
		1 << 8:  "EA",
		1 << 12: "BAA",
		1 << 16: "QAA",
	}
	for i, s := range cases {
		bs := base64EcodeInt64(i)
		if bs != s {
			t.Fatalf("%v must be encode to %s, but get %s", i, string(s), bs)
		}
	}
}

func TestBase64EncodeBytes(t *testing.T) {
	cases := map[string][]byte{
		"":         []byte{},
		"AAAA":     []byte{byte(0)},
		"AQAA":     []byte{byte(1)},
		"AQEA":     []byte{byte(1), byte(1)},
		"AQEB":     []byte{byte(1), byte(1), byte(1)},
		"AQEBAQAA": []byte{byte(1), byte(1), byte(1), byte(1)},
	}
	for s, bytes := range cases {
		bs := base64EcodeBytes(bytes)
		if bs != s {
			t.Fatalf("%v must be encode to %s, but get %s", bytes, string(s), bs)
		}
	}
}

func TestTimestampSection(t *testing.T) {
	{
		ts := &TimestampSection{Length: 7, Unit: 1 * time.Second}
		a, err := ts.New()
		t.Logf("gen %s, before %v", a, time.Now().Unix())
		if err != nil {
			t.Fatal(err)
		} else if len(a) != ts.Length {
			t.Fatal("invalid length")
		}

		time.Sleep(1 * time.Second)
		b, err := ts.New()
		t.Logf("gen %s, before %v", b, time.Now().Unix())
		if err != nil {
			t.Fatal(err)
		} else if CompareUUID(a, b) >= 0 {
			t.Fatalf("result not increase")
		}
	}
	{
		ts := &TimestampSection{Length: 10, Unit: 1 * time.Second, NoPadding: true}
		a, err := ts.New()
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("gen %s", a)
		if len(a) >= ts.Length {
			t.Fatalf("result must short than %v", ts.Length)
		}
	}
}

func TestRandomSection(t *testing.T) {
	rs := &RandomSection{Length: 3}
	a, err := rs.New()
	if err != nil {
		t.Fatal(err)
	} else if len(a) != rs.Length {
		t.Fatal("invalid length")
	}
	t.Logf("get slice: %s", a)
	b, err := rs.New()
	if err != nil {
		t.Fatal(err)
	} else if len(a) != rs.Length {
		t.Fatal("it can't be equal")
	}
	t.Logf("get slice: %s", b)
}

func TestCounterSection(t *testing.T) {
	{
		cs := &CounterSection{Length: 1, Unit: time.Nanosecond}
		for i := int64(0); i < 70; i++ {
			wanted := base64EcodeInt64(i % 64)
			got, err := cs.New()
			if err != nil {
				t.Fatal(err)
			}
			if wanted != got {
				t.Fatalf("want %s, but get %s", wanted, got)
			}
		}
	}
	{
		uuids := []string{}
		gen := &Generator{
			Sections: []Section{
				&TimestampSection{Length: 7, Unit: time.Second},
				&CounterSection{Length: 1, Unit: time.Second},
			},
		}
		for i := int64(0); i < 500; i++ {
			uuid, err := gen.New()
			if err != nil {
				t.Fatal(err)
			}
			uuids = append(uuids, uuid)
		}
		uuidSet := map[string]struct{}{}
		for _, uuid := range uuids {
			if _, ok := uuidSet[uuid]; ok {
				t.Fatalf("Gen same uuid: '%s'", uuid)
			}
			uuidSet[uuid] = struct{}{}
		}
	}
	{
		uuids := []string{}
		gen := &Generator{
			Sections: []Section{
				&TimestampSection{Length: 6, Unit: time.Second},
				&CounterSection{Length: 2, Unit: time.Second},
			},
		}
		for i := int64(0); i < 10000; i++ {
			uuid, err := gen.New()
			if err != nil {
				t.Fatal(err)
			}
			uuids = append(uuids, uuid)
		}
		uuidSet := map[string]struct{}{}
		for _, uuid := range uuids {
			if _, ok := uuidSet[uuid]; ok {
				t.Fatalf("Gen same uuid: '%s'", uuid)
			}
			uuidSet[uuid] = struct{}{}
		}
	}
}

func TestBenchmark(t *testing.T) {
	// time, count, random
	cases := [][]int{
		// ms
		[]int{7, 2, 2},
		[]int{7, 2, 1},
		[]int{7, 1, 2},
		[]int{7, 1, 1},
		// s
		[]int{6, 2, 2},
		[]int{6, 2, 1},
		[]int{6, 1, 2},
		[]int{6, 1, 1},
	}
	for i, setting := range cases {
		var unit time.Duration
		var count int
		if i < 4 {
			unit = time.Millisecond
		} else {
			unit = time.Second
		}
		if i < 6 {
			count = 10000
		} else {
			count = 500
		}
		gen := &Generator{
			Sections: []Section{
				&TimestampSection{Length: setting[0], Unit: unit},
				&CounterSection{Length: setting[1], Unit: unit},
				&RandomSection{Length: setting[2]},
			},
		}
		start := time.Now()
		for i := 0; i < count; i++ {
			if _, err := gen.New(); err != nil {
				t.Fatal(err)
			}
		}
		stop := time.Now()
		s := float64(stop.Sub(start)) / float64(time.Second)
		ps := float64(count) / s
		t.Logf("Gen(time[%v], count[%v], random[%v]) %v uuid, cost %v, %v uuids per second",
			setting[0], setting[1], setting[2], count, stop.Sub(start), ps)
	}
}
