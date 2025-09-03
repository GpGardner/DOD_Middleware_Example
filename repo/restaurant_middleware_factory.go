package repo

import (
	"context"
	"errors"
	"log"
	"sync"
	"time"

	"github.com/testingrepo/domain"
	"github.com/testingrepo/infra"
)

var ErrNotFound = errors.New("not found")

var (
	retries    = 2
	retryDelay = 100 * time.Millisecond
)

// ////////////////// CALLBACK FUNCTIONS ////////////////////
// func timerCallback(t time.Duration, e error) {
// 	log.Printf("TIMER: Operation took %s", t)
// }

// func retryCallback(attempt any, err error) {
// 	log.Printf("RETRY: Attempt %v failed with error: %v", attempt, err)
// }

func outputCallback(output []*domain.Restaurant, meta map[string]interface{}, err error) {
	if err != nil {
		infraLogger.Error("❌ OUTPUT: error: %+v", err)
		return
	}

	infraLogger.Info("📦 OUTPUT: Received %d restaurant(s)", len(output))
	for i, r := range output {
		infraLogger.Info("  ↪ [%d] Restaurant: %+v", i, r)
	}
}

func loggingCallback(ctx context.Context, msg string) {
	log.Printf(": %s", msg)
}

func maskingCallback(output []*domain.Restaurant) []*domain.Restaurant {
	for _, r := range output {
		// Mask email
		if r.Email != "" {
			r.Email = "****"
		}
		// Mask owner names
		for j, owner := range r.Owners {
			if len(owner) > 2 {
				r.Owners[j] = owner[:1] + "****" + owner[len(owner)-1:]
			}
		}
		// Mask employee names
		for j, emp := range r.Employees {
			if emp.Name != "" && len(emp.Name) > 2 {
				r.Employees[j].Name = emp.Name[:1] + "****" + emp.Name[len(emp.Name)-1:]
			}
		}
	}
	return output
}

//////////////////////////////////////////////////////////

type RestaurantMiddlewareFactory struct {
	RestaurantRepo domain.RestaurantReader
	once           sync.Once

	FindRestaurantByName     infra.RepoOp[string, []*domain.Restaurant]
	FindRestaurantByAddress  infra.RepoOp[string, []*domain.Restaurant]
	FindRestaurantByOwner    infra.RepoOp[string, []*domain.Restaurant]
	FindRestaurantByRating   infra.RepoOp[int, []*domain.Restaurant]
	FindRestaurantByMenuItem infra.RepoOp[string, []*domain.Restaurant]
}

func NewRestaurantMiddlewareFactory(repo domain.RestaurantReader) *RestaurantMiddlewareFactory {
	f := &RestaurantMiddlewareFactory{
		RestaurantRepo: repo,
		once:           sync.Once{},
	}

	f.initFindByName()
	f.initFindByAddress()
	f.initFindByOwner()
	f.initFindByRating()
	f.initFindByMenuItem()
	return f
}

func (f *RestaurantMiddlewareFactory) GetRestaurantReader() domain.RestaurantReader {
	return f.RestaurantRepo
}

type logger struct{}

func (l logger) Info(msg string, fields ...interface{}) {
	log.Printf("INFO: "+msg, fields...)
}
func (l logger) Error(msg string, fields ...interface{}) {
	log.Printf("ERROR: "+msg, fields...)
}

var infraLogger logger = logger{}

func (f *RestaurantMiddlewareFactory) initFindByName() {
	f.once.Do(func() {
		builder := infra.MiddlewareBuilder[string, []*domain.Restaurant]{}

		// Add middleware dynamically
		builder.Add(infra.Gate(infra.Logging[string, []*domain.Restaurant](loggingCallback), infra.IsTimingDisabled))
		builder.Add(infra.Gate(infra.Timer[string, []*domain.Restaurant](), infra.IsTimingDisabled))
		builder.Add(infra.Gate(infra.OutputResult[string](outputCallback), infra.IsOutputResultDisabled))
		builder.Add(infra.Gate(infra.MaskOutput[string](maskingCallback), infra.IsMaskingDisabled))
		builder.Add(infra.Gate(infra.Retry[string, []*domain.Restaurant](retries, retryDelay), infra.IsRetryDisabled))

		// Build the chain
		f.FindRestaurantByName = builder.Build(f.bindFindByName())
	})
}

// binder: repo method -> RepoOp
func (f *RestaurantMiddlewareFactory) bindFindByName() infra.RepoOp[string, []*domain.Restaurant] {
	return func(ctx context.Context, name string) (infra.OutputWithMeta[[]*domain.Restaurant], error) {
		data, err := f.RestaurantRepo.FindByName(ctx, name)
		return infra.OutputWithMeta[[]*domain.Restaurant]{Data: data}, err
	}
}

// FindByAddress(ctx context.Context, address Address) ([]*Restaurant, error)
func (f *RestaurantMiddlewareFactory) initFindByAddress() {
	f.once.Do(func() {
		builder := infra.MiddlewareBuilder[string, []*domain.Restaurant]{}

		// Add middleware dynamically
		builder.Add(infra.Gate(infra.Logging[string, []*domain.Restaurant](loggingCallback), infra.IsTimingDisabled))
		builder.Add(infra.Gate(infra.Timer[string, []*domain.Restaurant](), infra.IsTimingDisabled))
		builder.Add(infra.Gate(infra.OutputResult[string](outputCallback), infra.IsOutputResultDisabled))
		builder.Add(infra.Gate(infra.MaskOutput[string](maskingCallback), infra.IsMaskingDisabled))
		builder.Add(infra.Gate(infra.Retry[string, []*domain.Restaurant](retries, retryDelay), infra.IsRetryDisabled))

		// Build the chain
		f.FindRestaurantByAddress = builder.Build(f.bindFindByAddress())
	})
}

func (f *RestaurantMiddlewareFactory) bindFindByAddress() infra.RepoOp[string, []*domain.Restaurant] {
	return func(ctx context.Context, addr string) (infra.OutputWithMeta[[]*domain.Restaurant], error) {
		data, err := f.RestaurantRepo.FindByAddress(ctx, addr)
		return infra.OutputWithMeta[[]*domain.Restaurant]{Data: data}, err
	}
}

// FindByOwner(ctx context.Context, owner string) ([]*Restaurant, error)
func (f *RestaurantMiddlewareFactory) initFindByOwner() {
	f.once.Do(func() {
		builder := infra.MiddlewareBuilder[string, []*domain.Restaurant]{}
		// Add middleware dynamically
		builder.Add(infra.Gate(infra.Logging[string, []*domain.Restaurant](loggingCallback), infra.IsTimingDisabled))
		builder.Add(infra.Gate(infra.Timer[string, []*domain.Restaurant](), infra.IsTimingDisabled))
		builder.Add(infra.Gate(infra.OutputResult[string](outputCallback), infra.IsOutputResultDisabled))
		builder.Add(infra.Gate(infra.MaskOutput[string](maskingCallback), infra.IsMaskingDisabled))
		builder.Add(infra.Gate(infra.Retry[string, []*domain.Restaurant](retries, retryDelay), infra.IsRetryDisabled))
		// Build the chain
		f.FindRestaurantByOwner = builder.Build(f.bindFindByOwner())
	})
}

func (f *RestaurantMiddlewareFactory) bindFindByOwner() infra.RepoOp[string, []*domain.Restaurant] {
	return func(ctx context.Context, owner string) (infra.OutputWithMeta[[]*domain.Restaurant], error) {
		data, err := f.RestaurantRepo.FindByOwner(ctx, owner)
		return infra.OutputWithMeta[[]*domain.Restaurant]{Data: data}, err
	}
}

// FindByRating(ctx context.Context, score int) ([]*Restaurant, error)
func (f *RestaurantMiddlewareFactory) initFindByRating() {
	f.once.Do(func() {
		builder := infra.MiddlewareBuilder[int, []*domain.Restaurant]{}
		// Add middleware dynamically
		builder.Add(infra.Gate(infra.Logging[int, []*domain.Restaurant](loggingCallback), infra.IsTimingDisabled))
		builder.Add(infra.Gate(infra.Timer[int, []*domain.Restaurant](), infra.IsTimingDisabled))
		builder.Add(infra.Gate(infra.OutputResult[int](outputCallback), infra.IsOutputResultDisabled))
		builder.Add(infra.Gate(infra.MaskOutput[int](maskingCallback), infra.IsMaskingDisabled))
		builder.Add(infra.Gate(infra.Retry[int, []*domain.Restaurant](retries, retryDelay), infra.IsRetryDisabled))
		// Build the chain
		f.FindRestaurantByRating = builder.Build(f.bindFindByRating())
	})
}

func (f *RestaurantMiddlewareFactory) bindFindByRating() infra.RepoOp[int, []*domain.Restaurant] {
	return func(ctx context.Context, score int) (infra.OutputWithMeta[[]*domain.Restaurant], error) {
		data, err := f.RestaurantRepo.FindByRating(ctx, score)
		return infra.OutputWithMeta[[]*domain.Restaurant]{Data: data}, err
	}
}

// FindByMenuItem(ctx context.Context, itemName string) ([]*Restaurant, error)
func (f *RestaurantMiddlewareFactory) initFindByMenuItem() {
	f.once.Do(func() {
		builder := infra.MiddlewareBuilder[string, []*domain.Restaurant]{}
		// Add middleware dynamically
		builder.Add(infra.Gate(infra.Logging[string, []*domain.Restaurant](loggingCallback), infra.IsTimingDisabled))
		builder.Add(infra.Gate(infra.Timer[string, []*domain.Restaurant](), infra.IsTimingDisabled))
		builder.Add(infra.Gate(infra.OutputResult[string](outputCallback), infra.IsOutputResultDisabled))
		builder.Add(infra.Gate(infra.MaskOutput[string](maskingCallback), infra.IsMaskingDisabled))
		builder.Add(infra.Gate(infra.Retry[string, []*domain.Restaurant](retries, retryDelay), infra.IsRetryDisabled))
		// Build the chain
		f.FindRestaurantByMenuItem = builder.Build(f.bindFindByMenuItem())
	})
}

func (f *RestaurantMiddlewareFactory) bindFindByMenuItem() infra.RepoOp[string, []*domain.Restaurant] {
	return func(ctx context.Context, item string) (infra.OutputWithMeta[[]*domain.Restaurant], error) {
		data, err := f.RestaurantRepo.FindByMenuItem(ctx, item)
		return infra.OutputWithMeta[[]*domain.Restaurant]{Data: data}, err
	}
}
