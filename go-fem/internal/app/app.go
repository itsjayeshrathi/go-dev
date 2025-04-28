package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/itsjayeshrathi/go-fem/internal/api"
	"github.com/itsjayeshrathi/go-fem/internal/middleware"
	"github.com/itsjayeshrathi/go-fem/internal/migrations"
	"github.com/itsjayeshrathi/go-fem/internal/store"
)

type Application struct {
	Logger         *log.Logger
	WorkOutHandler *api.WorkOutHandler
	UserHandler    *api.UserHandler
	TokenHandler   *api.TokenHandler
	MiddleWare     middleware.UserMiddlware
	DB             *sql.DB
}

func NewApplication() (*Application, error) {

	pgDB, err := store.Open()
	if err != nil {
		return nil, err

	}
	err = store.MigrateFS(pgDB, migrations.FS, ".")
	if err != nil {
		panic(err)
	}
	logger := log.New(os.Stdout, "", log.Ldate|log.Ltime)

	//our store will go here
	workoutStore := store.NewPostgresWorkoutStore(pgDB)
	//our handler will go here
	workoutHandler := api.NewWorkoutHandler(workoutStore, logger)

	userStore := store.NewPostgresUserStore(pgDB)

	userHandler := api.NewUserHandler(userStore, logger)

	tokenStore := store.NewPostgresTokenStore(pgDB)

	tokenHandler := api.NewTokenHandler(tokenStore, userStore, logger)

	middlewareHandler := middleware.UserMiddlware{UserStore: userStore}

	app := &Application{
		Logger:         logger,
		WorkOutHandler: workoutHandler,
		UserHandler:    userHandler,
		TokenHandler:   tokenHandler,
		MiddleWare:     middlewareHandler,
		DB:             pgDB,
	}
	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is available\n")
}
