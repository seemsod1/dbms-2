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
)

// Assignments is the handler for the assignments page
func (m *Repository) Assignments(w http.ResponseWriter, r *http.Request) {

	var assignments []models.Assignment

	if err := m.App.DB.Order("id").Find(&assignments).Error; err != nil {
		helpers.ServerError(w, err)
	}

	data := make(map[string]any)
	data["Assignments"] = assignments

	if err := render.Template(w, r, "assignments.page.tmpl", &models.TemplateData{
		Data: data,
	}); err != nil {
		helpers.ServerError(w, err)
	}
}

// AssignmentsUpdate updates an assignment
func (m *Repository) AssignmentsUpdate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	id := r.Form.Get("id")
	title := r.Form.Get("title")

	form := forms.New(r.PostForm)
	form.Required("title")
	form.IsNumber("id")
	if !form.Valid() {
		helpers.ServerError(w, err)
		return
	}

	err = m.App.DB.Model(&models.Assignment{}).Where("id = ?", id).Updates(models.Assignment{Title: title}).Error
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	http.Redirect(w, r, "/assignments", http.StatusOK)
}

// AssignmentsDelete deletes an assignment
func (m *Repository) AssignmentsDelete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var aID int
	_, err := fmt.Sscanf(id, "%d", &aID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	err = m.App.DB.Table("assignments_junction").Where("assignment_id = ?", aID).Delete(nil).Error
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	err = m.App.DB.Where("id = ?", aID).Delete(&models.Assignment{}).Error
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	http.Redirect(w, r, "/assignments", http.StatusOK)
}

// AssignmentsCreate creates a new assignment
func (m *Repository) AssignmentsCreate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	title := r.Form.Get("title")

	form := forms.New(r.PostForm)
	form.Required("title")
	if !form.Valid() {
		helpers.ServerError(w, err)
		return
	}
	assignment := models.Assignment{Title: title}
	err = m.App.DB.Create(&assignment).Error
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	rend.JSON(w, r, assignment)

}
