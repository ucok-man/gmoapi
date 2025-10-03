package main

import (
	"context"
	"database/sql"
	"expvar"
	"log/slog"
	"os"
	"runtime"
	"sync"
	"time"

	_ "github.com/ucok-man/gmoapi/cmd/api/docs"

	_ "github.com/lib/pq"
	"github.com/ucok-man/gmoapi/cmd/api/config"
	"github.com/ucok-man/gmoapi/internal/data"
	"github.com/ucok-man/gmoapi/internal/mailer"
)

type application struct {
	config config.Config
	logger *slog.Logger
	models data.Models
	mailer *mailer.Mailer
	wg     sync.WaitGroup
}

// @title           Gmoapi - Movie Management API
// @version         1.0.1
// @description     A production-ready RESTful API for managing movies with comprehensive user authentication, role-based authorization, rate limiting, and email notifications.
// @description
// @description     ## Features
// @description     - Full CRUD operations for movies
// @description     - User registration and authentication
// @description     - Role-based access control (RBAC)
// @description     - Token-based authentication (Bearer)
// @description     - Email verification and password reset
// @description     - Rate limiting (2 req/s, burst: 4)
// @description     - Pagination and filtering
// @description     - CORS support
// @description
// @description     ## Authentication
// @description     Most endpoints require authentication. Use the `/v1/tokens/authentication` endpoint to obtain a token, then include it in the Authorization header as: `Bearer YOUR_TOKEN`
// @description
// @description     ## Rate Limiting
// @description     API requests are rate-limited to 2 requests per second with a burst of 4 requests.

// @contact.name   ucokman
// @contact.url    https://github.com/ucok-man/gmoapi
// @contact.email  support@ucokman.web.id

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:4000
// @BasePath  /v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Enter your bearer token in the format: Bearer {token}

// @tag.name Health
// @tag.description System health check and version information

// @tag.name Movies
// @tag.description Movie catalog management - requires authentication and appropriate permissions

// @tag.name Users
// @tag.description User account registration, activation, and password management

// @tag.name Tokens
// @tag.description Token generation for authentication, activation, and password reset

// @x-extension-openapi {"example": "value"}
func main() {
	cfg := config.MustNewConfig()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	defer db.Close()

	logger.Info("database connection pool established")

	mailer, err := mailer.New(cfg.SMTP.Host, cfg.SMTP.Port, cfg.SMTP.Username, cfg.SMTP.Password, cfg.SMTP.Sender)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	expvar.NewString("version").Set(config.APP_VERSION)

	// Publish the number of active goroutines.
	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	// Publish the database connection pool statistics.
	expvar.Publish("database", expvar.Func(func() any {
		return db.Stats()
	}))

	// Publish the current Unix timestamp.
	expvar.Publish("timestamp", expvar.Func(func() any {
		return time.Now().Unix()
	}))

	app := &application{
		config: cfg,
		logger: logger,
		models: data.NewModels(db),
		mailer: mailer,
	}

	err = app.serve()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
}

func openDB(cfg config.Config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.DB.DSN)
	if err != nil {
		return nil, err
	}

	// Set the maximum number of open (in-use + idle) connections in the pool. Note that
	// passing a value less than or equal to 0 will mean there is no limit.
	db.SetMaxOpenConns(cfg.DB.MaxOpenConn)

	// Set the maximum number of idle connections in the pool. Again, passing a value
	// less than or equal to 0 will mean there is no limit.
	db.SetMaxIdleConns(cfg.DB.MaxIdleConn)

	// Set the maximum idle timeout for connections in the pool. Passing a duration less
	// than or equal to 0 will mean that connections are not closed due to their idle time.
	db.SetConnMaxIdleTime(cfg.DB.MaxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Ping untuk benar-benar buka koneksi
	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
