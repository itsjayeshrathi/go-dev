package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/itsjayeshrathi/go-fem/internal/store"
	"github.com/itsjayeshrathi/go-fem/internal/tokens"
	"github.com/itsjayeshrathi/go-fem/internal/utils"
)

type UserMiddlware struct {
	UserStore store.UserStore
}

type contextKey string

const UserContextKey = contextKey("user")

func SetUser(r *http.Request, user *store.User) *http.Request {
	ctx := context.WithValue(r.Context(), UserContextKey, user)
	return r.WithContext(ctx)
}

func GetUser(r *http.Request) *store.User {
	user, ok := r.Context().Value(UserContextKey).(*store.User)
	if !ok {
		panic("missing user in request")
	}
	return user
}

func (um *UserMiddlware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//withtin this annonymous function we can interject any incoming request to our server
		w.Header().Add("Vary", "Authorization")
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			r = SetUser(r, store.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}
		headerParts := strings.Split(authHeader, " ") //Bearer <TOKEN>

		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelop{"error": "invalid authorization header"})
			return
		}
		token := headerParts[1]
		user, err := um.UserStore.GetUserToken(tokens.ScopeAuth, token)
		if err != nil {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelop{"error": "invalid token"})
			return
		}
		if user == nil {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelop{"error": "token expired or invaild"})
			return
		}
		r = SetUser(r, user)
		next.ServeHTTP(w, r)
	})
}

func (um *UserMiddlware) RequireUser(next http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := GetUser(r)
		if user.IsAnonymous() {
			utils.WriteJSON(w, http.StatusUnauthorized, utils.Envelop{"error": "you must be logged in to access this route"})
			return
		}

		next.ServeHTTP(w, r)
	})
}
