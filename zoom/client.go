package zoom

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/fterrag/go-zoom/zoom/tokenmutex"
	"github.com/golang-jwt/jwt/v4"
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
	tokenMutex   TokenMutex

	Users    *UsersService
	Meetings *MeetingsService
}

type PaginationOptions struct {
	NextPageToken *string `url:"next_page_token,omitempty"`
	PageSize      *int    `url:"page_size,omitempty"`
}

type PaginationResponse struct {
	NextPageToken string `json:"next_page_token"`
	PageCount     int    `json:"page_count"`
	PageSize      int    `json:"page_size"`
	TotalRecords  int    `json:"total_records"`
}

type TokenMutex interface {
	Lock(context.Context) error
	Unlock(context.Context) error
	Get(context.Context) (string, error)
	Set(context.Context, string, time.Time) error
	Clear(context.Context) error
}

// NewClient assumes the usage of Server-to-Server OAuth app (https://marketplace.zoom.us/docs/guides/build/server-to-server-oauth-app/)
func NewClient(httpClient *http.Client, accountID, clientID, clientSecret string, tokenMutex TokenMutex) *Client {
	if tokenMutex == nil {
		tokenMutex = tokenmutex.NewDefault()
	}

	c := &Client{
		httpClient:   httpClient,
		accountID:    accountID,
		clientID:     clientID,
		clientSecret: clientSecret,
		tokenMutex:   tokenMutex,
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
	err := c.tokenMutex.Lock(ctx)
	if err != nil {
		return nil, fmt.Errorf("locking token mutex: %w", err)
	}

	token, err := c.tokenMutex.Get(ctx)
	if err != nil {
		if !errors.Is(err, tokenmutex.ErrTokenNotExist) && !errors.Is(err, tokenmutex.ErrTokenExpired) {
			err = c.tokenMutex.Unlock(ctx)
			if err != nil {
				return nil, fmt.Errorf("unlocking token mutex: %w", err)
			}

			return nil, fmt.Errorf("getting token mutex: %w", err)
		}

		var expiresAt time.Time
		token, expiresAt, err = c.accessToken(ctx)
		if err != nil {
			err = c.tokenMutex.Unlock(ctx)
			if err != nil {
				return nil, fmt.Errorf("unlocking token mutex: %w", err)
			}

			return nil, fmt.Errorf("requesting access token from Zoom: %w", err)
		}

		err = c.tokenMutex.Set(context.Background(), token, expiresAt)
		if err != nil {
			err = c.tokenMutex.Unlock(ctx)
			if err != nil {
				return nil, fmt.Errorf("unlocking token mutex: %w", err)
			}

			return nil, fmt.Errorf("setting token mutex: %w", err)
		}
	}

	err = c.tokenMutex.Unlock(ctx)
	if err != nil {
		return nil, fmt.Errorf("unlocking token mutex: %w", err)
	}

	q, err := querystring.Values(query)
	if err != nil {
		return nil, fmt.Errorf("encoding URL query: %w", err)
	}

	u := baseURL + path
	if len(q) > 0 {
		u = u + "?" + q.Encode()
	}

	reader := bytes.NewReader(nil)
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshaling request body: %w", err)
		}

		reader = bytes.NewReader(b)
	}

	req, err := http.NewRequest(method, u, reader)
	if err != nil {
		return nil, fmt.Errorf("making new HTTP request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("doing HTTP request: %w", err)
	}

	if res.StatusCode > http.StatusIMUsed {
		if res.StatusCode == http.StatusUnauthorized {
			err = c.tokenMutex.Clear(ctx)
			if err != nil {
				return res, fmt.Errorf("clearing token mutex when receiving a 401 from Zoom: %w", err)
			}
		}

		errRes := &ErrorResponse{}
		err = json.NewDecoder(res.Body).Decode(errRes)
		if err != nil {
			return res, fmt.Errorf("decoding response body: %w", err)
		}

		return res, fmt.Errorf("Zoom API error: %w", errRes)
	}

	if out != nil {
		err = json.NewDecoder(res.Body).Decode(out)
		if err != nil {
			return res, fmt.Errorf("decoding response body: %w", err)
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

func (c *Client) accessToken(ctx context.Context) (string, time.Time, error) {
	query := url.Values{}
	query.Set("grant_type", "account_credentials")
	query.Set("account_id", c.accountID)

	req, err := http.NewRequestWithContext(ctx, "POST", authURL+"?"+query.Encode(), nil)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("making new HTTP request: %w", err)
	}

	auth := base64.URLEncoding.EncodeToString([]byte(c.clientID + ":" + c.clientSecret))
	req.Header.Set("Authorization", "Basic "+auth)

	res, err := c.httpClient.Do(req)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("doing HTTP request: %w", err)
	}

	if res.StatusCode != http.StatusOK {
		return "", time.Time{}, fmt.Errorf("received non-200 status code: %d", res.StatusCode)
	}

	authRes := &authResponse{}
	err = json.NewDecoder(res.Body).Decode(authRes)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("decoding HTTP response body: %w", err)
	}

	// Add a buffer to the expiration.
	expiresIn := authRes.ExpiresIn - 300

	return authRes.AccessToken, time.Now().Add(time.Duration(expiresIn) * time.Second), nil
}

// MeetingSDKJWT creates a Meeting SDK JWT, signs it, and returns the signed string (see https://marketplace.zoom.us/docs/sdk/native-sdks/auth/#meeting-sdk-auth).
// role is required for web, optional for native. 0 to specify participant or 1 to specify host.
// expiration is the duration or expiration of JWT from now. Minimum duration is 1800 seconds, maximum duration is 48 hours. Default duration is 24 hours.
func MeetingSDKJWT(meetingSDKKey, meetingSDKSecret string, meetingNumber int64, role int, expiration time.Duration) (string, error) {
	if expiration == 0 {
		expiration = 24 * time.Hour
	}

	now := time.Now().UTC()
	exp := now.Add(expiration).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"appKey":   meetingSDKKey,
		"sdkKey":   meetingSDKKey,
		"mn":       meetingNumber,
		"role":     role,
		"iat":      now.Unix(),
		"exp":      exp,
		"tokenExp": exp,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenStr, err := token.SignedString([]byte(meetingSDKSecret))
	if err != nil {
		return "", fmt.Errorf("signing JWT: %w", err)
	}

	return tokenStr, nil
}
