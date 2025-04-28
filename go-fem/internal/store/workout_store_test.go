package store

import (
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Helper to set up a clean test DB
func setUpTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("pgx", "host=localhost user=postgres password=postgres dbname=postgres port=5433 sslmode=disable")
	if err != nil {
		t.Fatalf("opening test db: %v", err)
	}

	// Run migrations
	err = Migrate(db, "../migrations/")
	if err != nil {
		t.Fatalf("migrating test db error: %v", err)
	}

	// Clear tables
	_, err = db.Exec(`TRUNCATE workouts, workout_entries CASCADE`)
	if err != nil {
		t.Fatalf("truncating tables: %v", err)
	}

	return db
}

// Test workout creation
func TestCreateWorkout(t *testing.T) {
	db := setUpTestDB(t)
	defer db.Close()

	store := NewPostgresWorkoutStore(db) // Fixed typo here

	tests := []struct {
		name    string
		workout *Workout
		wantErr bool
	}{
		{
			name: "valid workout",
			workout: &Workout{
				Title:           "push day",
				Description:     "upper body day",
				DurationMinutes: 60,
				CaloriesBurned:  300,
				Entries: []WorkoutEntry{
					{
						ExerciseName:    "Bench Press",
						Sets:            3,
						Reps:            IntPtr(3),
						DurationSeconds: IntPtr(3),
						Weight:          FloatPtr(135.5),
						Notes:           "warm up properly",
						OrderIndex:      1,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid workout",
			workout: &Workout{
				Title:           "push day",
				Description:     "upper body day",
				DurationMinutes: 60,
				CaloriesBurned:  300,
				Entries: []WorkoutEntry{
					{
						ExerciseName: "Bench Press",
						Sets:         3,
						Reps:         IntPtr(3),
						Notes:        "warm up properly",
						OrderIndex:   1,
					},
					{
						ExerciseName:    "Bench Press",
						Sets:            3,
						Reps:            IntPtr(3),
						DurationSeconds: IntPtr(50),
						Notes:           "warm up properly",
						OrderIndex:      2,
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			createdWorkout, err := store.CreateWorkout(tt.workout)

			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.workout.Title, createdWorkout.Title)
			assert.Equal(t, tt.workout.Description, createdWorkout.Description)
			assert.Equal(t, tt.workout.DurationMinutes, createdWorkout.DurationMinutes)

			retrieved, err := store.GetWorkoutById(int64(createdWorkout.ID))
			require.NoError(t, err)
			assert.Equal(t, createdWorkout.ID, retrieved.ID)
			assert.Equal(t, len(tt.workout.Entries), len(retrieved.Entries))

			for i := range retrieved.Entries {
				assert.Equal(t, tt.workout.Entries[i].ExerciseName, retrieved.Entries[i].ExerciseName)
				assert.Equal(t, tt.workout.Entries[i].Sets, retrieved.Entries[i].Sets)
				assert.Equal(t, tt.workout.Entries[i].OrderIndex, retrieved.Entries[i].OrderIndex)
			}
		})
	}
}

// Helpers to create pointer values
func IntPtr(i int) *int           { return &i }
func FloatPtr(f float64) *float64 { return &f }
