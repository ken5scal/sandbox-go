package middlewares

import (
	"context"
	"sync"
)

var (
	logNo = 1 //global
	mu    sync.Mutex
)

func newTraceID() int {
	var no int

	mu.Lock()
	no = logNo
	logNo++
	mu.Unlock()

	return no
}

func SetTraceID(ctx context.Context, traceID int) context.Context {
	return context.WithValue(ctx, "traceID", traceID)
}

func GetTraceID(ctx context.Context) int {
	id := ctx.Value("traceID")
	if idInt, ok := id.(int); ok {
		return idInt
	}
	return 0
}
