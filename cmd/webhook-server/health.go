package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"time"
)

func readinessCheckHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serveReadinessCheckFunc(w, r)
	})
}

func serveReadinessCheckFunc(w http.ResponseWriter, r *http.Request) {
	log.Print("Handling readiness check request ... ")
	var writeErr error
	err := doDNSResolveCheck("kubernetes.default.svc", 50*time.Millisecond)
	if err != nil {
		log.Printf("Error handling readiness check request: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, writeErr = w.Write([]byte(err.Error()))
	} else {
		log.Printf("Readiness check request handled successfully")
		_, writeErr = w.Write([]byte("ok"))
	}

	if writeErr != nil {
		log.Printf("Could not write response: %v", writeErr)
	}
}

// DNSResolveCheck returns a Check that makes sure the provided host can resolve
// to at least one IP address within the specified timeout.
func doDNSResolveCheck(host string, timeout time.Duration) error {
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
