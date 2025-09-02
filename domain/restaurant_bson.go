// restaurant_bson.go
package domain

type AddressBSON struct {
	Street string `bson:"street"`
	City   string `bson:"city"`
	State  string `bson:"state"`
	Zip    string `bson:"zip"`
}

type MenuItemBSON struct {
	Name        string  `bson:"name"`
	Description string  `bson:"description"`
	Price       float64 `bson:"price"`
}

type RatingBSON struct {
	Score int    `bson:"score"`
	User  string `bson:"user"`
	Note  string `bson:"note"`
}

type EmployeeBSON struct {
	Name string `bson:"name"`
	Role string `bson:"role"`
	Age  int    `bson:"age"`
}

type RestaurantBSON struct {
	ID        string         `bson:"_id,omitempty"`
	Name      string         `bson:"name"`
	Email     string         `bson:"email"`
	Age       int            `bson:"age"`
	Address   AddressBSON    `bson:"address"`
	Owners    []string       `bson:"owners"`
	Employees []EmployeeBSON `bson:"employees"`
	Menu      []MenuItemBSON `bson:"menu"`
	Ratings   []RatingBSON   `bson:"ratings"`
}
