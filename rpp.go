package rpp

import (
	"errors"
	"net/url"
	"strconv"
	"strings"

	"github.com/garyburd/redigo/redis"
)

var (
	// ErrNoHostnameProvided is returned when a provided Redis URI does not
	// contain a hostname.
	ErrNoHostnameProvided = errors.New("no hostname provided")

	dialFn = func(host string) (redis.Conn, error) {
		conn, err := redis.Dial("tcp", host)
		return conn, err
	}
)

// RPP returns a Redis connection pool that simplifies creating a connection
// pool when the URI contains authentication and specific database information.
func RPP(redisURI string, maxActive, maxIdle int) (rp *redis.Pool, err error) {
	var uri *url.URL
	if uri, err = url.Parse(redisURI); err != nil {
		return
	}

	if uri.Host == "" {
		err = ErrNoHostnameProvided
		return
	}

	// Time to build the pool! The most important thing here is the `Dial`
	// function.
	rp = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := dialFn(uri.Host)
			if err != nil {
				return nil, err
			}

			// See if there is any userinfo that was set in the URI.
			if uri.User != nil {
				if password, ok := uri.User.Password(); ok {
					if _, err := c.Do("AUTH", password); err != nil {
						_ = c.Close()
						return nil, err
					}
				}
			}

			// See if a specific DB was specified in the connection URI.
			path := uri.Path
			if path != "" {
				dbNum, err := strconv.Atoi(strings.TrimPrefix(path, "/"))
				if err != nil {
					return nil, err
				}

				// If we have a DB num then, use it.
				if _, err := c.Do("SELECT", dbNum); err != nil {
					_ = c.Close()
					return nil, err
				}

			}

			// Return the setup connection.
			return c, nil
		},
		MaxActive: maxActive,
		MaxIdle:   maxIdle,
	}

	return
}
