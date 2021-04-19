package access

import (
	"context"
	cons "core/internal/pkg/constants"
	"core/internal/pkg/jwt"
	srv "core/internal/pkg/server"
	"core/pkg/crypto"
	res "core/pkg/response"
	"encoding/json"
)

// GenerateToken generate an access token via an access pair
func GenerateToken(sctx *srv.Ctx, i KeySecretPair) (
	*res.Success,
	*res.Errors,
) {
	var e res.Errors
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	r, err := getResource(
		sctx,
		ctx,
		KeySecretPair{Key: i.Key, Secret: "******"},
	)
	if err != nil {
		e.Append(cons.ErrorAuth, err.Error())
		cancel()
	}

	if err := crypto.Compare(r.Secret, i.Secret); err != nil {
		e.Append(cons.ErrorAuth, "mismatching access key-secret pair")
	}

	ma, err := json.Marshal(r)
	if err != nil {
		e.Append(cons.ErrorInternal, err.Error())
	}

	atk, err := jwt.Sign(ma)
	if err != nil {
		e.Append(cons.ErrorAuth, "unable to sign token")
	}

	return &res.Success{
		Data: &Token{
			Token:  atk,
			Access: r,
		},
	}, &e
}

// Create creates new access resource.
func Create(sctx *srv.Ctx, i Access) (
	*res.Success,
	*res.Errors,
) {
	var e res.Errors
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	encrypted_secret, err := crypto.Encrypt(i.Secret)
	if err != nil {
		e.Append(cons.ErrorCrypto, err.Error())
		cancel()
	}

	original_secret := i.Secret
	i.Secret = encrypted_secret

	r, err := createResource(
		sctx,
		ctx,
		i,
	)
	if err != nil {
		e.Append(cons.ErrorInput, err.Error())
	}

	// display un-encrypted secret one time upon creation
	r.Secret = original_secret

	return &res.Success{Data: r}, &e
}
