package mongo

import (
	"context"
	"time"

	restaurant "github.com/testingrepo/domain"

	"go.mongodb.org/mongo-driver/bson"
)

// Collections
var (
	RESTAURANT_COLLECTION = "Restaurant"
)

type RestaurantRepo struct {
	Database *MongoClient
}

func NewRestaurantRepo() *RestaurantRepo {
	mongo, err := ConnectMongo(MongoConfig{"mongosrv://localhost:8080", "Example", 30 * time.Second})
	if err != nil {
		panic(err)
	}

	return &RestaurantRepo{
		Database: mongo,
	}
}

func (r *RestaurantRepo) FindByName(ctx context.Context, name string) ([]*restaurant.Restaurant, error) {
	filter := bson.M{"name": name}
	var docs []*restaurant.RestaurantBSON
	if err := r.Database.FindMany(ctx, RESTAURANT_COLLECTION, filter, &docs); err != nil {
		return nil, err
	}

	var restaurants = make([]*restaurant.Restaurant, len(docs))
	for i, restaurant := range docs {
		restaurants[i] = restaurant.RestaurantFromBSONToDTO()
	}

	return restaurants, nil
}

func (r *RestaurantRepo) FindByAddress(ctx context.Context, addr restaurant.Address) ([]*restaurant.Restaurant, error) {
	filter := bson.M{"address": addr}
	var docs []*restaurant.RestaurantBSON
	if err := r.Database.FindMany(ctx, RESTAURANT_COLLECTION, filter, &docs); err != nil {
		return nil, err
	}

	var restaurants = make([]*restaurant.Restaurant, len(docs))
	for i, restaurant := range docs {
		restaurants[i] = restaurant.RestaurantFromBSONToDTO()
	}

	return restaurants, nil
}

func (r *RestaurantRepo) FindByOwner(ctx context.Context, owner string) ([]*restaurant.Restaurant, error) {
	filter := bson.M{"owner": owner}
	var docs []*restaurant.RestaurantBSON
	if err := r.Database.FindMany(ctx, RESTAURANT_COLLECTION, filter, &docs); err != nil {
		return nil, err
	}

	var restaurants = make([]*restaurant.Restaurant, len(docs))
	for i, restaurant := range docs {
		restaurants[i] = restaurant.RestaurantFromBSONToDTO()
	}

	return restaurants, nil
}

func (r *RestaurantRepo) FindByRating(ctx context.Context, score int) ([]*restaurant.Restaurant, error) {
	filter := bson.M{"rating": score}
	var docs []*restaurant.RestaurantBSON
	if err := r.Database.FindMany(ctx, RESTAURANT_COLLECTION, filter, &docs); err != nil {
		return nil, err
	}

	var restaurants = make([]*restaurant.Restaurant, len(docs))
	for i, restaurant := range docs {
		restaurants[i] = restaurant.RestaurantFromBSONToDTO()
	}

	return restaurants, nil
}

func (r *RestaurantRepo) FindByMenuItem(ctx context.Context, item string) ([]*restaurant.Restaurant, error) {
	filter := bson.M{"menu": item}
	var docs []*restaurant.RestaurantBSON
	if err := r.Database.FindMany(ctx, RESTAURANT_COLLECTION, filter, &docs); err != nil {
		return nil, err
	}

	var restaurants = make([]*restaurant.Restaurant, len(docs))
	for i, restaurant := range docs {
		restaurants[i] = restaurant.RestaurantFromBSONToDTO()
	}

	return restaurants, nil
}

func (r *RestaurantRepo) InsertRestaurant(ctx context.Context, rest *restaurant.Restaurant) error {
	restBSON := restaurant.RestaurantFromDTOToBSON(*rest)
	_, err := r.Database.InsertOne(ctx, RESTAURANT_COLLECTION, restBSON)
	if err != nil {
		return err
	}
	return nil
}

func (r *RestaurantRepo) UpdateMenu(ctx context.Context, id string, menu []restaurant.MenuItem) error {
	filter := bson.M{"_id": id}
	update := restaurant.ConvertMenuItemsToBSON(menu)
	_, err := r.Database.UpdateOne(ctx, RESTAURANT_COLLECTION, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (r *RestaurantRepo) AddRating(ctx context.Context, id string, rating restaurant.Rating) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$push": bson.M{"rating": rating}}
	_, err := r.Database.UpdateOne(ctx, RESTAURANT_COLLECTION, filter, update)
	if err != nil {
		return err
	}
	return nil
}

func (r *RestaurantRepo) UpdateEmployee(ctx context.Context, id string, emp restaurant.Employee) error {
	filter := bson.M{"_id": id}
	update := bson.M{"$set": bson.M{"employee": emp}}
	_, err := r.Database.UpdateOne(ctx, RESTAURANT_COLLECTION, filter, update)
	if err != nil {
		return err
	}
	return nil
}
