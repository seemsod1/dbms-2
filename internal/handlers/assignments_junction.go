package handlers

import (
	"fmt"
	rend "github.com/go-chi/render"
	"laboratory_databases_2/internal/forms"
	"laboratory_databases_2/internal/helpers"
	"laboratory_databases_2/internal/models"
	"laboratory_databases_2/internal/render"
	"net/http"
	"strconv"
)

// AssignmentsJunction is the handler for the assignments junction page
func (m *Repository) AssignmentsJunction(w http.ResponseWriter, r *http.Request) {
	var cars []models.Car
	if err := m.App.DB.Order("id").Preload("Assignments").Find(&cars).Error; err != nil {
		helpers.ServerError(w, err)
		return
	}
	var res []models.Car
	for _, car := range cars {
		if len(car.Assignments) > 0 {
			res = append(res, car)
		}
	}
	data := make(map[string]interface{})
	data["Cars"] = res

	if err := render.Template(w, r, "assignments_junction.page.tmpl", &models.TemplateData{
		Data: data,
	}); err != nil {
		helpers.ServerError(w, err)
		return
	}

}

// AssignmentsJunctionUpdate updates an assignment junction
func (m *Repository) AssignmentsJunctionUpdate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		helpers.ServerError(w, err)
		return
	}

	carID := r.Form.Get("carID")
	assignmentID := r.Form.Get("assignmentID")
	carIDold := r.Form.Get("carIDold")
	assignmentIDold := r.Form.Get("assignmentIDold")

	form := forms.New(r.PostForm)
	form.Required("carID", "assignmentID")
	form.IsNumber("carID")
	form.IsNumber("carIDold")
	form.IsNumber("assignmentID")
	form.IsNumber("assignmentIDold")
	if !form.Valid() {
		helpers.ServerError(w, nil)
		return
	}
	carIDint, err := strconv.Atoi(carID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	carIDoldint, err := strconv.Atoi(carIDold)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	assignmentIDint, err := strconv.Atoi(assignmentID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	assignmentIDoldint, err := strconv.Atoi(assignmentIDold)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	var car models.Car
	if err = m.App.DB.First(&car, carIDint).Error; err != nil {
		helpers.ServerError(w, err)
		return
	}
	var assignment models.Assignment
	if err = m.App.DB.First(&assignment, assignmentIDint).Error; err != nil {
		helpers.ServerError(w, err)
		return
	}
	var existingAssociation []models.Car
	err = m.App.DB.Joins("JOIN assignments_junction ON assignments_junction.car_id = cars.id").
		Where("assignments_junction.assignment_id = ? AND assignments_junction.car_id = ?", assignmentIDint, carIDint).
		Find(&existingAssociation).Error
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	if len(existingAssociation) > 0 {
		helpers.ServerError(w, err)
		return
	}
	err = m.App.DB.Table("assignments_junction").Where("assignment_id = ? AND car_id = ?", assignmentIDoldint, carIDoldint).Delete(nil).Error
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	err = m.App.DB.Exec("INSERT INTO assignments_junction (car_id, assignment_id) VALUES (?, ?)", carIDint, assignmentIDint).Error
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	http.Redirect(w, r, "/assignments-junction", http.StatusSeeOther)
}

// AssignmentsJunctionDelete deletes an assignment junction
func (m *Repository) AssignmentsJunctionDelete(w http.ResponseWriter, r *http.Request) {
	carID := r.URL.Query().Get("carID")
	assignmentID := r.URL.Query().Get("assignmentID")

	var cID int
	_, err := fmt.Sscanf(carID, "%d", &cID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	var aID int
	_, err = fmt.Sscanf(assignmentID, "%d", &aID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	err = m.App.DB.Table("assignments_junction").Where("assignment_id = ? AND car_id = ?", aID, cID).Delete(nil).Error
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	http.Redirect(w, r, "/assignments-junction", http.StatusSeeOther)
}

// AssignmentsJunctionCreate creates a new assignment junction
func (m *Repository) AssignmentsJunctionCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		helpers.ServerError(w, err)
		return
	}
	carID := r.Form.Get("carID")
	assignmentID := r.Form.Get("assignmentID")

	form := forms.New(r.PostForm)
	form.Required("carID", "assignmentID")
	form.IsNumber("carID")
	form.IsNumber("assignmentID")
	if !form.Valid() {
		helpers.ServerError(w, nil)
		return
	}
	carIDint, err := strconv.Atoi(carID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	assignmentIDint, err := strconv.Atoi(assignmentID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	var car models.Car
	if err = m.App.DB.First(&car, carIDint).Error; err != nil {
		helpers.ServerError(w, err)
		return
	}
	var assignment models.Assignment
	if err = m.App.DB.First(&assignment, assignmentIDint).Error; err != nil {
		helpers.ServerError(w, err)
		return
	}
	var existingAssociation []models.Car
	err = m.App.DB.Joins("JOIN assignments_junction ON assignments_junction.car_id = cars.id").
		Where("assignments_junction.assignment_id = ? AND assignments_junction.car_id = ?", assignmentIDint, carIDint).
		Find(&existingAssociation).Error
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	if len(existingAssociation) > 0 {
		helpers.ServerError(w, err)
		return
	}
	err = m.App.DB.Exec("INSERT INTO assignments_junction (car_id, assignment_id) VALUES (?, ?)", carIDint, assignmentIDint).Error
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	w.WriteHeader(http.StatusCreated)
	//render data with car and assignment
	data := make(map[string]interface{})
	data["carID"] = carID
	data["assignmentID"] = assignmentID
	rend.JSON(w, r, data)

}
