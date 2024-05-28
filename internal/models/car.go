package models

// Car represents a car entity
type Car struct {
	ID          int
	Brand       string
	Model       string
	Year        int
	Assignments []Assignment `gorm:"many2many:assignments_junction;" json:",omitempty" `
	Rents       []Rent       `gorm:"constraint:OnDelete:CASCADE;" json:",omitempty"`
}
