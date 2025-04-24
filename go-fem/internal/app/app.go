package app

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/itsjayeshrathi/go-fem/internal/api"
	"github.com/itsjayeshrathi/go-fem/internal/migrations"
	"github.com/itsjayeshrathi/go-fem/internal/store"
)

type Application struct {
	Logger         *log.Logger
	WorkOutHandler *api.WorkOutHandler
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
	workoutStore := store.NewPostgresWokroutStore(pgDB)
	//our handler will go here
	workoutHandler := api.NewWorkoutHandler(workoutStore, logger)

	app := &Application{
		Logger:         logger,
		WorkOutHandler: workoutHandler,
		DB:             pgDB,
	}
	return app, nil
}

func (a *Application) HealthCheck(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Status is available\n")
}
