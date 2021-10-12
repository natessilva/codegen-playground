package authn

import (
	"context"
	"net/http"
)

type contextKey struct{}

var key contextKey

type Identity struct {
	UserID int
}

func requestWithIdentity(r *http.Request, a Identity) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), &key, a))
}

func IdentityFromFromContext(ctx context.Context) Identity {
	return ctx.Value(&key).(Identity)
}
