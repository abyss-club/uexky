package uexky

/*
import (
	"context"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"gitlab.com/abyss.club/uexky/mgmt"
)

func diffFlow(l, r *Flow) string {
	return cmp.Diff(l, r, cmp.AllowUnexported(Flow{}, limiter{}))
}

func TestNewFlow(t *testing.T) {
	type args struct {
		ip    string
		email string
	}
	cfg := mgmt.Config.RateLimit
	tests := []struct {
		name string
		args args
		want *Flow
	}{
		{
			name: "not logged in",
			args: args{
				ip:    "192.168.1.1",
				email: "",
			},
			want: &Flow{
				ip:    "192.168.1.1",
				email: "",
				limiters: []*limiter{
					&limiter{
						key:       "flow-ip-192.168.1.1",
						limit:     cfg.QueryLimit,
						ratio:     10,
						expire:    cfg.QueryResetTime,
						count:     0,
						remaining: cfg.QueryLimit,
					},
					&limiter{
						key:       "flow-ip-m-192.168.1.1",
						limit:     cfg.MutLimit,
						ratio:     1,
						expire:    cfg.MutResetTime,
						count:     0,
						remaining: cfg.MutLimit,
					},
				},
				queryIndex: []int{0},
				mutIndex:   []int{1},
			},
		},
		{
			name: "logged in",
			args: args{
				ip:    "192.168.1.1",
				email: "test@uexky.com",
			},
			want: &Flow{
				ip:    "192.168.1.1",
				email: "test@uexky.com",
				limiters: []*limiter{
					&limiter{
						key:       "flow-ip-192.168.1.1",
						limit:     cfg.QueryLimit,
						ratio:     10,
						expire:    cfg.QueryResetTime,
						count:     0,
						remaining: cfg.QueryLimit,
					},
					&limiter{
						key:       "flow-ip-m-192.168.1.1",
						limit:     cfg.MutLimit,
						ratio:     1,
						expire:    cfg.MutResetTime,
						count:     0,
						remaining: cfg.MutLimit,
					},
					&limiter{
						key:       "flow-email-test@uexky.com",
						limit:     cfg.QueryLimit,
						ratio:     10,
						expire:    cfg.QueryResetTime,
						count:     0,
						remaining: cfg.QueryLimit,
					},
					&limiter{
						key:       "flow-email-m-test@uexky.com",
						limit:     cfg.MutLimit,
						ratio:     1,
						expire:    cfg.MutResetTime,
						count:     0,
						remaining: cfg.MutLimit,
					},
				},
				queryIndex: []int{0, 2},
				mutIndex:   []int{1, 3},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newFlow(tt.args.ip, tt.args.email)
			if diff := diffFlow(got, tt.want); diff != "" {
				t.Errorf("NewFlow() want %v, diff = %s", tt.want, diff)
			}
		})
	}
}

func TestFlow_CostQuery(t *testing.T) {
	mgmt.LoadConfig("")
	mgmt.ReplaceConfigByEnv()
	conn := RedisPool.Get()
	errMsg := "rate limit exceeded"
	ctx := context.WithValue(context.Background(), ContextKeyRedis, conn)
	tests := []struct {
		name string
		flow   *Flow
		want *Flow
	}{
		{
			name: "not login",
			flow:   newFlow("192.168.1.1", ""),
			want: newFlow("192.168.1.1", ""),
		},
		{
			name: "logged in",
			flow:   newFlow("192.168.1.1", "test@uexky.com"),
			want: newFlow("192.168.1.1", "test@uexky.com"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn.Do("FLUSHDB")
			t.Log("cost small query")
			tt.flow.costQuery(ctx, 5)
			for _, idx := range tt.want.queryIndex {
				tt.want.limiters[idx].count = 5
			}
			if diff := diffFlow(tt.want, tt.flow); diff != "" {
				t.Fatalf("tt.want %v, diff = %v", tt.want, diff)
			}

			t.Log("cost big query")
			tt.flow.costQuery(ctx, 15)
			for _, idx := range tt.want.queryIndex {
				tt.want.limiters[idx].count = 0
				tt.want.limiters[idx].remaining -= 2
			}
			if diff := diffFlow(tt.want, tt.flow); diff != "" {
				t.Fatalf("tt.want %v, diff = %v", tt.want, diff)
			}

			t.Log("let rate limit exceeded")
			if err := tt.flow.costQuery(ctx, 3000); err == nil {
				t.Fatalf("must return err")
			} else if err.Error() != errMsg {
				t.Fatalf("err must be '%s', but get '%s'", errMsg, err.Error())
			}
			for _, idx := range tt.want.queryIndex {
				tt.want.limiters[idx].count = 0
				tt.want.limiters[idx].remaining -= 300
			}
			if diff := diffFlow(tt.want, tt.flow); diff != "" {
				t.Fatalf("tt.want %v, diff = %v", tt.want, diff)
			}
		})
	}
}

func TestFlow_CostMut(t *testing.T) {
	mgmt.LoadConfig("")
	mgmt.ReplaceConfigByEnv()
	conn := RedisPool.Get()
	errMsg := "rate limit exceeded"
	ctx := context.WithValue(context.Background(), ContextKeyRedis, conn)
	tests := []struct {
		name string
		flow   *Flow
		want *Flow
	}{
		{
			name: "not login",
			flow:   newFlow("192.168.1.1", ""),
			want: newFlow("192.168.1.1", ""),
		},
		{
			name: "logged in",
			flow:   newFlow("192.168.1.1", "test@uexky.com"),
			want: newFlow("192.168.1.1", "test@uexky.com"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conn.Do("FLUSHDB")
			t.Log("cost normal")
			tt.flow.costMut(ctx, 5)
			for _, idx := range tt.want.mutIndex {
				tt.want.limiters[idx].remaining -= 5
			}
			if diff := diffFlow(tt.want, tt.flow); diff != "" {
				t.Fatalf("tt.want %v, diff = %v", tt.want, diff)
			}

			t.Log("let rate limit exceeded")
			if err := tt.flow.costMut(ctx, 100); err == nil {
				t.Fatalf("must return err")
			} else if err.Error() != errMsg {
				t.Fatalf("err must be '%s', but get '%s'", errMsg, err.Error())
			}
			for _, idx := range tt.want.mutIndex {
				tt.want.limiters[idx].remaining -= 100
			}
			if diff := diffFlow(tt.want, tt.flow); diff != "" {
				t.Fatalf("tt.want %v, diff = %v", tt.want, diff)
			}
		})
	}
}

func TestFlow_Remaining(t *testing.T) {
	cfg := mgmt.Config.RateLimit
	tests := []struct {
		name string
		flow   *Flow
		want string
	}{
		{
			name: "not log in",
			flow:   newFlow("192.168.1.1", ""),
			want: fmt.Sprintf("%v,%v", cfg.QueryLimit, cfg.MutLimit),
		},
		{
			name: "logged in",
			flow:   newFlow("192.168.1.1", "test@uexky.com"),
			want: fmt.Sprintf("%v,%v,%v,%v", cfg.QueryLimit, cfg.MutLimit,
				cfg.QueryLimit, cfg.MutLimit),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.flow.remaining(); got != tt.want {
				t.Errorf("Flow.Remaining() = %v, want %v", got, tt.want)
			}
		})
	}
}
*/
