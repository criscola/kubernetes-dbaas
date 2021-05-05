package database

import (
	"fmt"
	"go.uber.org/ratelimit"
)

type RateLimitedDbmsConn struct {
	Driver
	limiter ratelimit.Limiter
}

func NewRateLimitedDbmsConn(driver string, dsn Dsn, rps int) (*RateLimitedDbmsConn, error) {
	if rps <= 0 {
		return nil, fmt.Errorf("rps cannot be less than or equal to 0. Rps found: %d", rps)
	}
	dbmsConn, err := NewDbmsConn(driver, dsn)
	if err != nil {
		return nil, err
	}
	limiter := ratelimit.New(rps)
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
