package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/itsjayeshrathi/go-fem/internal/app"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/health", app.HealthCheck)
	r.Get("/workouts/{id}", app.WorkOutHandler.HandleGetWokroutByID)
	r.Post("/workouts", app.WorkOutHandler.HandleCreateWorkout)
	r.Put("/workouts/{id}", app.WorkOutHandler.HandlWokroutUpdateByID)
	r.Delete("/workouts/{id}", app.WorkOutHandler.HandleWorkoutDeleteById)
	return r
}
