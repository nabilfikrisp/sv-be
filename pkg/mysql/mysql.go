package mysql

import (
	"database/sql"
	"fmt"
	"log"
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

	// MySQL menggunakan Question Mark (?) sebagai placeholder, bukan Dollar ($)
	m.Builder = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question)

	db, err := sql.Open("mysql", url)
	if err != nil {
		return nil, fmt.Errorf("mysql - New - sql.Open: %w", err)
	}

	// Konfigurasi Connection Pool
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
