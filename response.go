package bone

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// ResponseResolver gets API response and stores response body in the value pointed to by v
type ResponseResolver interface {
	CheckResponse(r *http.Response) error
	Do(ctx context.Context, req *http.Request, decode ResponseDecode, v interface{}) (response *http.Response, err error)
}

// ResponseDecode reads values from response and stores it in the value pointed to by v
type ResponseDecode = func(r io.Reader, v interface{}) error

type ResponseChecker = func(r *http.Response) error

// An ErrorResponse reports the error caused by an API request
type ErrorResponse struct {
	// HTTP response that caused this error
	Response *http.Response
}

// JSONDecode reads json value from response and stores it in the value pointed to by v
func JSONDecode(r io.Reader, v interface{}) error {
	return json.NewDecoder(r).Decode(v)
}

// CheckResponse checks the API response for errors, and returns them if present. A response is considered an
// error if it has a status code outside the 200 range.
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}
	return &ErrorResponse{Response: r}
}

func (er *ErrorResponse) Error() string {
	return fmt.Sprintf("fail to request: %v %v: %d",
		er.Response.Request.Method,
		er.Response.Request.URL,
		er.Response.StatusCode,
	)
}
