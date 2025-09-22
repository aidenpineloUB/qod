// Filename: cmd/api/main.go
package main

import (
	"context"
	"database/sql"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/aidenpineloUB/qod/internal/data"
	_ "github.com/lib/pq"
)

const appVersion = "1.0.0"

type serverConfig struct {
	port        int
	environment string
	db          struct {
		dsn string
	}
	cors struct {
		trustedOrigins []string
	}
	    limiter struct {
        rps float64                      // requests per second
        burst int                        // initial requests possible
        enabled bool                     // enable or disable rate limiter
    }

}

type applicationDependencies struct {
	config       serverConfig
	logger       *slog.Logger
	commentModel data.CommentModel
}

func main() {
	var settings serverConfig
	
	// Define all flags BEFORE flag.Parse()
	flag.IntVar(&settings.port, "port", 4000, "Server port")
	flag.StringVar(&settings.environment, "env", "development",
		"Environment(development|staging|production)")
	flag.StringVar(&settings.db.dsn, "db-dsn", "postgres://quotes:quotes@localhost/quotes?sslmode=require", "PostgreSQL DSN")
	
	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
		settings.cors.trustedOrigins = strings.Fields(val)
		return nil
	})
	    flag.Float64Var(&settings.limiter.rps, "limiter-rps", 2,
                  "Rate Limiter maximum requests per second")

    flag.IntVar(&settings.limiter.burst, "limiter-burst", 5,
                  "Rate Limiter maximum burst")

    flag.BoolVar(&settings.limiter.enabled, "limiter-enabled", true,
                  "Enable rate limiter")

	
	// Parse flags ONCE
	flag.Parse()
	
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	
	db, err := openDB(settings)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}
	
	// release the database resources before exiting
	defer db.Close()
	logger.Info("database connection pool established")
	
	// Add this after: logger.Info("database connection pool established")
	err = testDatabaseWrite(db)
	if err != nil {
		logger.Error("Database write test failed", "error", err)
	} else {
		logger.Info("Database write test successful!")
	}
	
	appInstance := &applicationDependencies{
		config:       settings,
		logger:       logger,
		commentModel: data.CommentModel{DB: db},
	}
	
	apiServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", settings.port),
		Handler:      appInstance.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}
	
	// Use the serve method which properly handles routing
	logger.Info("starting server", "address", apiServer.Addr, "environment", settings.environment)
	
	err = apiServer.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}

func testDatabaseWrite(db *sql.DB) error {
	// Create a simple test table
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS test_table (
			id SERIAL PRIMARY KEY,
			message TEXT,
			created_at TIMESTAMP DEFAULT NOW()
		)
	`)
	if err != nil {
		return fmt.Errorf("failed to create test table: %v", err)
	}
	
	// Insert a test record
	_, err = db.Exec("INSERT INTO test_table (message) VALUES ($1)", "Hello from Go!")
	if err != nil {
		return fmt.Errorf("failed to insert test data: %v", err)
	}
	
	// Read it back
	var message string
	err = db.QueryRow("SELECT message FROM test_table ORDER BY id DESC LIMIT 1").Scan(&message)
	if err != nil {
		return fmt.Errorf("failed to read test data: %v", err)
	}
	
	fmt.Printf("Successfully wrote and read: %s\n", message)
	return nil
}

func openDB(settings serverConfig) (*sql.DB, error) {
	// open a connection pool
	db, err := sql.Open("postgres", settings.db.dsn)
	if err != nil {
		return nil, err
	}
	
	// set a context to ensure DB operations don't take too long
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	// let's test if the connection pool was created
	// we trying pinging it with a 5-second timeout
	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}
	
	// return the connection pool (sql.DB)
	return db, nil
}