package store

import (
	"database/sql"
	"time"

	"github.com/itsjayeshrathi/go-fem/internal/tokens"
)

type PostgresTokenStore struct {
	db *sql.DB
}

func NewPostgresTokenStore(db *sql.DB) *PostgresTokenStore {
	return &PostgresTokenStore{
		db: db,
	}
}

type Tokenstore interface {
	Insert(token *tokens.Token) error
	CreateNewToken(userId int, ttl time.Duration, scope string) (*tokens.Token, error)
	DeleteAllTokensForUser(userId int, scope string) error
}

func (pg *PostgresTokenStore) Insert(token *tokens.Token) error {
	query := `INSERT INTO tokens (hash, user_id,expiry, scope) VALUES ($1, $2, $3, $4)`
	_, err := pg.db.Exec(query, token.Hash, token.UserId, token.Expiry, token.Scope)
	return err
}

func (pg *PostgresTokenStore) CreateNewToken(userId int, ttl time.Duration, scope string) (*tokens.Token, error) {
	token, err := tokens.GenerateToken(userId, ttl, scope)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (pg *PostgresTokenStore) DeleteAllTokensForUser(userId int, scope string) error {
	query := `DELETE FROM tokens WHERE scope = $1 and user_id = $2`
	_, err := pg.db.Exec(query, scope, userId)
	return err
}
