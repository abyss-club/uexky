package uexky

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"gitlab.com/abyss.club/uexky/config"
)

func diffFlow(l, r *FlowImpl) string {
	if l.u != r.u {
		return "uexky is not same"
	}
	if l.ip != r.ip {
		return fmt.Sprintf("ip: -%s\nip: +%s\n", l.ip, r.ip)
	}
	if l.email != r.email {
		return fmt.Sprintf("email: -%s\nemail: +%s\n", l.email, r.email)
	}
	if diff := cmp.Diff(l.limiters, r.limiters, cmp.AllowUnexported(limiter{})); diff != "" {
		return fmt.Sprintf("limiters' diff: %s", diff)
	}
	if diff := cmp.Diff(l.queryIndex, r.queryIndex); diff != "" {
		return fmt.Sprintf("queryIndex' diff: %s", diff)
	}
	if diff := cmp.Diff(l.mutIndex, r.mutIndex); diff != "" {
		return fmt.Sprintf("mutIndex' diff: %s", diff)
	}
	return ""
}

func TestNewFlow(t *testing.T) {
	config.LoadConfig("")
	config.ReplaceConfigByEnv()
	pool := InitPool()
	u := pool.NewUexky()

	type args struct {
		ip    string
		email string
	}
	cfg := config.Config.RateLimit
	tests := []struct {
		name string
		args args
		want *FlowImpl
	}{
		{
			name: "not logged in",
			args: args{
				ip:    "192.168.1.1",
				email: "",
			},
			want: &FlowImpl{
				u:     u,
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
			want: &FlowImpl{
				u:     u,
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
			got := NewUexkyFlow(u, tt.args.ip, tt.args.email)
			if diff := diffFlow(got, tt.want); diff != "" {
				t.Errorf("NewUexkyFlow() want %v, diff = %s", tt.want, diff)
			}
		})
	}
}

func TestFlow_CostQuery(t *testing.T) {
	config.LoadConfig("")
	config.ReplaceConfigByEnv()
	pool := InitPool()
	u := pool.NewUexky()
	errMsg := "rate limit exceeded"

	tests := []struct {
		name string
		flow *FlowImpl
		want *FlowImpl
	}{
		{
			name: "not login",
			flow: NewUexkyFlow(u, "192.168.1.1", ""),
			want: NewUexkyFlow(u, "192.168.1.1", ""),
		},
		{
			name: "logged in",
			flow: NewUexkyFlow(u, "192.168.1.1", "test@uexky.com"),
			want: NewUexkyFlow(u, "192.168.1.1", "test@uexky.com"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u.Redis.Do("FLUSHDB")
			t.Log("cost small query")
			tt.flow.CostQuery(5)
			for _, idx := range tt.want.queryIndex {
				tt.want.limiters[idx].count = 5
			}
			if diff := diffFlow(tt.want, tt.flow); diff != "" {
				t.Fatalf("tt.want %v, diff = %v", tt.want, diff)
			}

			t.Log("cost big query")
			tt.flow.CostQuery(15)
			for _, idx := range tt.want.queryIndex {
				tt.want.limiters[idx].count = 0
				tt.want.limiters[idx].remaining -= 2
			}
			if diff := diffFlow(tt.want, tt.flow); diff != "" {
				t.Fatalf("tt.want %v, diff = %v", tt.want, diff)
			}

			t.Log("let rate limit exceeded")
			if err := tt.flow.CostQuery(3000); err == nil {
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
	config.LoadConfig("")
	config.ReplaceConfigByEnv()
	pool := InitPool()
	u := pool.NewUexky()
	errMsg := "rate limit exceeded"

	tests := []struct {
		name string
		flow *FlowImpl
		want *FlowImpl
	}{
		{
			name: "not login",
			flow: NewUexkyFlow(u, "192.168.1.1", ""),
			want: NewUexkyFlow(u, "192.168.1.1", ""),
		},
		{
			name: "logged in",
			flow: NewUexkyFlow(u, "192.168.1.1", "test@uexky.com"),
			want: NewUexkyFlow(u, "192.168.1.1", "test@uexky.com"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u.Redis.Do("FLUSHDB")
			t.Log("cost normal")
			tt.flow.CostMut(5)
			for _, idx := range tt.want.mutIndex {
				tt.want.limiters[idx].remaining -= 5
			}
			if diff := diffFlow(tt.want, tt.flow); diff != "" {
				t.Fatalf("tt.want %v, diff = %v", tt.want, diff)
			}

			t.Log("let rate limit exceeded")
			if err := tt.flow.CostMut(100); err == nil {
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
	config.LoadConfig("")
	config.ReplaceConfigByEnv()
	pool := InitPool()
	u := pool.NewUexky()
	cfg := config.Config.RateLimit
	tests := []struct {
		name string
		flow *FlowImpl
		want string
	}{
		{
			name: "not log in",
			flow: NewUexkyFlow(u, "192.168.1.1", ""),
			want: fmt.Sprintf("%v,%v", cfg.QueryLimit, cfg.MutLimit),
		},
		{
			name: "logged in",
			flow: NewUexkyFlow(u, "192.168.1.1", "test@uexky.com"),
			want: fmt.Sprintf("%v,%v,%v,%v", cfg.QueryLimit, cfg.MutLimit,
				cfg.QueryLimit, cfg.MutLimit),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.flow.Remaining(); got != tt.want {
				t.Errorf("Flow.Remaining() = %v, want %v", got, tt.want)
			}
		})
	}
}
