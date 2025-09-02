package infra

import (
	"context"
	"time"
)

type MiddlewareConfig struct {
	RetryCount int
	RetryDelay time.Duration

	TimerCallback   func()
	OutputCallback  func(output any, OutOutputWithMetaa map[string]interface{}, err error)
	MaskingCallback func(output any) any
	RetryCallback   func(attempt any, err error)
	LoggerCallback  func(ctx context.Context, msg string)
}

func BuildMiddlewarechain(config MiddlewareConfig) []Middleware[any, any] {
	var chain []Middleware[any, any]
	if config.TimerCallback != nil {
		chain = append(chain, Timer[any, any]())
	}
	if config.OutputCallback != nil {
		chain = append(chain, OutputResult[any](config.OutputCallback))
	}
	if config.MaskingCallback != nil {
		chain = append(chain, MaskOutput[any](config.MaskingCallback))
	}
	if config.RetryCount > 0 {
		chain = append(chain, Retry[any, any](config.RetryCount, config.RetryDelay))
	}
	if config.LoggerCallback != nil {
		chain = append(chain, Logging[any, any](config.LoggerCallback))
	}

	return chain
}
