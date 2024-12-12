package main

import (
	"context"
	"net/http"
)

func handleRequest(req *http.Request) {
	// Do some work
	ctx := req.Context()
	uncancelableCtx := context.WithoutCancel(ctx)

	var logRequestDetails func(ctx context.Context, req *http.Request)
	go func() {
		// This logging operation won't be interrupted if the parent context is canceled.
		logRequestDetails(uncancelableCtx, req)
	}()
}
