package models

// Rent represents a rent entity
type Rent struct {
	ID     int
	UserID int
	User   User
	CarID  int
	Car    Car
	Total  float64
}
