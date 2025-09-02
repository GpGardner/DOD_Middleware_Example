// Package domain defines the data structures and interfaces for the restaurant domain.
package domain

type Address struct {
	Street string
	City   string
	State  string
	Zip    string
}

type MenuItem struct {
	Name        string
	Description string
	Price       float64
}

type Rating struct {
	Score int
	User  string
	Note  string
}

type Employee struct {
	Name string
	Role string
	Age  int
}

type Restaurant struct {
	ID        string
	Name      string
	Email     string
	Age       int
	Address   Address
	Owners    []string
	Employees []Employee
	Menu      []MenuItem
	Ratings   []Rating
}
