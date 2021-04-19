package access

import (
	"context"
	rsc "core/internal/pkg/resource"
	srv "core/internal/pkg/server"
	"core/pkg/dbutil"

	"github.com/lib/pq"
)

func createResource(
	sctx *srv.Ctx,
	ctx context.Context,
	i Access,
) (*Access, error) {
	var o Access

	sqlStatement := `
	INSERT INTO
		access(
		key,
		encrypted_secret,
		type,
		expires_at,
		name,
		description,
		tags
		)
	VALUES
		($1, $2, $3, $4, $5, $6, $7)
	RETURNING
		key, type, expires_at, name, description, tags;
	`
	err := dbutil.ParseError(
		rsc.Access.String(),
		i,
		sctx.DB.QueryRow(
			ctx,
			sqlStatement,
			i.Key,
			i.Secret,
			i.Type,
			i.ExpiresAt,
			i.Name,
			i.Description,
			pq.Array(i.Tags),
		).Scan(
			&o.Key,
			&o.Type,
			&o.ExpiresAt,
			&o.Name,
			&o.Description,
			&o.Tags,
		),
	)
	return &o, err
}

func getResource(
	sctx *srv.Ctx,
	ctx context.Context,
	i KeySecretPair,
) (*Access, error) {
	var o Access

	sqlStatement := `
	SELECT
		id,
		key,
		name,
		description,
		tags,
		type,
		expires_at,
		encrypted_secret
	FROM
		access
	WHERE
		key = $1
	`
	err := dbutil.ParseError(
		rsc.Access.String(),
		i,
		sctx.DB.QueryRow(
			ctx,
			sqlStatement,
			i.Key,
		).Scan(
			&o.ID,
			&o.Key,
			&o.Name,
			&o.Description,
			&o.Tags,
			&o.Type,
			&o.ExpiresAt,
			&o.Secret,
		),
	)
	return &o, err
}
