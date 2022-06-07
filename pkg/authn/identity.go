package authn

import (
	"context"
	"net/http"
)

type contextKey struct{}

var key contextKey

type Identity struct {
	WorkspaceID int
	UserID      int
}

func RequestWithIdentity(r *http.Request, i Identity) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), &key, i))
}

func IdentityFromFromContext(ctx context.Context) Identity {
	return ctx.Value(&key).(Identity)
}
