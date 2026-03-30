// Package postgres provides support for access the postgres database.
package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

const (
	uniqueViolation = "23505"
	undefinedTable  = "42P01"
)

// Set of error variables for CRUD operations.
var (
	ErrDBNotFound        = sql.ErrNoRows
	ErrDBDuplicatedEntry = errors.New("duplicated entry")
	ErrUndefinedTable    = errors.New("undefined table")
)

// Config is the required properties to use the database.
type Config struct {
	User            string
	Password        string
	Host            string
	Port            string
	Name            string
	Schema          string
	MaxIdleConns    int
	MaxOpenConns    int
	IdleConnTimeout time.Duration
	EnableTLS       bool
	CACert          string
	ClientCert      string
	ClientKey       string
	ApplicationName string
}

// Open knows how to open a database connection based on the configuration.
func Open(cfg Config) (*sqlx.DB, error) {
	sslMode, err := getSSLMode(cfg)
	if err != nil {
		return nil, err
	}

	q := make(url.Values)
	q.Set("sslmode", sslMode)
	q.Set("timezone", "utc")
	q.Set("application_name", cfg.ApplicationName)

	if cfg.EnableTLS {
		q.Set("sslrootcert", cfg.CACert)
		q.Set("sslcert", cfg.ClientCert)
		q.Set("sslkey", cfg.ClientKey)
	}

	if cfg.Schema != "" {
		q.Set("search_path", cfg.Schema)
	}

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     cfg.Host + ":" + cfg.Port,
		Path:     cfg.Name,
		RawQuery: q.Encode(),
	}

	db, err := sqlx.Open("pgx", u.String())
	if err != nil {
		return nil, err
	}
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetConnMaxIdleTime(cfg.IdleConnTimeout)

	// Status check

	t := time.Second * 5
	ctx, cancel := context.WithTimeout(context.Background(), t)
	defer cancel()

	if err := StatusCheck(ctx, db); err != nil {
		return nil, fmt.Errorf("database status check: %w", err)
	}

	return db, nil
}

func getSSLMode(cfg Config) (string, error) {
	if !cfg.EnableTLS {
		return "disable", nil
	}

	if cfg.CACert == "" || cfg.ClientCert == "" || cfg.ClientKey == "" {
		return "", fmt.Errorf("SSL certificates not properly configured")
	}

	return "require", nil
}

// RunQuery is a helper function for executing queries that return a
// single value to be unmarshalled into a struct type.
func RunQuery(ctx context.Context, db sqlx.ExtContext, query string, dest any) error {
	var rows *sqlx.Rows
	var err error

	rows, err = sqlx.NamedQueryContext(ctx, db, query, struct{}{})

	if err != nil {
		if pqerr, ok := err.(*pgconn.PgError); ok && pqerr.Code == undefinedTable {
			return ErrUndefinedTable
		}
		return err
	}
	defer rows.Close()

	if !rows.Next() {
		return ErrDBNotFound
	}

	if err := rows.StructScan(dest); err != nil {
		return err
	}

	return nil
}

// RunQuerySlice is a helper function for executing queries that return a
// collection of data to be unmarshalled into a slice.
func RunQuerySlice[T any](ctx context.Context, db sqlx.ExtContext, query string, dest *[]T) error {
	var rows *sqlx.Rows
	var err error

	rows, err = sqlx.NamedQueryContext(ctx, db, query, struct{}{})

	if err != nil {
		if pqerr, ok := err.(*pgconn.PgError); ok && pqerr.Code == undefinedTable {
			return ErrUndefinedTable
		}
		return err
	}
	defer rows.Close()

	var slice []T
	for rows.Next() {
		v := new(T)
		if err := rows.StructScan(v); err != nil {
			return err
		}
		slice = append(slice, *v)
	}
	*dest = slice

	return nil
}

// RunCUD is a helper function to execute create, update, or delete operation.
func RunCUD(ctx context.Context, db sqlx.ExtContext, query string, data any) (any, error) {
	result, err := sqlx.NamedExecContext(ctx, db, query, data)
	if err != nil {
		if pqerr, ok := err.(*pgconn.PgError); ok {
			switch pqerr.Code {
			case undefinedTable:
				return "", ErrUndefinedTable
			case uniqueViolation:
				return "", ErrDBDuplicatedEntry
			}
		}
		return "", err
	}

	return result, nil
}

// QueryRunCUD is a helper function to execute create, update, or delete operation.
func QueryRunCUD(ctx context.Context, db sqlx.ExtContext, query string, data any, dest any) error {
	var rows *sqlx.Rows
	var err error

	rows, err = sqlx.NamedQueryContext(ctx, db, query, data)

	if err != nil {
		if pqerr, ok := err.(*pgconn.PgError); ok && pqerr.Code == undefinedTable {
			return ErrUndefinedTable
		}
		return err
	}
	defer rows.Close()

	if !rows.Next() {
		return ErrDBNotFound
	}

	if err := rows.StructScan(dest); err != nil {
		return err
	}

	return nil
}

// StatusCheck returns nil if it can successfully talk to the database. It
// returns a non-nil error otherwise.
func StatusCheck(ctx context.Context, db *sqlx.DB) error {
	if _, ok := ctx.Deadline(); !ok {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, time.Second*5)
		defer cancel()
	}

	var pingError error
	for attempts := 1; ; attempts++ {
		pingError = db.PingContext(ctx)
		if pingError == nil {
			break
		}
		time.Sleep(time.Duration(attempts) * 1 * time.Second)
		if ctx.Err() != nil {
			return fmt.Errorf("%w : database: %w", ctx.Err(), pingError)
		}
	}

	if ctx.Err() != nil {
		return ctx.Err()
	}

	// Run a simple query to determine connectivity.
	// Running this query forces a round trip through the database.
	const q = `SELECT true`
	var tmp bool
	return db.QueryRowContext(ctx, q).Scan(&tmp)
}

// ParseQuery provides a pretty version of the query and parameters.
func ParseQuery(query string, args any) string {
	query, params, err := sqlx.Named(query, args)
	if err != nil {
		return err.Error()
	}

	for _, param := range params {
		var value string
		switch v := param.(type) {
		case string:
			value = fmt.Sprintf("'%s'", v)
		case []byte:
			value = fmt.Sprintf("'%s'", string(v))
		case uuid.UUID:
			value = fmt.Sprintf("'%s'", v)
		default:
			value = fmt.Sprintf("%v", v)
		}
		query = strings.Replace(query, "?", value, 1)
	}

	query = strings.ReplaceAll(query, "\t", "")
	query = strings.ReplaceAll(query, "\n", " ")

	return strings.Trim(query, " ")
}
