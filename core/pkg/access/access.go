package access

import (
	"context"
	"core/internal/constants"
	"core/internal/crypto"
	"core/internal/db"
	"core/internal/jwt"
	"core/internal/response"

	"github.com/lib/pq"
)

// Access is used to represent the relationship between the API user and the service. Access objects are attached to the resources, which are used to authorise users.
type Access struct {
	ID          string   `json:"id,omitempty"`
	Key         string   `json:"key,omitempty"`
	Description string   `json:"description,omitempty"`
	ExpiresAt   int64    `json:"expiresAt,omitempty"`
	Name        string   `json:"name,omitempty"`
	Secret      string   `json:"secret,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	Type        string   `json:"type,omitempty"`
}

// PairInput access secret-key pair, used for login
type PairInput struct {
	Key    string `json:"key,omitempty"`
	Secret string `json:"secret,omitempty"`
}

// Token access token
type Token struct {
	Token   string `json:"token,omitempty"`
	*Access `json:"access,omitempty"`
}

// GenAccessToken generate an access token via an access pair
func GenAccessToken(i PairInput) (
	*response.Success,
	*response.Errors,
) {
	var e response.Errors
	var a Access

	var encryptedSecret string
	row := db.Pool.QueryRow(context.Background(), `
	SELECT
	  id, key, type, expires_at, name, description, tags,  encrypted_secret
	FROM
	  access
	WHERE
	  key = $1
	`, i.Key)
	if err := row.Scan(
		&a.ID,
		&a.Key,
		&a.Type,
		&a.ExpiresAt,
		&a.Name,
		&a.Description,
		&a.Tags,
		&encryptedSecret,
	); err != nil {
		e.Append(constants.AuthError, "Can't find access pair")
	}

	if err := crypto.Compare(encryptedSecret, i.Secret); err != nil {
		e.Append(constants.AuthError, "Mismatching access key-secret pair")
	}

	tokenString, err := jwt.Sign(a)
	if err != nil {
		e.Append(constants.AuthError, "Unable to sign JWT")
	}

	return &response.Success{
		Data: &Token{
			Token:  tokenString,
			Access: &a,
		},
	}, &e
}

// CreateAccess create a new access key-secret pair.
func CreateAccess(i Access) (
	*response.Success,
	*response.Errors,
) {
	var e response.Errors
	var a Access
	// encrypt secret
	encryptedSecret, err := crypto.Encrypt(i.Secret)
	if err != nil {
		e.Append(constants.CryptoError, err.Error())
	}

	// create root user
	row := db.Pool.QueryRow(context.Background(), `
  INSERT INTO access
    (key, encrypted_secret, type, expires_at, name, description, tags)
  VALUES
    ($1, $2, $3, $4, $5, $6, $7)
  RETURNING
    key, type, expires_at, name, description, tags;`,
		i.Key,
		encryptedSecret,
		i.Type,
		i.ExpiresAt,
		i.Name,
		i.Description,
		pq.Array(i.Tags),
	)
	if err := row.Scan(
		&a.Key,
		&a.Type,
		&a.ExpiresAt,
		&a.Name,
		&a.Description,
		&a.Tags,
	); err != nil {
		e.Append(constants.InputError, err.Error())
	}

	// display unencrypted secret one time upon creation
	a.Secret = i.Secret

	return &response.Success{Data: a}, &e
}