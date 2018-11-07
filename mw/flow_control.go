package mw

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gomodule/redigo/redis"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
	"gitlab.com/abyss.club/uexky/mgmt"
)

const remoteIPHeader = "Remote-IP"

// WithFlowControl is middleware for rate limit
func WithFlowControl(handle httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		ip := req.Header.Get(remoteIPHeader)
		if ip == "" {
			httpErrorf(w, http.StatusBadRequest,
				"Not found header '%s'", remoteIPHeader)
			return
		}
		flc := &FlowController{ip: ip}
		if email, ok := req.Context().Value(ContextKeyEmail).(string); ok {
			flc.email = email
		}
		req = reqWithValue(req, ContextKeyFlowControl, flc)

		handle(w, req, p)

		remaining := flc.Remaining()
		w.Header().Set("RateLimitRemaining", remaining)
	}
}

// FlowController manage flowcontrol amount
type FlowController struct {
	ip         string
	email      string
	limiters   []*limiter
	queryIndex []int
	mutIndex   []int
}

// NewFlowController ...
func NewFlowController(ip, email string) *FlowController {
	fc := &FlowController{ip: ip, email: email}
	cfg := &mgmt.Config.RateLimit
	fc.limiters = []*limiter{
		newLimiter(fc.ipKey(), cfg.QueryLimit, cfg.QueryResetTime, 10),
		newLimiter(fc.ipMutKey(), cfg.MutLimit, cfg.MutResetTime, 1),
	}

	if email == "" {
		fc.queryIndex = []int{0}
		fc.mutIndex = []int{1}
		return fc
	}
	fc.limiters = append(fc.limiters,
		newLimiter(fc.emailKey(), cfg.QueryLimit, cfg.QueryResetTime, 10),
		newLimiter(fc.emailMutKey(), cfg.MutLimit, cfg.MutResetTime, 1),
	)
	fc.queryIndex = []int{0, 2}
	fc.mutIndex = []int{1, 3}
	return fc
}

// CostQuery ...
func (fc *FlowController) CostQuery(ctx context.Context, count int) error {
	for _, idx := range fc.queryIndex {
		if err := fc.limiters[idx].cost(ctx, count); err != nil {
			return err
		}
	}
	return nil
}

// CostMut ...
func (fc *FlowController) CostMut(ctx context.Context, count int) error {
	for _, idx := range fc.mutIndex {
		if err := fc.limiters[idx].cost(ctx, count); err != nil {
			return err
		}
	}
	return nil
}

// Remaining ...
func (fc *FlowController) Remaining() string {
	strs := []string{}
	for _, l := range fc.limiters {
		strs = append(strs, l.getRemaining())
	}
	return strings.Join(strs, ",")
}

func (fc *FlowController) ipKey() string {
	return fmt.Sprintf("fc-ip-%s", fc.ip)
}

func (fc *FlowController) emailKey() string {
	return fmt.Sprintf("fc-email-%s", fc.email)
}

func (fc *FlowController) ipMutKey() string {
	return fmt.Sprintf("fc-ip-m-%s", fc.ip)
}

func (fc *FlowController) emailMutKey() string {
	return fmt.Sprintf("fc-email-m-%s", fc.email)
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

func (l *limiter) cost(ctx context.Context, count int) error {
	rd := GetRedis(ctx)
	l.count += count
	if l.count < l.ratio {
		return nil
	}

	cost := l.count / l.ratio
	l.count -= cost * l.ratio
	if _, err := rd.Do("SET", l.key, l.limit, "EX", l.expire, "NX"); err != nil {
		return errors.Wrap(err, "set rate limit")
	}
	remaining, err := redis.Int(rd.Do("DECRBY", l.key, cost))
	if err != nil {
		return errors.Wrap(err, "cost flow control")
	}
	l.remaining = remaining
	if l.remaining < 0 {
		return errors.New("rate limit exceeded")
	}
	return nil
}

func (l *limiter) getRemaining() string {
	return fmt.Sprint(l.remaining)
}
