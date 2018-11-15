package uexky

import (
	"fmt"
	"log"
	"strings"

	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
	"gitlab.com/abyss.club/uexky/mgmt"
)

// NewUexkyFlow make a new Flow, and add to Uexky
func NewUexkyFlow(u *Uexky, ip, email string) *Flow {
	flow := &Flow{u: u, ip: ip, email: email}
	u.Flow = flow

	cfg := &mgmt.Config.RateLimit
	log.Printf("DEBUG!!! ratelimit cfg is %+v", cfg)
	flow.limiters = []*limiter{
		newLimiter(flow.ipKey(), cfg.QueryLimit, cfg.QueryResetTime, 10),
		newLimiter(flow.ipMutKey(), cfg.MutLimit, cfg.MutResetTime, 1),
	}

	if email == "" {
		flow.queryIndex = []int{0}
		flow.mutIndex = []int{1}
		return flow
	}
	flow.limiters = append(flow.limiters,
		newLimiter(flow.emailKey(), cfg.QueryLimit, cfg.QueryResetTime, 10),
		newLimiter(flow.emailMutKey(), cfg.MutLimit, cfg.MutResetTime, 1),
	)
	flow.queryIndex = []int{0, 2}
	flow.mutIndex = []int{1, 3}
	return flow
}

// Flow manage tool
type Flow struct {
	u          *Uexky
	ip         string
	email      string
	limiters   []*limiter
	queryIndex []int
	mutIndex   []int
}

// CostQuery ...
func (flow *Flow) CostQuery(count int) error {
	log.Printf("DEBUG!!! Flow getRemaining = %v", flow.remaining())
	exceeded := false
	for _, idx := range flow.queryIndex {
		e, err := flow.limiters[idx].cost(flow.u, count)
		if err != nil {
			return err
		}
		exceeded = exceeded || e
	}
	if exceeded {
		return errors.New("rate limit exceeded")
	}
	return nil
}

// CostMut ...
func (flow *Flow) CostMut(count int) error {
	log.Printf("DEBUG!!! Flow getRemaining = %v", flow.remaining())
	exceeded := false
	for _, idx := range flow.mutIndex {
		e, err := flow.limiters[idx].cost(flow.u, count)
		if err != nil {
			return err
		}
		exceeded = exceeded || e
	}
	if exceeded {
		return errors.New("rate limit exceeded")
	}
	return nil
}

// remaining ...
func (flow *Flow) remaining() string {
	strs := []string{}
	for _, l := range flow.limiters {
		strs = append(strs, l.getRemaining())
	}
	return strings.Join(strs, ",")
}

func (flow *Flow) ipKey() string {
	return fmt.Sprintf("flow-ip-%s", flow.ip)
}

func (flow *Flow) emailKey() string {
	return fmt.Sprintf("flow-email-%s", flow.email)
}

func (flow *Flow) ipMutKey() string {
	return fmt.Sprintf("flow-ip-m-%s", flow.ip)
}

func (flow *Flow) emailMutKey() string {
	return fmt.Sprintf("flow-email-m-%s", flow.email)
}

type limiter struct {
	// setting
	key    string
	limit  int
	ratio  int
	expire int

	// runtime
	count     int // count of not-dealed
	remaining int
}

func newLimiter(key string, limit, expire, ratio int) *limiter {
	return &limiter{key, limit, ratio, expire, 0, limit}
}

// bool: return true if rate limit exceeded
func (l *limiter) cost(u *Uexky, count int) (bool, error) {
	l.count += count
	if l.count < l.ratio {
		return false, nil
	}

	cost := l.count / l.ratio
	l.count -= cost * l.ratio
	if _, err := u.Redis.Do("SET", l.key, l.limit, "EX", l.expire, "NX"); err != nil {
		return false, errors.Wrap(err, "set rate limit")
	}
	remaining, err := redis.Int(u.Redis.Do("DECRBY", l.key, cost))
	if err != nil {
		return false, errors.Wrap(err, "cost flow control")
	}
	l.remaining = remaining
	if l.remaining < 0 {
		return true, nil
	}
	return false, nil
}

func (l *limiter) getRemaining() string {
	return fmt.Sprint(l.remaining)
}
