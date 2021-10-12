package apimux

import (
	"encoding/json"
	"net/http"

	"github.com/pkg/errors"
)

func DecodeEncode(w http.ResponseWriter, r *http.Request, input interface{}, handleFunc func() (interface{}, error)) {
	if err := json.NewDecoder(r.Body).Decode(input); err != nil {
		handleError(w, errors.Wrap(err, "failed to decode"))
		return
	}
	res, err := handleFunc()
	if err != nil {
		handleError(w, err)
		return
	}
	json.NewEncoder(w).Encode(res)
}

func handleError(w http.ResponseWriter, err error) {
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(err.Error()))
}
