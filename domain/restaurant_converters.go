package domain

func ConvertRestaurantToJSON(dto Restaurant) RestaurantJSON {
	return RestaurantJSON{
		ID:        dto.ID,
		Name:      dto.Name,
		Email:     dto.Email,
		Age:       dto.Age,
		Address:   AddressJSON(dto.Address),
		Owners:    dto.Owners,
		Employees: convertEmployeesToJSON(dto.Employees),
		Menu:      ConvertMenuItemsToJSON(dto.Menu),
		Ratings:   convertRatingsToJSON(dto.Ratings),
	}
}

func convertEmployeesToJSON(src []Employee) []EmployeeJSON {
	out := make([]EmployeeJSON, len(src))
	for i, e := range src {
		out[i] = EmployeeJSON(e)
	}
	return out
}

func ConvertMenuItemsToJSON(src []MenuItem) []MenuItemJSON {
	out := make([]MenuItemJSON, len(src))
	for i, m := range src {
		out[i] = MenuItemJSON(m)
	}
	return out
}

func convertRatingsToJSON(src []Rating) []RatingJSON {
	out := make([]RatingJSON, len(src))
	for i, r := range src {
		out[i] = RatingJSON(r)
	}
	return out
}

func RestaurantFromDTOToBSON(dto Restaurant) RestaurantBSON {
	return RestaurantBSON{
		ID:        dto.ID,
		Name:      dto.Name,
		Email:     dto.Email,
		Age:       dto.Age,
		Address:   AddressBSON(dto.Address),
		Owners:    dto.Owners,
		Employees: convertEmployeesToBSON(dto.Employees),
		Menu:      ConvertMenuItemsToBSON(dto.Menu),
		Ratings:   convertRatingsToBSON(dto.Ratings),
	}
}

func convertEmployeesToBSON(src []Employee) []EmployeeBSON {
	out := make([]EmployeeBSON, len(src))
	for i, e := range src {
		out[i] = EmployeeBSON(e)
	}
	return out
}

func ConvertMenuItemsToBSON(src []MenuItem) []MenuItemBSON {
	out := make([]MenuItemBSON, len(src))
	for i, m := range src {
		out[i] = MenuItemBSON(m)
	}
	return out
}

func convertRatingsToBSON(src []Rating) []RatingBSON {
	out := make([]RatingBSON, len(src))
	for i, r := range src {
		out[i] = RatingBSON(r)
	}
	return out
}

func (src *RestaurantBSON) RestaurantFromBSONToDTO() *Restaurant {
	return &Restaurant{
		ID:        src.ID,
		Name:      src.Name,
		Email:     src.Email,
		Age:       src.Age,
		Address:   Address(src.Address),
		Owners:    src.Owners,
		Employees: convertEmployeesToDTO(src.Employees),
		Menu:      ConvertMenuItemsToDTO(src.Menu),
		Ratings:   convertRatingsToDTO(src.Ratings),
	}
}

func convertEmployeesToDTO(employeeBSON []EmployeeBSON) []Employee {
	out := make([]Employee, len(employeeBSON))
	for i, e := range employeeBSON {
		out[i] = Employee(e)
	}
	return out
}

func ConvertMenuItemsToDTO(menuItemBSON []MenuItemBSON) []MenuItem {
	out := make([]MenuItem, len(menuItemBSON))
	for i, m := range menuItemBSON {
		out[i] = MenuItem(m)
	}
	return out
}

func convertRatingsToDTO(ratingBSON []RatingBSON) []Rating {
	out := make([]Rating, len(ratingBSON))
	for i, r := range ratingBSON {
		out[i] = Rating(r)
	}
	return out
}
