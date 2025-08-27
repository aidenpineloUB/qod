// // Filename: cmd/api/server.go

package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"
)


// serve starts the HTTP server (server.go)
func (app *applicationDependencies) serve() error {
       srv := &http.Server {
       Addr:         fmt.Sprintf(":%d", app.config.port),
       Handler:      app.routes(),
       IdleTimeout:  time.Minute,
       ReadTimeout:  5 * time.Second,
       WriteTimeout: 10 * time.Second,
       ErrorLog:     slog.NewLogLogger(app.logger.Handler(), slog.LevelError),
    }
	app.logger.Info("starting server", "addr", srv.Addr, "env", app.config.environment)
	return srv.ListenAndServe()
}
