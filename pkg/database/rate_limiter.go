package database

import (
	"fmt"
	"go.uber.org/ratelimit"
)

type RateLimitedDbmsConn struct {
	Driver
	limiter ratelimit.Limiter
}

// NewRateLimitedDbmsConn returns a new rate-limited dbms connection. Rps specifies the number of allowed requests per
// second for this dbms connection. If rps is equal to 0, it returns a connection that is not rate-limited.
// Rps cannot be a negative number.
func NewRateLimitedDbmsConn(dbmsConn Driver, rps int) (*RateLimitedDbmsConn, error) {
	if rps <= 0 {
		return nil, fmt.Errorf("rps cannot be a negative number. Rps found: %d", rps)
	}
	var limiter ratelimit.Limiter
	if rps == 0 {
		limiter = ratelimit.NewUnlimited()
	} else {
		limiter = ratelimit.New(rps)
	}

	return &RateLimitedDbmsConn{
		Driver:  dbmsConn,
		limiter: limiter,
	}, nil
}

func (conn *RateLimitedDbmsConn) CreateDb(operation Operation) OpOutput {
	conn.limiter.Take()
	return conn.Driver.CreateDb(operation)
}

func (conn *RateLimitedDbmsConn) DeleteDb(operation Operation) OpOutput {
	conn.limiter.Take()
	return conn.Driver.CreateDb(operation)
}

func (conn *RateLimitedDbmsConn) Ping() error {
	conn.limiter.Take()
	return conn.Driver.Ping()
}
