package middleware

import (
	"fmt"
	"github.com/kubemq-hub/kubemq-bridges/config"
	"github.com/kubemq-hub/kubemq-bridges/pkg/ratelimit"

	"math"
)

type RateLimitMiddleware struct {
	rateLimiter ratelimit.Limiter
}

func NewRateLimitMiddleware(meta config.Metadata) (*RateLimitMiddleware, error) {
	rpc, err := meta.ParseIntWithRange("rate_per_second", 0, 0, math.MaxInt32)
	if err != nil {
		return nil, fmt.Errorf("invalid rate limiter rate per second value, %w", err)
	}
	rl := &RateLimitMiddleware{}
	if rpc > 0 {
		rl.rateLimiter = ratelimit.New(rpc, ratelimit.WithoutSlack)
	} else {
		rl.rateLimiter = ratelimit.NewUnlimited()
	}
	return rl, nil
}

func (rl *RateLimitMiddleware) Take() {
	_ = rl.rateLimiter.Take()

}
