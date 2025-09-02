package repo

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testingrepo/domain"
	"github.com/testingrepo/infra"
)

func TestNewRestaurantMiddlewareFactory(t *testing.T) {
	// Arrange
	mockRepo := &mockRestaurantReader{}

	// Act
	factory := NewRestaurantMiddlewareFactory(mockRepo)

	// Assert
	assert.NotNil(t, factory)
	assert.Equal(t, mockRepo, factory.RestaurantRepo)

	// Test FindByName initialization
	assert.NotNil(t, factory.FindRestaurantByName)

	// Test FindByName functionality
	ctx := context.Background()
	ctx = infra.EnableAll(ctx)
	output, err := factory.FindRestaurantByName(ctx, "test")
	assert.NoError(t, err)
	assert.Len(t, output.Data, 1)
	//Check test resturant name matches with OutputWithMeta data
	assert.Equal(t, "Test Restaurant", output.Data[0].Name)
	// Check that the email is masked
	assert.Equal(t, "****", output.Data[0].Email)

	// Test FindByName with a non-existent name
	output, err = factory.FindRestaurantByName(ctx, "nonexistent")
	assert.Error(t, err)
	assert.Nil(t, output.Data)
}

type mockRestaurantReader struct{}

func (m *mockRestaurantReader) FindByName(ctx context.Context, name string) ([]*domain.Restaurant, error) {
	if name == "test" {
		//RETURN DOMAIN OBJECT WITH FULL DATA
		return []*domain.Restaurant{
			{
				ID:   "1",
				Name: "Test Restaurant",
				Address: domain.Address{
					Street: "123 Test St",
					City:   "Testville",
					State:  "TS",
					Zip:    "12345",
				},
				Owners: []string{
					"Test Owner",
				},
				Ratings: []domain.Rating{
					{Score: 5, Note: "Excellent!"},
				},
				Menu: []domain.MenuItem{
					{Name: "Test Dish", Price: 9.99},
				},
				Employees: []domain.Employee{
					{Name: "John Doe", Role: "Chef", Age: 30},
				},
				Email: "TEST@TEST.COM",
			},
		}, nil
	}
	return nil, errors.New("not found")
}

func (m *mockRestaurantReader) FindByAddress(ctx context.Context, address string) ([]*domain.Restaurant, error) {
	return nil, nil
}

func (m *mockRestaurantReader) FindByOwner(ctx context.Context, owner string) ([]*domain.Restaurant, error) {
	return nil, nil
}

func (m *mockRestaurantReader) FindByRating(ctx context.Context, score int) ([]*domain.Restaurant, error) {
	return nil, nil
}

func (m *mockRestaurantReader) FindByMenuItem(ctx context.Context, itemName string) ([]*domain.Restaurant, error) {
	return nil, nil
}
