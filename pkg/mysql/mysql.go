package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql"
)

const (
	_defaultMaxPoolSize  = 10
	_defaultConnAttempts = 10
	_defaultConnTimeout  = time.Second
)

// Mysql -.
type Mysql struct {
	maxPoolSize  int
	connAttempts int
	connTimeout  time.Duration

	Builder squirrel.StatementBuilderType
	DB      *sql.DB
}

// New -.
func New(url string, opts ...Option) (*Mysql, error) {
	m := &Mysql{
		maxPoolSize:  _defaultMaxPoolSize,
		connAttempts: _defaultConnAttempts,
		connTimeout:  _defaultConnTimeout,
	}

	// Custom options
	for _, opt := range opts {
		opt(m)
	}

	m.Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question)

	dsn := strings.TrimPrefix(url, "mysql://")

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("mysql - New - sql.Open: %w", err)
	}

	db.SetMaxOpenConns(m.maxPoolSize)
	db.SetMaxIdleConns(m.maxPoolSize)
	db.SetConnMaxLifetime(time.Hour)

	for m.connAttempts > 0 {
		err = db.Ping()
		if err == nil {
			break
		}

		log.Printf("MySQL is trying to connect, attempts left: %d", m.connAttempts)
		time.Sleep(m.connTimeout)
		m.connAttempts--
	}

	if err != nil {
		return nil, fmt.Errorf("mysql - New - connAttempts == 0: %w", err)
	}

	m.DB = db
	return m, nil
}

// Close -.
func (m *Mysql) Close() {
	if m.DB != nil {
		m.DB.Close()
	}
}
