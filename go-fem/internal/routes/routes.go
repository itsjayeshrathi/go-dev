package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/itsjayeshrathi/go-fem/internal/app"
)

func SetupRoutes(app *app.Application) *chi.Mux {
	r := chi.NewRouter()

	r.Group(func(r chi.Router) {
		r.Use(app.MiddleWare.Authenticate)
		r.Get("/workouts/{id}", app.MiddleWare.RequireUser(app.WorkOutHandler.HandleGetWokroutByID))
		r.Post("/workouts", app.MiddleWare.RequireUser(app.WorkOutHandler.HandleCreateWorkout))
		r.Put("/workouts/{id}", app.MiddleWare.RequireUser(app.WorkOutHandler.HandlWokroutUpdateByID))
		r.Delete("/workouts/{id}", app.MiddleWare.RequireUser(app.WorkOutHandler.HandleWorkoutDeleteById))

	})

	r.Get("/health", app.HealthCheck)
	r.Post("/users", app.UserHandler.HandleRegisterUser)
	r.Post("/tokens/authenticate", app.TokenHandler.HandleCreateToken)

	return r
}
