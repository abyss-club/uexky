package uexky

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

type ModelA struct {
	VI int    `json:"vi"`
	VS string `json:"vs"`
	iv bool
}

type ModelB struct {
	ID string  `json:"id"`
	A  *ModelA `json:"a"`
}

func TestCache(t *testing.T) {
	pool := NewRedisPool()
	conn := pool.Get()

	u := &Uexky{Redis: conn}
	m := &ModelB{"test", &ModelA{12, "test", true}}
	want := &ModelB{"test", &ModelA{12, "test", false}}

	if err := SetCache(u, "TestCache", m, 3600); err != nil {
		t.Errorf("SetCache error: %v", err)
	}

	got := &ModelB{}
	if exist, err := GetCache(u, "TestCacheX", got); err != nil {
		t.Errorf("GetCache(X) error: %v", err)
	} else if exist {
		t.Errorf("GetCache(X) should be empty")
	}

	if exist, err := GetCache(u, "TestCache", got); err != nil {
		t.Errorf("GetCache() error: %v", err)
	} else if !exist {
		t.Errorf("GetCache() should not be empty")
	}

	if !cmp.Equal(got, want, cmp.AllowUnexported(ModelA{})) {
		t.Errorf("GetCache() = %v, want = %v", got, want)
	}
}
