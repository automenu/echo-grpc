// Copyright 2024 The Echo gRPC Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"errors"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/automenu/echo-grpc/echo"
	"go.uber.org/zap"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	// Initialize global zap logger
	zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
	if strings.EqualFold(os.Getenv("DEBUG"), "true") {
		zap.ReplaceGlobals(zap.Must(zap.NewDevelopment()))
	}
	defer func() {
		_ = zap.L().Sync()
	}()

	// Create initial context and cancel it on exit for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create HTTP request multiplexer and register API handlers
	mux := http.NewServeMux()
	mux.Handle(newHTTPHealthCheckHandler())
	mux.Handle(echo.NewEchoAPIHandler())

	addr := "localhost:8080"
	if port := os.Getenv("PORT"); port != "" {
		addr = "localhost:" + port
	}
	zap.L().Info("Starting server", zap.String("addr", addr))

	// Create HTTP server
	srv := &http.Server{
		Addr:              addr,
		Handler:           h2c.NewHandler(mux, &http2.Server{}),
		BaseContext:       func(net.Listener) context.Context { return ctx },
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       15 * time.Second,
		MaxHeaderBytes:    8 * 1024, // 8KB
	}

	// Create a tcp keep-alive listener
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		zap.L().Fatal("net.Listen", zap.Error(err))
	}
	keepAliveListener := tcpKeepAliveListener{ln.(*net.TCPListener)}

	// Start HTTP server
	go func() {
		err := srv.Serve(keepAliveListener)
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			zap.L().Fatal("HTTP listen and serve", zap.Error(err))
		}
	}()

	// Graceful shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM)
	<-signals

	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		zap.L().Fatal("HTTP shutdown", zap.Error(err))
	}
	zap.L().Info("Server shutdown")
}

// tcpKeepAliveListener sets TCP keep-alive timeouts on accepted
// connections. It's used by ListenAndServe and ListenAndServeTLS so
// dead TCP connections (e.g. closing laptop mid-download) eventually
// go away.
type tcpKeepAliveListener struct {
	*net.TCPListener
}

func (ln tcpKeepAliveListener) Accept() (c net.Conn, err error) {
	tc, err := ln.AcceptTCP()
	if err != nil {
		return nil, err
	}
	if err := tc.SetKeepAlive(true); err != nil {
		return nil, err
	}
	if err := tc.SetKeepAlivePeriod(30 * time.Second); err != nil {
		return nil, err
	}
	return tc, nil
}

func newHTTPHealthCheckHandler() (string, http.Handler) {
	return "/healthz", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	})
}
