// Copyright © 2020 Alexandre KOVAC <contact@kovacou.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package strava

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/kovacou/go-convert"
	"github.com/kovacou/go-types"
)

const (
	GrantAuthorizationCode = "authorization_code"
	GrantRefreshToken      = "refresh_token"
)

// AuthorizationError describe the token authorization error.
var AuthorizationError = errors.New("authorization error")

// Client is the client interface of Strava service.
type Client interface {
	// AuthorizationURL returns the URL to get the authorization from the user.
	AuthorizationURL(string) string

	// AuthorizationAccessToken get an access token from the code or refresh token.
	AuthorizationAccessToken(tok, grant string) (at AccessToken, err error)

	// Activity returns an activity
	Activity(uint64) (Activity, error)

	// Activities returns a list of activities.
	Activities(ActivitiesRequest) ([]Activity, error)

	// SetAccessToken set a new token.
	SetAccessToken(tok string)

	// SetUserID set a new default user id for user's requests.
	SetUserID(id uint64)
}

// RequestParams define the parameters to request the API.
type RequestParams struct {
	Queries            types.Map
	Values             types.Map
	WithBearer         bool
	WithFormURLEncoded bool
}

// AccessToken is the response of the Authorization.
type AccessToken struct {
	Type         string `json:"token_type"`
	ExpiresAt    int64  `json:"expires_at"`
	ExpiresIn    int64  `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
	Athlete      struct {
		ID        uint64 `json:"id"`
		Username  string `json:"username"`
		Firstname string `json:"firstname"`
		Lastname  string `json:"lastname"`
	} `json:"athlete"`
}

// close is used as defer to automatically close the body and prevent memory leak.
func closeHTTPResponse(r *http.Response) {
	if r != nil && r.Body != nil {
		_ = r.Body.Close()
	}
}

// NewEnv create a new Strava client from the environment variables.
func NewEnv() Client {
	return New(cfgEnviron)
}

// New create a new Strava client from the given config.
func New(cfg Config) Client {
	return &strava{
		cfg: cfg,
		Client: &http.Client{
			Timeout: time.Duration(cfg.Timeout) * time.Second,
		},
	}
}

// strava is the HTTP Client of the service.
type strava struct {
	*http.Client

	cfg         Config
	accessToken string
	userID      uint64
}

// SetAccessToken set a new token.
func (s *strava) SetAccessToken(tok string) {
	s.accessToken = tok
}

// SetUserID set a new user id.
func (s *strava) SetUserID(id uint64) {
	s.userID = id
}

// AuthorizationURL returns the URL to get the authorization from the user.
func (s *strava) AuthorizationURL(state string) string {
	return fmt.Sprintf("https://www.strava.com/oauth/authorize?client_id=%s&response_type=code&redirect_uri=http://runjoy.kovacou.com&approval_prompt=force&scope=%s&state=%s", s.cfg.ClientID, s.cfg.Scope, state)
}

// AuthorizationAccessToken get an access token from code or refresh_token.
func (s *strava) AuthorizationAccessToken(tok, grant string) (at AccessToken, err error) {
	p := RequestParams{
		Queries: types.Map{
			"client_id":     s.cfg.ClientID,
			"client_secret": s.cfg.ClientSecret,
		},
	}

	// Adding information based on grant type.
	switch grant {
	case GrantAuthorizationCode:
		p.Queries.Set("code", tok)

	case GrantRefreshToken:
		p.Queries.Set("refresh_token", tok)

	default:
		err = fmt.Errorf("grant_type `%s` not supported", grant)
		return
	}

	p.Queries.Set("grant_type", grant)
	resp, err := s.Request(http.MethodPost, "https://www.strava.com/api/v3/oauth/token", p)
	defer closeHTTPResponse(resp)
	if err != nil {
		return
	}

	body, _ := ioutil.ReadAll(resp.Body)
	err = json.Unmarshal(body, &at)
	return
}

// Request build a new request from the input and return the response.
func (s *strava) Request(method, uri string, p RequestParams) (r *http.Response, err error) {
	var values io.Reader

	if p.Queries == nil {
		p.Queries = types.Map{}
	}

	if p.Values == nil {
		p.Values = types.Map{}
	}

	if method == http.MethodPost {
	}

	req, err := http.NewRequest(method, uri, values)
	if err != nil {
		return nil, err
	}

	// Manage authorization to use.
	if p.WithBearer {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.accessToken))
	}

	// Managing the content type to use : some endpoint need JSON and some need form encoded.
	// To indicate Values must be encoded as FormURLEncoded, please pass WithFormURLEncoded with true.
	contentType := "application/x-www-form-urlencoded"
	if method == http.MethodPost && !p.WithFormURLEncoded {
		contentType = "application/json"
	}

	// Encoding the queries and updating the raw query.
	q := url.Values{}
	for k, val := range p.Queries {
		q.Set(k, convert.String(val))
	}
	req.URL.RawQuery = q.Encode()

	// Setting up the headers.
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Accept", "application/json;charset=UTF-8")

	r, err = s.Do(req)
	if r.StatusCode == http.StatusUnauthorized {
		err = AuthorizationError
	}
	return
}

// POST creates a new POST request.
func (s *strava) POST(endpoint string, p RequestParams) (*http.Response, error) {
	return s.Request(http.MethodPost, s.cfg.Host+endpoint, p)
}

// GET creates a new GET request.
func (s *strava) GET(endpoint string, p RequestParams) (*http.Response, error) {
	return s.Request(http.MethodGet, s.cfg.Host+endpoint, p)
}
