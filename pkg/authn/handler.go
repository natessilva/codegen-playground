package authn

import (
	"codegen/app/db/model"
	"database/sql"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/pkg/errors"
)

func Handle(q *model.Queries, key string, h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		a := r.Header.Get("Authorization")
		if a == "" || !strings.HasPrefix(a, "Bearer ") {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("no token provided"))
			return
		}
		a = a[7:]
		var claims userClaims
		token, err := jwt.ParseWithClaims(a, &claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(key), nil
		})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errors.Wrap(err, "failed to parse token").Error()))
			return
		}
		if !token.Valid {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("invalid token"))
			return
		}
		u, err := q.GetWorkspaceUser(r.Context(), int32(claims.WSUserID))
		if err != nil {
			if err == sql.ErrNoRows {
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("invalid token"))
				return
			}
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errors.Wrap(err, "query error").Error()))
			return
		}
		h.ServeHTTP(w, RequestWithIdentity(r, Identity{
			UserID:      int(u.UserID),
			WorkspaceID: int(u.WorkspaceID),
		}))
	}
}
