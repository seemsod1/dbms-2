package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"laboratory_databases_2/internal/config"
	"laboratory_databases_2/internal/handlers"
	"laboratory_databases_2/internal/models"
	"laboratory_databases_2/internal/render"
	"log"
	"os"
)

func setup(app *config.AppConfig) error {
	env, err := loadEnv()
	if err != nil {
		return err
	}

	app.Env = env

	db, err := connectDB(env)
	if err != nil {
		return err
	}

	app.DB = db

	if err = runSchemasMigration(db); err != nil {
		return err
	}

	repo := handlers.NewRepo(app)
	handlers.NewHandlers(repo)
	render.NewRenderer(app)
	return nil

}

func connectDB(env *config.EnvVariables) (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=disable",
		env.PostgresHost, env.PostgresUser, env.PostgresDBName, env.PostgresPass)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("error connecting to database: ", err)
	}

	return db, nil
}

func runSchemasMigration(db *gorm.DB) error {

	if err := db.AutoMigrate(&models.Assignment{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&models.User{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&models.Car{}); err != nil {
		return err
	}

	if err := db.AutoMigrate(&models.Rent{}); err != nil {
		return err
	}

	if err := migrateUsers(db); err != nil {
		return err
	}

	if err := migrateAssignments(db); err != nil {
		return err
	}

	if err := migrateCars(db); err != nil {
		return err
	}

	if err := migrateRents(db); err != nil {
		return err
	}

	return nil

}

func migrateUsers(db *gorm.DB) error {
	if count := db.Find(&models.User{}).RowsAffected; count > 0 {
		return nil
	}

	users := []models.User{
		{
			FirstName: "Vadym",
			LastName:  "Ripa",
			Email:     "vadym_ripa@mail.com",
		},
		{
			FirstName: "Vlad",
			LastName:  "Pavlenko",
			Email:     "vlad_pvlnk@mail.com",
		},
		{
			FirstName: "Timur",
			LastName:  "Tkachenko",
			Email:     "timur_tkachenko@mail.com",
		},
		{
			FirstName: "Dima",
			LastName:  "Nakonechniy",
			Email:     "dima_modifiks@mail.com",
		},
		{
			FirstName: "Oleksandr",
			LastName:  "Mododtsov",
			Email:     "karateodoria@mail.com",
		},
	}

	if err := db.Create(&users).Error; err != nil {
		return err
	}

	return nil
}
func migrateAssignments(db *gorm.DB) error {
	if count := db.Find(&models.Assignment{}).RowsAffected; count > 0 {
		return nil
	}

	assignments := []models.Assignment{
		{
			Title: "Family",
		},
		{
			Title: "Personal use",
		},
		{
			Title: "Weeding",
		},
		{
			Title: "Business",
		},
		{
			Title: "Sport",
		},
	}

	if err := db.Create(&assignments).Error; err != nil {
		return err
	}

	return nil
}
func migrateCars(db *gorm.DB) error {
	if count := db.Find(&models.Car{}).RowsAffected; count > 0 {
		return nil
	}

	cars := []models.Car{
		{
			Brand:       "Toyota",
			Model:       "Camry",
			Year:        2019,
			Assignments: []models.Assignment{{ID: 1}, {ID: 2}},
		},
		{
			Brand:       "BMW",
			Model:       "X5",
			Year:        2020,
			Assignments: []models.Assignment{{ID: 3}, {ID: 4}},
		},
		{
			Brand:       "Audi",
			Model:       "A6",
			Year:        2018,
			Assignments: []models.Assignment{{ID: 5}},
		},
		{
			Brand:       "Mercedes",
			Model:       "E-class",
			Year:        2017,
			Assignments: []models.Assignment{{ID: 1}, {ID: 2}},
		},
		{
			Brand:       "Ford",
			Model:       "Focus",
			Year:        2016,
			Assignments: []models.Assignment{{ID: 3}, {ID: 4}},
		},
	}

	if err := db.Create(&cars).Error; err != nil {
		return err
	}

	return nil
}
func migrateRents(db *gorm.DB) error {
	if count := db.Find(&models.Rent{}).RowsAffected; count > 0 {
		return nil
	}

	rents := []models.Rent{
		{
			UserID: 1,
			CarID:  1,
			Total:  100,
		},
		{
			UserID: 1,
			CarID:  2,
			Total:  200,
		},
		{
			UserID: 3,
			CarID:  3,
			Total:  300,
		},
		{
			UserID: 2,
			CarID:  4,
			Total:  400,
		},
		{
			UserID: 5,
			CarID:  5,
			Total:  500,
		},
		{
			UserID: 4,
			CarID:  1,
			Total:  222,
		},
		{
			UserID: 3,
			CarID:  2,
			Total:  333,
		},
		{
			UserID: 2,
			CarID:  3,
			Total:  444,
		},
		{
			UserID: 1,
			CarID:  4,
			Total:  555,
		},
		{
			UserID: 5,
			CarID:  5,
			Total:  666,
		},
	}

	if err := db.Create(&rents).Error; err != nil {
		return err
	}

	return nil
}

func loadEnv() (*config.EnvVariables, error) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	postgresHost := os.Getenv("POSTGRES_HOST")
	postgresUser := os.Getenv("POSTGRES_USER")
	postgresPass := os.Getenv("POSTGRES_PASS")
	postgresDBName := os.Getenv("POSTGRES_DBNAME")

	return &config.EnvVariables{
		PostgresHost:   postgresHost,
		PostgresUser:   postgresUser,
		PostgresPass:   postgresPass,
		PostgresDBName: postgresDBName,
	}, nil
}
