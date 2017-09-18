package {{.PackageName}}

import (
	"context"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
)

var StopOperation = errors.New("StopOperation")

type contextKey int

const bodyKey contextKey = iota

func parseInt(s string, base int, bitSize int) (i int64, err error) {
	return strconv.ParseInt(s, base, bitSize)
}
func getBody(r *http.Request) (bytes []byte, err error) {
	if r.Body == nil {
		return nil, nil
	}
	body := r.Context().Value(bodyKey)
	if body == nil {
		bytes, err = ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		*r = *r.WithContext(context.WithValue(r.Context(), bodyKey, bytes))
		return bytes, err
	}
	return body.([]byte), nil
}
