package handlers

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	rend "github.com/go-chi/render"
	"gorm.io/gorm"
	"laboratory_databases_2/internal/forms"
	"laboratory_databases_2/internal/helpers"
	"laboratory_databases_2/internal/models"
	"laboratory_databases_2/internal/render"
	"net/http"
	"strconv"
)

// Cars is the handler for the car page
func (m *Repository) Cars(w http.ResponseWriter, r *http.Request) {
	var cars []models.Car
	m.App.DB.Order("id").Find(&cars)

	data := make(map[string]interface{})
	data["Cars"] = cars

	err := render.Template(w, r, "cars.page.tmpl", &models.TemplateData{
		Data: data,
	})
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
}

// CarsUpdate updates a car
func (m *Repository) CarsUpdate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	carID := r.Form.Get("id")
	brandName := r.Form.Get("brandName")
	modelName := r.Form.Get("modelName")
	year := r.Form.Get("year")

	form := forms.New(r.PostForm)
	form.Required("id", "brandName", "modelName", "year")
	form.IsNumber("id")
	form.IsNumber("year")

	if !form.Valid() {
		helpers.ServerError(w, err)
		return
	}
	yearInt, err := strconv.Atoi(year)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	err = m.App.DB.Model(&models.Car{}).
		Where("id = ?", carID).Updates(models.Car{
		Brand: brandName,
		Model: modelName,
		Year:  yearInt,
	},
	).Error
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	http.Redirect(w, r, "/cars", http.StatusSeeOther)
}

// CarsDelete deletes a car
func (m *Repository) CarsDelete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var cID uint
	_, err := fmt.Sscanf(id, "%d", &cID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	err = m.App.DB.Table("assignments_junction").Where("car_id = ?", cID).Delete(nil).Error
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	err = m.App.DB.Where("id = ?", cID).Delete(&models.Car{}).Error
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	http.Redirect(w, r, "/cars", http.StatusSeeOther)
}

// CarsCreate creates a new car
func (m *Repository) CarsCreate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	brandName := r.Form.Get("brandName")
	modelName := r.Form.Get("modelName")
	year := r.Form.Get("year")

	form := forms.New(r.PostForm)
	form.Required("brandName", "modelName", "year")
	form.IsNumber("year")

	if !form.Valid() {
		helpers.ServerError(w, err)
		return
	}
	yearInt, err := strconv.Atoi(year)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	car := models.Car{
		Brand: brandName,
		Model: modelName,
		Year:  yearInt,
	}
	err = m.App.DB.Create(&car).Error
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		fmt.Println("Car already exists")
		return
	}
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	//write id of the car to response
	w.WriteHeader(http.StatusCreated)
	rend.JSON(w, r, car)
}
