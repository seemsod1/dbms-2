package models

// User represents a user entity
type User struct {
	ID        int
	FirstName string
	LastName  string
	Email     string `gorm:"unique"`
	Rents     []Rent `gorm:"constraint:OnDelete:CASCADE;"`
}
