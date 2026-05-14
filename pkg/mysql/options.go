// Package mysql provides MySQL connection options.
package mysql

import "time"

// Option -.
type Option func(*Mysql)

// MaxPoolSize sets the maximum number of open connections to the database.
func MaxPoolSize(size int) Option {
	return func(m *Mysql) {
		m.maxPoolSize = size
	}
}

// ConnAttempts sets the number of connection attempts before giving up.
func ConnAttempts(attempts int) Option {
	return func(m *Mysql) {
		m.connAttempts = attempts
	}
}

// ConnTimeout sets the duration between connection attempts.
func ConnTimeout(timeout time.Duration) Option {
	return func(m *Mysql) {
		m.connTimeout = timeout
	}
}
