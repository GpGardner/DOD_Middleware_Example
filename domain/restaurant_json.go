// restaurant_json.go
package domain

type AddressJSON struct {
	Street string `json:"street"`
	City   string `json:"city"`
	State  string `json:"state"`
	Zip    string `json:"zip"`
}

type MenuItemJSON struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}

type RatingJSON struct {
	Score int    `json:"score"`
	User  string `json:"user"`
	Note  string `json:"note"`
}

type EmployeeJSON struct {
	Name string `json:"name"`
	Role string `json:"role"`
	Age  int    `json:"age"`
}

type RestaurantJSON struct {
	ID        string         `json:"id"`
	Name      string         `json:"name"`
	Email     string         `json:"email"`
	Age       int            `json:"age"`
	Address   AddressJSON    `json:"address"`
	Owners    []string       `json:"owners"`
	Employees []EmployeeJSON `json:"employees"`
	Menu      []MenuItemJSON `json:"menu"`
	Ratings   []RatingJSON   `json:"ratings"`
}
