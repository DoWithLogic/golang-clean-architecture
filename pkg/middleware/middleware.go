package middleware

import (
	"github.com/DoWithLogic/golang-clean-architecture/pkg/jwt"
	"github.com/DoWithLogic/golang-clean-architecture/pkg/types"
	"github.com/labstack/echo/v4"
)

type embedClaimedDataIntoContextOpts struct {
	claimedData *jwt.JWTClaims
}

// EmbedClaimedDataIntoContext embeds the claimed JWT data and bearer token
// into the request context. It also embeds the user login ID and action performed
// into the context, making this data available for downstream handlers.
func embedClaimedDataIntoContext(c echo.Context, opts embedClaimedDataIntoContextOpts) {
	// Store the token claims in the request context for later use
	c.Set(types.CredentialDataContextKey.String(), opts.claimedData)
}
