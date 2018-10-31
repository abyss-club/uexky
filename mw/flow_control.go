package mw

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gomodule/redigo/redis"
	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

const remoteIPHeader = "Remote-IP"

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

		if remaining, err := flc.remaining(req.Context()); err != nil {
			httpError(w, http.StatusInternalServerError,
				"read rete limit error")
			return
		} else {
			w.Header().Set("RateLimitRemaining", remaining)
		}
	}
}

// FlowController manage flowcontrol amount
type FlowController struct {
	ip    string
	email string
}

func (fc *FlowController) ipKey() string {
	return fmt.Sprintf("fc-ip-%s", fc.ip)
}

func (fc *FlowController) emailKey() string {
	return fmt.Sprintf("fc-email-%s", fc.email)
}

func (fc *FlowController) ipMutationKey() string {
	return fmt.Sprintf("fc-ip-m-%s", fc.ip)
}

func (fc *FlowController) emailMutationKey() string {
	return fmt.Sprintf("fc-email-m-%s", fc.email)
}

const expireSeconds = 3600

func costLimit(ctx context.Context, key string, init, count int) error {
	rd := GetRedis(ctx)
	if _, err := rd.Do("SET", key, init, "EX", expireSeconds, "NX"); err != nil {
		return errors.Wrap(err, "set rate limit")
	}
	if remaining, err := redis.Int(rd.Do("DECRBY", key, count)); err != nil {
		return errors.Wrap(err, "cost flow control")
	} else if remaining < 0 {
		return errors.New("Rate limit exceeded")
	}
	return nil
}

// Cost flowcontrol amount by count, if amount < 0, return error
func (fc *FlowController) Cost(ctx context.Context, count int) error {
	if err := costLimit(ctx, fc.ipKey(), 1000, count); err != nil {
		return err
	}
	if fc.email == "" {
		return nil
	}
	if err := costLimit(ctx, fc.emailKey(), 1000, count); err != nil {
		return err
	}
	return nil
}

// MutationCost cost rate limit when do mutation
func (fc *FlowController) MutationCost(ctx context.Context, count int) error {
	if fc.email == "" {
		return errors.New("you should login before do mutation")
	}
	if err := costLimit(ctx, fc.ipMutationKey(), 1000, count); err != nil {
		return err
	}
	if err := costLimit(ctx, fc.emailMutationKey(), 1000, count); err != nil {
		return err
	}
	return nil
}

// remaining of Rete limit:
// Not Login: "<queryByIP>"
// Loginned: "<queryByIP>,<queryByEmail>,<mutationByIP>,<mutationByEmail>
func (fc *FlowController) remaining(ctx context.Context) (string, error) {
	rd := GetRedis(ctx)
	keys := []string{fc.ipKey()}
	remainings := []string{}
	if fc.email != "" {
		keys = append(
			keys, fc.emailKey(), fc.ipMutationKey(), fc.emailMutationKey(),
		)
	}
	for _, key := range keys {
		if r, err := redis.Int(rd.Do("GET", key)); err != nil {
			return "", errors.New("query rate limit remaining")
		} else {
			remainings = append(remainings, fmt.Sprint(r))
		}
	}
	return strings.Join(remainings, ","), nil
}
