package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"laboratory_databases_2/internal/config"
	"laboratory_databases_2/internal/handlers"
	"net/http"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)

	mux.Get("/", handlers.Repo.Home)

	mux.Get("/users", handlers.Repo.Users)
	mux.Put("/users/update", handlers.Repo.UsersUpdate)
	mux.Delete("/users/delete/{id}", handlers.Repo.UsersDelete)
	mux.Post("/users/create", handlers.Repo.UsersCreate)

	mux.Get("/assignments", handlers.Repo.Assignments)
	mux.Put("/assignments/update", handlers.Repo.AssignmentsUpdate)
	mux.Delete("/assignments/delete/{id}", handlers.Repo.AssignmentsDelete)
	mux.Post("/assignments/create", handlers.Repo.AssignmentsCreate)

	mux.Get("/cars", handlers.Repo.Cars)
	mux.Put("/cars/update", handlers.Repo.CarsUpdate)
	mux.Delete("/cars/delete/{id}", handlers.Repo.CarsDelete)
	mux.Post("/cars/create", handlers.Repo.CarsCreate)

	mux.Get("/rents", handlers.Repo.Rents)
	mux.Put("/rents/update", handlers.Repo.RetsUpdate)
	mux.Delete("/rents/delete/{id}", handlers.Repo.RentsDelete)
	mux.Post("/rents/create", handlers.Repo.RentsCreate)

	mux.Get("/assignments-junction", handlers.Repo.AssignmentsJunction)
	mux.Put("/assignments-junction/update", handlers.Repo.AssignmentsJunctionUpdate)
	mux.Delete("/assignments-junction/delete/", handlers.Repo.AssignmentsJunctionDelete)
	mux.Post("/assignments-junction/create", handlers.Repo.AssignmentsJunctionCreate)

	mux.Get("/queries", handlers.Repo.Queries)
	mux.Post("/simple-query/1", handlers.Repo.SimpleQuery1)
	mux.Post("/simple-query/2", handlers.Repo.SimpleQuery2)
	mux.Post("/simple-query/3", handlers.Repo.SimpleQuery3)
	mux.Post("/simple-query/4", handlers.Repo.SimpleQuery4)
	mux.Post("/simple-query/5", handlers.Repo.SimpleQuery5)

	mux.Post("/complex-query/1", handlers.Repo.ComplexQuery1)
	mux.Post("/complex-query/2", handlers.Repo.ComplexQuery2)
	mux.Post("/complex-query/3", handlers.Repo.ComplexQuery3)
	return mux
}
