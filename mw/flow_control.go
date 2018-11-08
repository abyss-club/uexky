package mw

import (
	"context"
	"fmt"
	"log"
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
		email := ""
		if ip == "" {
			httpErrorf(w, http.StatusBadRequest,
				"Not found header '%s'", remoteIPHeader)
			return
		}
		if userEmail, ok := req.Context().Value(ContextKeyEmail).(string); ok {
			email = userEmail
		}
		fc := newFlowController(ip, email)
		req = reqWithValue(req, ContextKeyFlowControl, fc)

		handle(w, req, p)

		remaining := fc.remaining()
		w.Header().Set("RateLimitRemaining", remaining)
	}
}

// FlowCostQuery ...
func FlowCostQuery(ctx context.Context, count int) error {
	fc, ok := ctx.Value(ContextKeyFlowControl).(*flowController)
	if !ok {
		log.Fatal("Can't find flow controller in context")
	}
	return fc.costQuery(ctx, count)
}

// FlowCostMut ...
func FlowCostMut(ctx context.Context, count int) error {
	fc, ok := ctx.Value(ContextKeyFlowControl).(*flowController)
	if !ok {
		log.Fatal("Can't find flow controller in context")
	}
	return fc.costMut(ctx, count)
}

// flowController manage flowcontrol amount
type flowController struct {
	ip         string
	email      string
	limiters   []*limiter
	queryIndex []int
	mutIndex   []int
}

// newFlowController ...
func newFlowController(ip, email string) *flowController {
	fc := &flowController{ip: ip, email: email}
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

// costQuery ...
func (fc *flowController) costQuery(ctx context.Context, count int) error {
	exceeded := false
	for _, idx := range fc.queryIndex {
		e, err := fc.limiters[idx].cost(ctx, count)
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

// costMut ...
func (fc *flowController) costMut(ctx context.Context, count int) error {
	exceeded := false
	for _, idx := range fc.mutIndex {
		e, err := fc.limiters[idx].cost(ctx, count)
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
func (fc *flowController) remaining() string {
	strs := []string{}
	for _, l := range fc.limiters {
		strs = append(strs, l.getRemaining())
	}
	return strings.Join(strs, ",")
}

func (fc *flowController) ipKey() string {
	return fmt.Sprintf("fc-ip-%s", fc.ip)
}

func (fc *flowController) emailKey() string {
	return fmt.Sprintf("fc-email-%s", fc.email)
}

func (fc *flowController) ipMutKey() string {
	return fmt.Sprintf("fc-ip-m-%s", fc.ip)
}

func (fc *flowController) emailMutKey() string {
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

// bool: return true if rate limit exceeded
func (l *limiter) cost(ctx context.Context, count int) (bool, error) {
	rd := GetRedis(ctx)
	l.count += count
	if l.count < l.ratio {
		return false, nil
	}

	cost := l.count / l.ratio
	l.count -= cost * l.ratio
	if _, err := rd.Do("SET", l.key, l.limit, "EX", l.expire, "NX"); err != nil {
		return false, errors.Wrap(err, "set rate limit")
	}
	remaining, err := redis.Int(rd.Do("DECRBY", l.key, cost))
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
