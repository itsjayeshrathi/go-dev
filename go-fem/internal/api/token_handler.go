package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/itsjayeshrathi/go-fem/internal/store"
	"github.com/itsjayeshrathi/go-fem/internal/tokens"
	"github.com/itsjayeshrathi/go-fem/internal/utils"
)

type TokenHandler struct {
	tokenStore store.Tokenstore
	userStore  store.UserStore
	logger     *log.Logger
}

type createTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewTokenHandler(token store.Tokenstore, user store.UserStore, logger *log.Logger) *TokenHandler {
	return &TokenHandler{
		tokenStore: token,
		userStore:  user,
		logger:     logger,
	}
}

func (h *TokenHandler) HandleCreateToken(w http.ResponseWriter, r *http.Request) {
	var req createTokenRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		h.logger.Printf("ERROR: create request token: %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelop{"error": "invalid request payload"})
		return
	}
	//users, err
	user, err := h.userStore.GetUserByUsername(req.Username)
	if err != nil || user == nil {
		h.logger.Printf("ERROR: GetUserByUsername %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelop{"error": "internal server error"})
		return
	}

	passwordMatch, err := user.PasswordHash.Matches(req.Password)

	if err != nil {
		h.logger.Printf("ERROR: passwordHash.Matches %v", err)
		utils.WriteJSON(w, http.StatusBadRequest, utils.Envelop{"error": "invalid request payload"})
		return
	}

	if !passwordMatch {
		utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelop{"error": "invalid credantials"})
		return
	}
	token, err := h.tokenStore.CreateNewToken(user.ID, 24*time.Hour, tokens.ScopeAuth)

	if err != nil {
		h.logger.Printf("ERROR: createToken %v", err)
		utils.WriteJSON(w, http.StatusInternalServerError, utils.Envelop{"error": "internal server error"})
		return
	}
	utils.WriteJSON(w, http.StatusCreated, utils.Envelop{"auth_token": token })
}

func (h *TokenHandler) HandleDeleteToken(w http.ResponseWriter, r *http.Request) {}
