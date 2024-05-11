package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	rend "github.com/go-chi/render"
	"laboratory_databases_2/internal/forms"
	"laboratory_databases_2/internal/helpers"
	"laboratory_databases_2/internal/models"
	"laboratory_databases_2/internal/render"
	"net/http"
	"strconv"
)

// Rents is the handler for the rents page
func (m *Repository) Rents(w http.ResponseWriter, r *http.Request) {
	var rents []models.Rent
	if err := m.App.DB.Order("id").Find(&rents).Error; err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["Rents"] = rents

	err := render.Template(w, r, "rents.page.tmpl", &models.TemplateData{
		Data: data,
	})
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
}

// RetsUpdate updates a rent
func (m *Repository) RetsUpdate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		helpers.ServerError(w, err)
		return
	}
	rentID := r.Form.Get("id")
	userID := r.Form.Get("userID")
	carID := r.Form.Get("carID")
	total := r.Form.Get("total")

	form := forms.New(r.PostForm)
	form.Required("id", "userID", "carID", "total")
	form.IsNumber("id")
	form.IsNumber("userID")
	form.IsNumber("carID")
	form.IsFloat("total")

	if !form.Valid() {
		helpers.ServerError(w, nil)
		return
	}
	rentIDint, err := strconv.Atoi(rentID)
	if err != nil || rentIDint < 1 {
		helpers.ServerError(w, err)
		return
	}
	userIDint, err := strconv.Atoi(userID)
	if err != nil || userIDint < 1 {
		helpers.ServerError(w, err)
		return
	}
	carIDint, err := strconv.Atoi(carID)
	if err != nil || carIDint < 1 {
		helpers.ServerError(w, err)
		return
	}
	totalInt, err := strconv.ParseFloat(total, 64)
	if err != nil || totalInt < 1 {
		helpers.ServerError(w, err)
		return
	}
	//check if user exists
	var user models.User
	if err = m.App.DB.Where("id = ?", userIDint).First(&user).Error; err != nil {
		helpers.ServerError(w, err)
		return
	}
	//check if car exists

	var car models.Car
	if err = m.App.DB.Where("id = ?", carIDint).First(&car).Error; err != nil {
		helpers.ServerError(w, err)
		return
	}

	if err = m.App.DB.Model(&models.Rent{}).
		Where("id = ?", rentIDint).Updates(models.Rent{
		UserID: userIDint,
		CarID:  carIDint,
		Total:  totalInt,
	}).Error; err != nil {
		helpers.ServerError(w, err)
		return
	}

	http.Redirect(w, r, "/rents", http.StatusSeeOther)

}

// RentsDelete deletes a rent
func (m *Repository) RentsDelete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var rID uint
	_, err := fmt.Sscanf(id, "%d", &rID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	if err = m.App.DB.Where("id = ?", rID).Delete(&models.Rent{}).Error; err != nil {
		helpers.ServerError(w, err)
		return
	}

	http.Redirect(w, r, "/rents", http.StatusSeeOther)
}

// RentsCreate creates a new rent
func (m *Repository) RentsCreate(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		helpers.ServerError(w, err)
		return
	}
	userID := r.Form.Get("userID")
	carID := r.Form.Get("carID")
	total := r.Form.Get("total")

	form := forms.New(r.PostForm)
	form.Required("userID", "carID", "total")
	form.IsNumber("userID")
	form.IsNumber("carID")
	form.IsFloat("total")

	if !form.Valid() {
		helpers.ServerError(w, nil)
		return
	}
	userIDint, err := strconv.Atoi(userID)
	if err != nil || userIDint < 1 {
		helpers.ServerError(w, err)
		return
	}
	carIDint, err := strconv.Atoi(carID)
	if err != nil || carIDint < 1 {
		helpers.ServerError(w, err)
		return
	}
	totalInt, err := strconv.ParseFloat(total, 64)
	if err != nil || totalInt < 1 {
		helpers.ServerError(w, err)
		return
	}
	//check if user exists
	var user models.User
	if err = m.App.DB.Where("id = ?", userIDint).First(&user).Error; err != nil {
		helpers.ServerError(w, err)
		return
	}
	//check if car exists

	var car models.Car
	if err = m.App.DB.Where("id = ?", carIDint).First(&car).Error; err != nil {
		helpers.ServerError(w, err)
		return
	}

	rent := models.Rent{
		UserID: userIDint,
		CarID:  carIDint,
		Total:  totalInt,
	}

	if err = m.App.DB.Create(&rent).Error; err != nil {
		helpers.ServerError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	rend.JSON(w, r, rent)
}
