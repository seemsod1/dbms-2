package handlers

import (
	"encoding/json"
	"fmt"
	"laboratory_databases_2/internal/helpers"
	"laboratory_databases_2/internal/models"
	"laboratory_databases_2/internal/render"
	"net/http"
)

// Queries is the handler for the queries page
func (m *Repository) Queries(w http.ResponseWriter, r *http.Request) {
	var users []models.User
	if err := m.App.DB.Find(&users).Error; err != nil {
		helpers.ServerError(w, err)
		return
	}

	var uniqueBrands []string

	result := m.App.DB.Model(&models.Car{}).Distinct("brand").Pluck("brand", &uniqueBrands)

	if result.Error != nil {
		helpers.ServerError(w, result.Error)
	}

	var cars []models.Car
	if err := m.App.DB.Find(&cars).Error; err != nil {
		helpers.ServerError(w, err)
		return
	}

	var assignments []models.Assignment
	if err := m.App.DB.Find(&assignments).Error; err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["Users"] = users
	data["Brands"] = uniqueBrands
	data["Cars"] = cars
	data["Assignments"] = assignments

	if err := render.Template(w, r, "queries.page.tmpl", &models.TemplateData{
		Data: data,
	}); err != nil {
		m.App.ErrorLog.Println(err)
	}
}

// SimpleQuery1 returns all cars rented by a user
func (m *Repository) SimpleQuery1(w http.ResponseWriter, r *http.Request) {

	userID := r.FormValue("userID")
	var uID int
	if _, err := fmt.Sscanf(userID, "%d", &uID); err != nil {
		helpers.ServerError(w, err)
		return
	}

	var cars []models.Car
	if err := m.App.DB.Joins("JOIN rents ON cars.id = rents.car_id").
		Joins("JOIN users ON rents.user_id = users.id").
		Where("users.id = ?", userID).
		Distinct().
		Find(&cars).Error; err != nil {
		helpers.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(cars); err != nil {
		helpers.ServerError(w, err)
	}
}

// SimpleQuery2 returns all users that rented a car from a brand
func (m *Repository) SimpleQuery2(w http.ResponseWriter, r *http.Request) {

	brandName := r.FormValue("brandName")

	var users []models.User
	if err := m.App.DB.Joins("JOIN rents ON users.id = rents.user_id").
		Joins("JOIN cars ON rents.car_id = cars.id").
		Where("cars.brand = ?", brandName).
		Distinct().
		Find(&users).Error; err != nil {
		helpers.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		helpers.ServerError(w, err)
	}
}

// SimpleQuery3 returns the total amount of money spent by a user
func (m *Repository) SimpleQuery3(w http.ResponseWriter, r *http.Request) {
	userID := r.FormValue("userID")
	var uID int
	if _, err := fmt.Sscanf(userID, "%d", &uID); err != nil {
		helpers.ServerError(w, err)
		return
	}

	type TotalRent struct {
		Total float64 `json:"total"`
	}

	var totalRent TotalRent

	if err := m.App.DB.Model(&models.Rent{}).
		Joins("JOIN users ON rents.user_id = users.id").
		Where("users.id = ?", userID).
		Pluck("SUM(total)", &totalRent.Total).Error; err != nil {
		helpers.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(totalRent); err != nil {
		helpers.ServerError(w, err)
	}
}

// SimpleQuery4 returns all assignments for a car
func (m *Repository) SimpleQuery4(w http.ResponseWriter, r *http.Request) {
	carID := r.FormValue("carID")
	var cID int
	if _, err := fmt.Sscanf(carID, "%d", &cID); err != nil {
		helpers.ServerError(w, err)
		return
	}

	var assignments []models.Assignment

	query := `
        SELECT a.id, a.title
        FROM assignments a
        JOIN assignments_junction aj ON a.id = aj.assignment_id
        JOIN cars c ON c.id = aj.car_id
        WHERE c.id = ?
    `

	rows, err := m.App.DB.Raw(query, cID).Rows()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var assignment models.Assignment
		if err = rows.Scan(&assignment.ID, &assignment.Title); err != nil {
			helpers.ServerError(w, err)
			return
		}
		assignments = append(assignments, assignment)
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(assignments); err != nil {
		helpers.ServerError(w, err)
	}
}

// SimpleQuery5 returns all cars for an assignment
func (m *Repository) SimpleQuery5(w http.ResponseWriter, r *http.Request) {
	assignmentID := r.FormValue("assignmentID")
	var aID int
	if _, err := fmt.Sscanf(assignmentID, "%d", &aID); err != nil {
		helpers.ServerError(w, err)
		return
	}

	var cars []models.Car

	query := `
        SELECT c.id, c.brand, c.model, c.year
        FROM cars c
		JOIN assignments_junction aj ON c.id = aj.car_id
        JOIN assignments a ON a.id = aj.assignment_id
        WHERE a.id = ?
    `

	rows, err := m.App.DB.Raw(query, aID).Rows()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var car models.Car
		if err = rows.Scan(&car.ID, &car.Brand, &car.Model, &car.Year); err != nil {
			helpers.ServerError(w, err)
			return
		}
		cars = append(cars, car)
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(cars); err != nil {
		helpers.ServerError(w, err)
	}
}

// ComplexQuery1 returns all rents for a user and a car
func (m *Repository) ComplexQuery1(w http.ResponseWriter, r *http.Request) {
	userId := r.FormValue("userID")
	var uID int
	if _, err := fmt.Sscanf(userId, "%d", &uID); err != nil {
		helpers.ServerError(w, err)
		return
	}

	carId := r.FormValue("carID")
	var cID int
	if _, err := fmt.Sscanf(carId, "%d", &cID); err != nil {
		helpers.ServerError(w, err)
		return
	}

	var rents []models.Rent

	query := `
		SELECT r.id, r.total, r.user_id, r.car_id
		FROM rents r
		JOIN users u ON r.user_id = u.id
		JOIN cars c ON r.car_id = c.id
		WHERE u.id = ? AND c.id = ?
	`

	rows, err := m.App.DB.Raw(query, uID, cID).Rows()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var rent models.Rent
		if err = rows.Scan(&rent.ID, &rent.Total, &rent.UserID, &rent.CarID); err != nil {
			helpers.ServerError(w, err)
			return
		}
		rents = append(rents, rent)
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(rents); err != nil {
		helpers.ServerError(w, err)
	}
}

// ComplexQuery2 returns all cars rented by two users
func (m *Repository) ComplexQuery2(w http.ResponseWriter, r *http.Request) {
	firstUserID := r.FormValue("firstUserID")
	var fUID int
	if _, err := fmt.Sscanf(firstUserID, "%d", &fUID); err != nil {
		helpers.ServerError(w, err)
		return
	}

	secondUserID := r.FormValue("secondUserID")
	var sUID int
	if _, err := fmt.Sscanf(secondUserID, "%d", &sUID); err != nil {
		helpers.ServerError(w, err)
		return
	}

	if fUID == sUID {
		helpers.ServerError(w, fmt.Errorf("user's must be different"))
		return
	}

	var cars []models.Car

	query := `
		SELECT c.id, c.brand, c.model, c.year
		FROM cars c
		JOIN rents r ON c.id = r.car_id
		JOIN users u ON r.user_id = u.id
		WHERE u.id = ? AND c.id IN (
		    			SELECT c.id
		    			FROM cars c
		    			JOIN rents r ON c.id = r.car_id
		    			JOIN users u ON r.user_id = u.id
		    			WHERE u.id = ?
		)
	`

	rows, err := m.App.DB.Raw(query, fUID, sUID).Rows()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var car models.Car
		if err = rows.Scan(&car.ID, &car.Brand, &car.Model, &car.Year); err != nil {
			helpers.ServerError(w, err)
			return
		}
		cars = append(cars, car)
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(cars); err != nil {
		helpers.ServerError(w, err)
	}
}

// ComplexQuery3 returns all assignments for two cars
func (m *Repository) ComplexQuery3(w http.ResponseWriter, r *http.Request) {
	firstCarID := r.FormValue("firstCarID")
	var fCID int
	if _, err := fmt.Sscanf(firstCarID, "%d", &fCID); err != nil {
		helpers.ServerError(w, err)
		return
	}

	secondCarID := r.FormValue("secondCarID")
	var sCID int
	if _, err := fmt.Sscanf(secondCarID, "%d", &sCID); err != nil {
		helpers.ServerError(w, err)
		return
	}

	if fCID == sCID {
		helpers.ServerError(w, fmt.Errorf("cars must be different"))
		return
	}

	var assignments []models.Assignment

	query := `
		SELECT a.id, a.title
		FROM assignments a
		JOIN assignments_junction aj ON a.id = aj.assignment_id
		JOIN cars c ON c.id = aj.car_id
		WHERE c.id = ? AND a.id IN (
		    			SELECT a.id
		    			FROM assignments a
		    			JOIN assignments_junction aj ON a.id = aj.assignment_id
		    			JOIN cars c ON c.id = aj.car_id
		    			WHERE c.id = ?
		)
	`

	rows, err := m.App.DB.Raw(query, fCID, sCID).Rows()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var assignment models.Assignment
		if err = rows.Scan(&assignment.ID, &assignment.Title); err != nil {
			helpers.ServerError(w, err)
			return
		}
		assignments = append(assignments, assignment)
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(assignments); err != nil {
		helpers.ServerError(w, err)
	}
}
