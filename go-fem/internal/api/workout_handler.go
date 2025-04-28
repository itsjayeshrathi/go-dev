package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/itsjayeshrathi/go-fem/internal/store"
	"github.com/itsjayeshrathi/go-fem/internal/utils"
)

type WorkOutHandler struct {
	workoutStore store.WorkoutStore
	logger       *log.Logger
}

func NewWorkoutHandler(workoutStore store.WorkoutStore, logger *log.Logger) *WorkOutHandler {
	return &WorkOutHandler{
		workoutStore: workoutStore,
		logger:       logger,
	}
}

func (wh *WorkOutHandler) HandleGetWokroutByID(w http.ResponseWriter, r *http.Request) {
	workoutId, err := utils.ReadIDParam(r)

	if err != nil {
		wh.logger.Printf("ERROR: readIdParam: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelop{"error": "invalid workout id"})
		return
	}

	workout, err := wh.workoutStore.GetWorkoutById(workoutId)
	if err != nil {
		wh.logger.Printf("ERROR: getWorkoutByID: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelop{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelop{"workout": workout})
}

func (wh *WorkOutHandler) HandleCreateWorkout(w http.ResponseWriter, r *http.Request) {
	var workout store.Workout
	err := json.NewDecoder(r.Body).Decode(&workout)
	if err != nil {
		wh.logger.Printf("ERROR: error decodingCreatWorkout : %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelop{"error": "internal server error"})
		return
	}
	createdWorkout, err := wh.workoutStore.CreateWorkout(&workout)
	if err != nil {
		wh.logger.Printf("ERROR: error createWokrout: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelop{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusCreated, utils.Envelop{"workout": createdWorkout})

}

func (wh *WorkOutHandler) HandlWokroutUpdateByID(w http.ResponseWriter, r *http.Request) {

	workoutId, err := utils.ReadIDParam(r)

	if err != nil {
		wh.logger.Printf("ERROR: readIdParam: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelop{"error": "invalid workout id"})
		return
	}
	existingWorkout, err := wh.workoutStore.GetWorkoutById(workoutId)

	if err != nil {
		wh.logger.Printf("ERROR: getWorkoutById: %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelop{"error": "internal server error"})
		return
	}
	if existingWorkout == nil {
		http.NotFound(w, r)
		return
	}

	// at this point we can assume we are able to find existing workout
	var updateWorkoutRequest struct {
		Title           *string              `json:"title"`
		Description     *string              `json:"description"`
		DurationMinutes *int                 `json:"duration_minutes"`
		CaloriesBurned  *int                 `json:"calories_burned"`
		Entries         []store.WorkoutEntry `json:"entries"`
	}
	err = json.NewDecoder(r.Body).Decode(&updateWorkoutRequest)

	if err != nil {
		wh.logger.Printf("ERROR: decodingUpdateRequest: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelop{"error": "internal server error"})
		return
	}

	if updateWorkoutRequest.Title != nil {
		existingWorkout.Title = *updateWorkoutRequest.Title
	}
	if updateWorkoutRequest.Description != nil {
		existingWorkout.Description = *updateWorkoutRequest.Description
	}
	if updateWorkoutRequest.DurationMinutes != nil {
		existingWorkout.DurationMinutes = *updateWorkoutRequest.DurationMinutes
	}
	if updateWorkoutRequest.CaloriesBurned != nil {
		existingWorkout.CaloriesBurned = *updateWorkoutRequest.CaloriesBurned
	}
	if updateWorkoutRequest.Entries != nil {
		existingWorkout.Entries = updateWorkoutRequest.Entries
	}
	err = wh.workoutStore.UpdateWorkout(existingWorkout)
	if err != nil {
		wh.logger.Printf("ERROR: updatingWork: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelop{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusOK, utils.Envelop{"workout": existingWorkout})
}

func (wh *WorkOutHandler) HandleWorkoutDeleteById(w http.ResponseWriter, r *http.Request) {
	paramsWorkoutId := chi.URLParam(r, "id")

	if paramsWorkoutId == "" {
		http.NotFound(w, r)
		return
	}
	workoutId, err := strconv.ParseInt(paramsWorkoutId, 10, 64)

	if err != nil {
		http.NotFound(w, r)
		return
	}
	err = wh.workoutStore.DeleteWorkoutById(workoutId)
	if err == sql.ErrNoRows {
		http.Error(w, "workout not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "error deleting workout", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)

}
