package api

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"regexp"

	"github.com/itsjayeshrathi/go-fem/internal/store"
	"github.com/itsjayeshrathi/go-fem/internal/utils"
)

type registeredUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}

type UserHandler struct {
	userStore store.UserStore
	logger    *log.Logger
}

func NewUserHandler(userStore store.UserStore, logger *log.Logger) *UserHandler {
	return &UserHandler{
		userStore: userStore,
		logger:    logger,
	}
}

func (h *UserHandler) validateRegisterRequest(req *registeredUserRequest) error {
	if req.Username == "" {
		return errors.New("username can't be empty")
	}
	if len(req.Username) > 50 {
		return errors.New("username cannot be greater than 50")
	}

	if req.Email == "" {
		return errors.New("email can't be empty")
	}

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(req.Email) {
		return errors.New("invvalid email format")
	}

	if req.Password == "" {
		return errors.New("password can't be empty")
	}

	return nil

}

func (h *UserHandler) HandleRegisterUser(w http.ResponseWriter, r *http.Request) {

	var req registeredUserRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Printf("ERROR: decoding register user: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelop{"error": "invalid request payload"})
		return
	}

	err = h.validateRegisterRequest(&req)
	if err != nil {
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelop{"error": err.Error()})
		return
	}

	user := &store.User{
		Username: req.Username,
		Email:    req.Email,
	}
	if req.Bio != "" {
		user.Bio = req.Bio
	}
	// how do we deal with their password

	err = user.PasswordHash.Set(req.Password)
	if err != nil {
		h.logger.Printf("ERROR: error in hashing %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelop{"error": "internal server error"})
		return
	}
	err = h.userStore.CreateUser(user)
	if err != nil {
		h.logger.Printf("ERROR: error in hashing %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelop{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusCreated, utils.Envelop{"user": user})

}
func (h *UserHandler) HandleGetUserByUsername(w http.ResponseWriter, r *http.Request) {

}

func (h *UserHandler) HandleUpdateUser(w http.ResponseWriter, r *http.Request) {}
