package domain

import (
	"context"
)

//This will not exists in the actual repo package. It will be defined per application.
// Each application will Choose their own methods to be used for reading.
// EX: Dental Claims only needs to know how to read the dental claims data.

// RestaurantReader defines the interface for reading restaurants, This is a interface
// that is only using domain methods. It does not specify which database or how the data is stored.
type RestaurantReader interface {
	FindByAddress(ctx context.Context, address string) ([]*Restaurant, error)
	FindByName(ctx context.Context, name string) ([]*Restaurant, error)
	FindByOwner(ctx context.Context, owner string) ([]*Restaurant, error)
	FindByRating(ctx context.Context, score int) ([]*Restaurant, error)
	FindByMenuItem(ctx context.Context, itemName string) ([]*Restaurant, error)
}

// RestaurantWriter defines the interface for writing restaurants, This is a interface
// that is only using domain methods. It does not specify which database or how the data is stored.
type RestaurantWriter interface {
	InsertRestaurant(ctx context.Context, r *Restaurant) error
	UpdateMenu(ctx context.Context, id string, menu []MenuItem) error
	AddRating(ctx context.Context, id string, rating Rating) error
	UpdateEmployee(ctx context.Context, id string, emp Employee) error
}

// RestaurantRepository repo methods required to read and write restaurants
type RestaurantRepository interface {
	RestaurantReader
	RestaurantWriter
}
