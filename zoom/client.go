package zoom

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/eleanorhealth/go-common/pkg/errs"
	querystring "github.com/google/go-querystring/query"
)

const (
	authURL = "https://zoom.us/oauth/token"
	baseURL = "https://api.zoom.us/v2"
)

type Client struct {
	httpClient   *http.Client
	accountID    string
	clientID     string
	clientSecret string

	tokenLock       sync.Mutex
	token           string
	tokenExpiration time.Time

	Users    *UsersService
	Meetings *MeetingsService
}

type paginationOpts struct {
	NextPageToken string `url:"next_page_token"`
	PageNumber    int    `url:"page_number"`
	PageSize      int    `url:"page_size"`
}

// NewClient assumes the usage of Server-to-Server OAuth app (https://marketplace.zoom.us/docs/guides/build/server-to-server-oauth-app/)
func NewClient(httpClient *http.Client, accountID, clientID, clientSecret string) *Client {
	c := &Client{
		httpClient:   httpClient,
		accountID:    accountID,
		clientID:     clientID,
		clientSecret: clientSecret,
	}

	c.Users = &UsersService{c}
	c.Meetings = &MeetingsService{c}

	return c
}

type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Errors  []FieldError
}

func (e *ErrorResponse) Error() string {
	return e.Message
}

type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (c *Client) request(ctx context.Context, method string, path string, query any, body any, out any) (*http.Response, error) {
	token, err := c.accessToken()
	if err != nil {
		return nil, errs.Wrap(err, "getting access token")
	}

	q, err := querystring.Values(query)
	if err != nil {
		return nil, errs.Wrap(err, "encoding URL query")
	}

	u := baseURL + path
	if len(q) > 0 {
		u = u + "?" + q.Encode()
	}

	reader := bytes.NewReader(nil)
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, errs.Wrap(err, "marshaling request body")
		}

		reader = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, u, reader)
	if err != nil {
		return nil, errs.Wrap(err, "making new HTTP request")
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, errs.Wrap(err, "doing HTTP request")
	}

	if res.StatusCode > http.StatusIMUsed {
		errRes := &ErrorResponse{}
		err = json.NewDecoder(res.Body).Decode(errRes)
		if err != nil {
			return res, errs.Wrap(err, "decoding response body")
		}

		return res, errs.Wrap(errRes, "Zoom API error")
	}

	if out != nil {
		err = json.NewDecoder(res.Body).Decode(out)
		if err != nil {
			return res, errs.Wrap(err, "decoding response body")
		}
	}

	return res, nil
}

type authResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int64  `json:"expires_in"`
	Scope       string `json:"scope"`
}

func (c *Client) accessToken() (string, error) {
	c.tokenLock.Lock()
	defer c.tokenLock.Unlock()

	if c.tokenExpiration.After(time.Now()) {
		return c.token, nil
	}

	query := url.Values{}
	query.Set("grant_type", "account_credentials")
	query.Set("account_id", c.accountID)

	req, err := http.NewRequest("POST", authURL+"?"+query.Encode(), nil)
	if err != nil {
		return "", errs.Wrap(err, "making new HTTP request")
	}

	auth := base64.URLEncoding.EncodeToString([]byte(c.clientID + ":" + c.clientSecret))
	req.Header.Set("Authorization", "Basic "+auth)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return "", errs.Wrap(err, "doing HTTP request")
	}

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("received non-200 status code: %d", res.StatusCode)
	}

	authRes := &authResponse{}
	err = json.NewDecoder(res.Body).Decode(authRes)
	if err != nil {
		return "", errs.Wrap(err, "decoding HTTP response body")
	}

	// Add a 60 second buffer to the expiration.
	expiresIn := authRes.ExpiresIn - 60

	c.token = authRes.AccessToken
	c.tokenExpiration = time.Now().Add(time.Duration(expiresIn) * time.Second)

	return c.token, nil
}
