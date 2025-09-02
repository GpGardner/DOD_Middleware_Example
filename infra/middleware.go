package infra

import (
	"context"
	"fmt"
	"time"
)

var (
	DURATION    = "duration"
	RETRY_COUNT = "retry_count"
	MASKED      = "masked"
)

type RepoOp[In any, Out any] func(ctx context.Context, input In) (OutputWithMeta[Out], error)
type Middleware[In any, Out any] func(RepoOp[In, Out]) RepoOp[In, Out]

type OutputWithMeta[T any] struct {
	Data T
	Meta map[string]interface{} // or structured fields for meta
}

func SetMeta(ctx context.Context, key string, value interface{}) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	meta, ok := ctx.Value("meta").(map[string]interface{})
	if !ok {
		meta = make(map[string]interface{})
		ctx = context.WithValue(ctx, "meta", meta)
	}
	meta[key] = value
	return ctx
}

func GetMeta(ctx context.Context, key string) (interface{}, bool) {
	if ctx == nil {
		return nil, false
	}
	meta, ok := ctx.Value("meta").(map[string]interface{})
	if !ok {
		return nil, false
	}
	val, exists := meta[key]
	return val, exists
}

// Per-request overrides
type ctxKey int

const (
	ckDisableLogging ctxKey = iota
	ckDisableTiming
	ckDisableOutputResult
	ckDisableMasking
	ckDisableTracing
)

func DisableLogging(ctx context.Context) context.Context {
	return context.WithValue(ctx, ckDisableLogging, true)
}
func DisableTiming(ctx context.Context) context.Context {
	return context.WithValue(ctx, ckDisableTiming, true)
}
func DisableOutputResult(ctx context.Context) context.Context {
	return context.WithValue(ctx, ckDisableOutputResult, true)
}
func DisableTracing(ctx context.Context) context.Context {
	return context.WithValue(ctx, ckDisableTracing, true)
}
func DisableMasking(ctx context.Context) context.Context {
	return context.WithValue(ctx, ckDisableMasking, true)
}
func DisableAll(ctx context.Context) context.Context {
	ctx = DisableLogging(ctx)
	ctx = DisableTiming(ctx)
	ctx = DisableOutputResult(ctx)
	ctx = DisableMasking(ctx)
	ctx = DisableTracing(ctx)
	return ctx
}
func EnableLogging(ctx context.Context) context.Context {
	return context.WithValue(ctx, ckDisableLogging, false)
}
func EnableTiming(ctx context.Context) context.Context {
	return context.WithValue(ctx, ckDisableTiming, false)
}
func EnableOutputResult(ctx context.Context) context.Context {
	return context.WithValue(ctx, ckDisableOutputResult, false)
}
func EnableMasking(ctx context.Context) context.Context {
	return context.WithValue(ctx, ckDisableMasking, false)
}
func EnableTracing(ctx context.Context) context.Context {
	return context.WithValue(ctx, ckDisableTracing, false)
}
func EnableAll(ctx context.Context) context.Context {
	ctx = EnableLogging(ctx)
	ctx = EnableTiming(ctx)
	ctx = EnableOutputResult(ctx)
	ctx = EnableMasking(ctx)
	ctx = EnableTracing(ctx)
	return ctx
}

func IsLoggingDisabled(ctx context.Context) bool {
	disabled, ok := ctx.Value(ckDisableLogging).(bool)
	return ok && disabled
}
func IsTimingDisabled(ctx context.Context) bool {
	disabled, ok := ctx.Value(ckDisableTiming).(bool)
	return ok && disabled
}
func IsOutputResultDisabled(ctx context.Context) bool {
	disabled, ok := ctx.Value(ckDisableOutputResult).(bool)
	return ok && disabled
}
func IsMaskingDisabled(ctx context.Context) bool {
	disabled, ok := ctx.Value(ckDisableMasking).(bool)
	return ok && disabled
}
func IsTracingDisabled(ctx context.Context) bool {
	disabled, ok := ctx.Value(ckDisableTracing).(bool)
	return ok && disabled
}
func IsRetryDisabled(ctx context.Context) bool {
	disabled, ok := ctx.Value(ckDisableTracing).(bool)
	return ok && disabled
}

// Gate composes a middleware but short-circuits to `next` when disabledFn(ctx) == true.
// Name : Gate
// [In any, Out any] : Generic types for input and output (any type)
// mw Middleware[In, Out] : Is a function that takes a RepoOp[In, Out] and returns a RepoOp[In, Out]
// disabledFn func(ctx context.Context) bool : A function that determines if the middleware should be disabled based on the context
// next RepoOp[In, Out] : The next middleware in the chain
// returns RepoOp[In, Out] : The composed middleware operation
func Gate[In any, Out any](mw Middleware[In, Out], disabledFn func(ctx context.Context) bool) Middleware[In, Out] {
	// returns a middleware that checks if it should be disabled
	return func(next RepoOp[In, Out]) RepoOp[In, Out] {
		// wraps the next middleware in a new middleware that checks if it should be disabled
		wrapped := mw(next)
		// returns a RepoOp that checks if the middleware is disabled
		return func(ctx context.Context, in In) (OutputWithMeta[Out], error) {
			// checks if the middleware is disabled
			if disabledFn != nil && disabledFn(ctx) {
				return next(ctx, in)
			}
			// if not disabled, calls the wrapped middleware
			return wrapped(ctx, in)
		}
	}
}

// Timer measures time around the ENTIRE downstream chain and reports via callback.
func Timer[In any, Out any]() Middleware[In, Out] {
	return func(next RepoOp[In, Out]) RepoOp[In, Out] {
		return func(ctx context.Context, input In) (OutputWithMeta[Out], error) {
			start := time.Now()
			out, err := next(ctx, input)
			if out.Meta == nil {
				out.Meta = make(map[string]interface{})
			}
			out.Meta[DURATION] = time.Since(start)
			return out, err
		}
	}
}

func Logging[In any, Out any](logger func(ctx context.Context, msg string)) Middleware[In, Out] {
	return func(next RepoOp[In, Out]) RepoOp[In, Out] {
		return func(ctx context.Context, input In) (OutputWithMeta[Out], error) {
			logger(ctx, fmt.Sprintf("🟢 [START] Operation\n  ↳ Input: %+v", input))

			out, err := next(ctx, input)

			// Log any known metadata
			if len(out.Meta) > 0 {
				logger(ctx, "[META] Collected metadata:")
				for k, v := range out.Meta {
					logger(ctx, fmt.Sprintf("  - %s: %v", k, v))
				}
			}

			// Log result
			if err != nil {
				logger(ctx, fmt.Sprintf("[END] Operation FAILED\n  ↳ Error: %v", err))
			} else {
				logger(ctx, fmt.Sprintf("[END] Operation SUCCESS\n  ↳ Output: %+v", out.Data))
			}

			return out, err
		}
	}
}

func Tracing[In any, Out any]() Middleware[In, Out] {
	return func(next RepoOp[In, Out]) RepoOp[In, Out] {
		return func(ctx context.Context, input In) (OutputWithMeta[Out], error) {
			// inject tracing logic here
			return next(ctx, input)
		}
	}
}

func Retry[In any, Out any](maxRetries int, retryDelay time.Duration) Middleware[In, Out] {
	return func(next RepoOp[In, Out]) RepoOp[In, Out] {
		return func(ctx context.Context, input In) (OutputWithMeta[Out], error) {
			var out OutputWithMeta[Out]
			var err error
			retries := 0

			for {
				out, err = next(ctx, input)
				if err == nil || retries >= maxRetries {
					break
				}
				retries++
				time.Sleep(retryDelay)
			}

			if out.Meta == nil {
				out.Meta = make(map[string]interface{})
			}
			out.Meta[RETRY_COUNT] = retries

			return out, err
		}
	}
}

// OutputResult processes the output of the RepoOp and logs or modifies it as needed.
func OutputResult[In any, Out any](callback func(output Out, meta map[string]interface{}, err error)) Middleware[In, Out] {
	return func(next RepoOp[In, Out]) RepoOp[In, Out] {
		return func(ctx context.Context, input In) (OutputWithMeta[Out], error) {
			out, err := next(ctx, input)
			if callback != nil {
				callback(out.Data, out.Meta, err)
			}
			return out, err
		}
	}
}

// MaskOutput applies a masking function to the output of the RepoOp.
func MaskOutput[In any, Out any](maskFunc func(output Out) Out) Middleware[In, Out] {
	return func(next RepoOp[In, Out]) RepoOp[In, Out] {
		return func(ctx context.Context, input In) (OutputWithMeta[Out], error) {
			// Call the next middleware or base operation
			out, err := next(ctx, input)
			if err == nil && maskFunc != nil {
				out.Data = maskFunc(out.Data)
				if out.Meta == nil {
					out.Meta = make(map[string]interface{})
				}
				out.Meta[MASKED] = true
			}
			return out, err
		}
	}
}
