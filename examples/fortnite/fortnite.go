package fortnite

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/usk81/bone"
)

const (
	version        = "0.0.1"
	defaultBaseURL = "https://api.fortnitetracker.com/v1/"
	userAgent      = "go-fortnite/" + version
	tokenKey       = "TRN-Api-Key"
)

func New(httpClient *http.Client, token string) (c *bone.DefaultClient, err error) {
	c = &bone.DefaultClient{
		TokenKey: tokenKey,
		Token:    token,
	}
	if err = bone.NewClient(c, httpClient); err != nil {
		return nil, err
	}
	c.SetBaseURL(defaultBaseURL)
	c.SetUserAgent(userAgent)
	c.SetResponseChecker(CheckResponse)

	ps := &ProfileService{}
	ps.SetClient(c)
	c.SetService("profile", ps)
	return
}

// An ErrorResponse reports the error caused by an API request
type ErrorResponse struct {
	Response *http.Response
	Message  string `json:"message"`
}

func (r *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %s",
		r.Response.Request.Method,
		r.Response.Request.URL,
		r.Response.StatusCode,
		r.Message,
	)
}

// CheckResponse checks the API response for errors, and returns them if present
func CheckResponse(r *http.Response) error {
	if c := r.StatusCode; c >= 200 && c <= 299 {
		return nil
	}

	errorResponse := &ErrorResponse{Response: r}
	data, err := ioutil.ReadAll(r.Body)
	if err == nil && len(data) > 0 {
		json.Unmarshal(data, errorResponse)
	}
	return errorResponse
}
