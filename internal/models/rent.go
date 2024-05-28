package models

// Rent represents a rent entity
type Rent struct {
	ID     int
	UserID int
	User   User `json:",omitempty"`
	CarID  int
	Car    Car `json:",omitempty"`
	Total  float64
}
