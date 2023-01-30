package zoom

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/eleanorhealth/go-common/pkg/errs"
)

type UsersServicer interface {
	List(ctx context.Context, opts *UsersListOptions) (*UsersListResponse, *http.Response, error)
	Create(ctx context.Context, opts *UsersCreateOptions) (*UsersCreateResponse, *http.Response, error)
	Delete(ctx context.Context, userID string, opts *UsersDeleteOptions) (*http.Response, error)
}

type UsersService struct {
	client *Client
}

var _ UsersServicer = (*UsersService)(nil)

type UsersListOptions struct {
	paginationOpts

	IncludeFields *[]string `url:"include_fields,omitempty"`
	RoleID        *string   `url:"role_id,omitempty"`
	Status        *[]string `url:"status,omitempty"`
}

type UsersListResponse struct {
	Users []*UsersListItem `json:"users"`
}

type UsersListItem struct {
	CustomAttributes  []*UsersListItemCustomAttribute `json:"custom_attributes"`
	Dept              string                          `json:"dept"`
	DisplayName       string                          `json:"display_name"`
	Email             string                          `json:"email"`
	EmployeeUniqueID  string                          `json:"employee_unique_id"`
	FirstName         string                          `json:"first_name"`
	GroupIds          []string                        `json:"group_ids"`
	ID                string                          `json:"id"`
	ImGroupIds        []string                        `json:"im_group_ids"`
	LastClientVersion string                          `json:"last_client_version"`
	LastLoginTime     time.Time                       `json:"last_login_time"`
	LastName          string                          `json:"last_name"`
	PlanUnitedType    string                          `json:"plan_united_type"`
	Pmi               int64                           `json:"pmi"`
	RoleID            string                          `json:"role_id"`
	Status            string                          `json:"status"`
	Timezone          string                          `json:"timezone"`
	Type              int                             `json:"type"`
	UserCreatedAt     time.Time                       `json:"user_created_at"`
	Verified          int                             `json:"verified"`
}

type UsersListItemCustomAttribute struct {
	Key   string `json:"key"`
	Name  string `json:"name"`
	Value string `json:"value"`
}

func (u *UsersService) List(ctx context.Context, opts *UsersListOptions) (*UsersListResponse, *http.Response, error) {
	out := &UsersListResponse{}

	res, err := u.client.request(ctx, http.MethodGet, "/users", opts, nil, out)
	if err != nil {
		return nil, res, errs.Wrap(err, "making request")
	}

	return out, res, nil
}

type UsersCreateOptions struct {
	Action   string                      `json:"action"`
	UserInfo *UsersCreateOptionsUserInfo `json:"user_info"`
}

type UsersCreateOptionsUserInfo struct {
	DisplayName *string `json:"display_name,omitempty"`
	Email       string  `json:"email"`
	FirstName   *string `json:"first_name,omitempty"`
	LastName    *string `json:"last_name,omitempty"`
	Password    *string `json:"password,omitempty"`
	Type        int     `json:"type"`
}

type UsersCreateResponse struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	ID        string `json:"id"`
	LastName  string `json:"last_name"`
	Type      int    `json:"type"`
}

func (u *UsersService) Create(ctx context.Context, opts *UsersCreateOptions) (*UsersCreateResponse, *http.Response, error) {
	out := &UsersCreateResponse{}

	res, err := u.client.request(ctx, http.MethodPost, "/users", nil, opts, out)
	if err != nil {
		return nil, res, errs.Wrap(err, "making request")
	}

	return out, res, nil
}

type UsersDeleteOptions struct {
	Action *string `url:"action,omitempty"`
}

func (u *UsersService) Delete(ctx context.Context, userID string, opts *UsersDeleteOptions) (*http.Response, error) {
	res, err := u.client.request(ctx, http.MethodDelete, "/users/"+url.QueryEscape(userID), opts, nil, nil)
	if err != nil {
		return res, errs.Wrap(err, "making request")
	}

	return res, nil
}
