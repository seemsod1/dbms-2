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
)

// Users is the handler for the users page
func (m *Repository) Users(w http.ResponseWriter, r *http.Request) {
	var users []models.User

	if err := m.App.DB.Order("id").Find(&users).Error; err != nil {
		helpers.ServerError(w, err)
		return
	}

	data := make(map[string]interface{})
	data["Users"] = users

	if err := render.Template(w, r, "users.page.tmpl", &models.TemplateData{
		Data: data,
	}); err != nil {
		helpers.ServerError(w, err)
		return
	}
}

// UsersUpdate updates a user
func (m *Repository) UsersUpdate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	userID := r.Form.Get("id")
	firstName := r.Form.Get("firstName")
	lastName := r.Form.Get("lastName")
	email := r.Form.Get("email")

	form := forms.New(r.PostForm)
	form.Required("id", "firstName", "lastName", "email")
	form.IsNumber("id")
	form.IsEmail("email")
	if !form.Valid() {
		helpers.ServerError(w, err)
		return
	}
	if err = m.App.DB.Model(&models.User{}).
		Where("id = ?", userID).Updates(models.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	},
	).Error; err != nil {
		helpers.ServerError(w, err)
		return
	}

	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

// UsersDelete deletes a user
func (m *Repository) UsersDelete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var uID uint
	_, err := fmt.Sscanf(id, "%d", &uID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	if err = m.App.DB.Where("id = ?", uID).Delete(&models.User{}).Error; err != nil {
		helpers.ServerError(w, err)
		return
	}

	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

// UsersCreate creates a new user
func (m *Repository) UsersCreate(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	firstName := r.Form.Get("firstName")
	lastName := r.Form.Get("lastName")
	email := r.Form.Get("email")

	form := forms.New(r.PostForm)
	form.Required("firstName", "lastName", "email")
	form.IsEmail("email")
	if !form.Valid() {
		helpers.ServerError(w, err)
		return
	}
	user := models.User{
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
	}
	err = m.App.DB.Create(&user).Error
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		fmt.Println("Email already exists")
		return
	}
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	rend.JSON(w, r, user)
}
