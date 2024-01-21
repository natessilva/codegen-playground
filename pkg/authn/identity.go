package authn

import (
	"context"
	"net/http"
)

type contextKey struct{}

var key contextKey

type User struct {
	SpaceID int32
	ID      int32
}

func RequestWithUser(r *http.Request, i User) *http.Request {
	return r.WithContext(context.WithValue(r.Context(), &key, i))
}

func UserFromFromContext(ctx context.Context) User {
	return ctx.Value(&key).(User)
}
