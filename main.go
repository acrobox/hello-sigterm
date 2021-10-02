package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

const timeout = 10 * time.Second

func main() {
	mux := http.NewServeMux()
	srv := &http.Server{Addr: ":8080", Handler: mux}
	mux.HandleFunc("/", index)
	done := make(chan struct{})
	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
		<-quit
		ctx := context.Background()
		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()
		err := srv.Shutdown(ctx)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
		}
		close(done)
	}()
	err := srv.ListenAndServe()
	if err != http.ErrServerClosed {
		fmt.Fprintf(os.Stderr, "%v\n", err)
	}
	<-done
}

func index(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	n := req.URL.Path[1:]
	if n == "" {
		n = "0"
	}
	s, err := strconv.Atoi(n)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%v\n", err)
		return
	}
	d := time.Duration(s)
	time.Sleep(d * time.Second)
	fmt.Fprintf(w, "Hello, 世界\n")
}
