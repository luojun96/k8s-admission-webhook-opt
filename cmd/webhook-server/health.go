package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"time"
)

type doHealthCheck func() error

func newHealthCheckServer() *http.Server {
	mux := http.NewServeMux()
	mux.Handle("/ready", healthCheckHandler(doDNSResolveCheck))
	mux.Handle("/live", healthCheckHandler(doLivenessCheck))
	return &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
}

func healthCheckHandler(check doHealthCheck) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serveHealthCheckFunc(w, r, check)
	})
}

func serveHealthCheckFunc(w http.ResponseWriter, r *http.Request, check doHealthCheck) {
	var writeErr error
	err := check()
	if err != nil {
		log.Printf("Error handling health check request: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, writeErr = w.Write([]byte(err.Error()))
	} else {
		_, writeErr = w.Write([]byte("ok"))
	}

	if writeErr != nil {
		log.Printf("Could not write response: %v", writeErr)
	}
}

func doLivenessCheck() error {
	return nil
}

// DNSResolveCheck returns a Check that makes sure the provided host can resolve
// to at least one IP address within the specified timeout.
func doDNSResolveCheck() error {
	host := "kubernetes.default.svc"
	timeout := 50 * time.Millisecond
	resolver := net.Resolver{}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	addrs, err := resolver.LookupHost(ctx, host)
	if err != nil {
		return err
	}

	if len(addrs) < 1 {
		return errors.New("Could not resolve host")
	}
	return nil
}
