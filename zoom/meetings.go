package zoom

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	MeetingTypeInstant        MeetingType = 1
	MeetingTypeScheduled      MeetingType = 2
	MeetingTypeRecurring      MeetingType = 3
	MeetingTypeRecurringFixed MeetingType = 8

	JoinBeforeHostAnytime   JoinBeforeHostTime = 0
	JoinBeforeHost5Minutes  JoinBeforeHostTime = 5
	JoinBeforeHost10Minutes JoinBeforeHostTime = 10
)

type MeetingType int

func (m MeetingType) Int() int {
	return int(m)
}

// JoinBeforeHostTime indicates the time limit when a participant can join a meeting before the meeting's host if join_before_host is true.
// Values are in minutes except 0 which represents anytime.
type JoinBeforeHostTime int

func (j JoinBeforeHostTime) Int() int {
	return int(j)
}

type MeetingsServicer interface {
	List(ctx context.Context, userID string, opts *MeetingsListOptions) (*MeetingsListResponse, *http.Response, error)
	Create(ctx context.Context, userID string, opts *MeetingsCreateOptions) (*MeetingsCreateResponse, *http.Response, error)
	Delete(ctx context.Context, meetingID int64, opts *MeetingsDeleteOptions) (*http.Response, error)
}

type MeetingsService struct {
	client *Client
}

var _ MeetingsServicer = (*MeetingsService)(nil)

type MeetingsListOptions struct {
	*PaginationOptions `url:",omitempty"`

	Type *string `url:"type,omitempty"`
}

type MeetingsListResponse struct {
	*PaginationResponse

	Meetings []*MeetingsListItem `json:"meetings"`
}

type MeetingsListItem struct {
	Agenda    string    `json:"agenda"`
	CreatedAt time.Time `json:"created_at"`
	Duration  int       `json:"duration"`
	HostID    string    `json:"host_id"`
	ID        int64     `json:"id"`
	JoinURL   string    `json:"join_url"`
	Pmi       string    `json:"pmi"`
	StartTime time.Time `json:"start_time"`
	Timezone  string    `json:"timezone"`
	Topic     string    `json:"topic"`
	Type      int       `json:"type"`
	UUID      string    `json:"uuid"`
}

func (m *MeetingsService) List(ctx context.Context, userID string, opts *MeetingsListOptions) (*MeetingsListResponse, *http.Response, error) {
	out := &MeetingsListResponse{}

	res, err := m.client.request(ctx, http.MethodGet, "/users/"+url.QueryEscape(userID)+"/meetings", opts, nil, out)
	if err != nil {
		return nil, res, fmt.Errorf("making HTTP request: %w", err)
	}

	return out, res, nil
}

type MeetingsCreateOptions struct {
	DefaultPassword *bool                           `json:"default_password,omitempty"`
	Duration        *int                            `json:"duration,omitempty"`
	Settings        *MeetingsCreateOptionsSettings  `json:"settings,omitempty"`
	StartTime       *MeetingsCreateOptionsStartTime `json:"start_time,omitempty"`
	Type            *int                            `json:"type,omitempty"`
}

type MeetingsCreateOptionsStartTime time.Time

func (m MeetingsCreateOptionsStartTime) MarshalJSON() ([]byte, error) {
	format := time.Time(m).UTC().Format("2006-01-02T15:04:00Z")
	b, err := json.Marshal(format)
	if err != nil {
		return nil, err
	}

	return b, nil
}

type MeetingsCreateOptionsSettings struct {
	JBHTime        *int  `json:"jbh_time,omitempty"`
	JoinBeforeHost *bool `json:"join_before_host,omitempty"`
}

type MeetingsCreateResponseOccurances struct {
	Duration     int       `json:"duration,omitempty"`
	OccurrenceID string    `json:"occurrence_id,omitempty"`
	StartTime    time.Time `json:"start_time,omitempty"`
	Status       string    `json:"status,omitempty"`
}

type MeetingsCreateResponseRecurrance struct {
	EndDateTime    time.Time `json:"end_date_time,omitempty"`
	EndTimes       int       `json:"end_times,omitempty"`
	MonthlyDay     int       `json:"monthly_day,omitempty"`
	MonthlyWeek    int       `json:"monthly_week,omitempty"`
	MonthlyWeekDay int       `json:"monthly_week_day,omitempty"`
	RepeatInterval int       `json:"repeat_interval,omitempty"`
	Type           int       `json:"type"`
	WeeklyDays     string    `json:"weekly_days,omitempty"`
}

type MeetingCreateResponseSettings struct {
	AllowMultipleDevices               bool                                                              `json:"allow_multiple_devices,omitempty"`
	AlternativeHosts                   string                                                            `json:"alternative_hosts,omitempty"`
	AlternativeHostsEmailNotification  bool                                                              `json:"alternative_hosts_email_notification,omitempty"`
	AlternativeHostUpdatePolls         bool                                                              `json:"alternative_host_update_polls,omitempty"`
	ApprovalType                       int                                                               `json:"approval_type,omitempty"`
	ApprovedOrDeniedCountriesOrRegions *MeetingsCreateResponseSettingsApprovedOrDeniedCountriesOrRegions `json:"approved_or_denied_countries_or_regions,omitempty"`
	Audio                              string                                                            `json:"audio,omitempty"`
	AuthenticationDomains              string                                                            `json:"authentication_domains,omitempty"`
	AuthenticationException            []*MeetingsCreateResponseSettingsAuthenticationException          `json:"authentication_exception,omitempty"`
	AuthenticationName                 string                                                            `json:"authentication_name,omitempty"`
	AuthenticationOption               string                                                            `json:"authentication_option,omitempty"`
	AutoRecording                      string                                                            `json:"auto_recording,omitempty"`
	BreakoutRoom                       *MeetingsCreateResponseSettingsBreakoutRoom                       `json:"breakout_room,omitempty"`
	CalendarType                       int                                                               `json:"calendar_type,omitempty"`
	CloseRegistration                  bool                                                              `json:"close_registration,omitempty"`
	ContactEmail                       string                                                            `json:"contact_email,omitempty"`
	ContactName                        string                                                            `json:"contact_name,omitempty"`
	CustomKeys                         []*MeetingsCreateResponseSettingsCustomKey                        `json:"custom_keys,omitempty"`
	EmailNotification                  bool                                                              `json:"email_notification,omitempty"`
	EncryptionType                     string                                                            `json:"encryption_type,omitempty"`
	FocusMode                          bool                                                              `json:"focus_mode,omitempty"`
	GlobalDialInCountries              []string                                                          `json:"global_dial_in_countries,omitempty"`
	GlobalDialInNumbers                []*MeetingsCreateResponseSettingsGlobalDialInNumber               `json:"global_dial_in_numbers,omitempty"`
	HostSaveVideoOrder                 bool                                                              `json:"host_save_video_order,omitempty"`
	HostVideo                          bool                                                              `json:"host_video,omitempty"`
	JbhTime                            int                                                               `json:"jbh_time,omitempty"`
	JoinBeforeHost                     bool                                                              `json:"join_before_host,omitempty"`
	LanguageInterpretation             *MeetingsCreateResponseSettingsLanguageInterpretation             `json:"language_interpretation,omitempty"`
	MeetingAuthentication              bool                                                              `json:"meeting_authentication,omitempty"`
	MuteUponEntry                      bool                                                              `json:"mute_upon_entry,omitempty"`
	ParticipantVideo                   bool                                                              `json:"participant_video,omitempty"`
	PrivateMeeting                     bool                                                              `json:"private_meeting,omitempty"`
	RegistrantsConfirmationEmail       bool                                                              `json:"registrants_confirmation_email,omitempty"`
	RegistrantsEmailNotification       bool                                                              `json:"registrants_email_notification,omitempty"`
	RegistrationType                   int                                                               `json:"registration_type,omitempty"`
	ShowShareButton                    bool                                                              `json:"show_share_button,omitempty"`
	UsePmi                             bool                                                              `json:"use_pmi,omitempty"`
	WaitingRoom                        bool                                                              `json:"waiting_room,omitempty"`
	Watermark                          bool                                                              `json:"watermark,omitempty"`
}

type MeetingsCreateResponseSettingsApprovedOrDeniedCountriesOrRegions struct {
	ApprovedList []string `json:"approved_list,omitempty"`
	DeniedList   []string `json:"denied_list,omitempty"`
	Enable       bool     `json:"enable,omitempty"`
	Method       string   `json:"method,omitempty"`
}

type MeetingsCreateResponseSettingsAuthenticationException struct {
	Email   string `json:"email,omitempty"`
	JoinURL string `json:"join_url,omitempty"`
	Name    string `json:"name,omitempty"`
}

type MeetingsCreateResponseSettingsBreakoutRoom struct {
	Enable bool                                              `json:"enable,omitempty"`
	Rooms  []*MeetingsCreateResponseSettingsBreakoutRoomRoom `json:"rooms,omitempty"`
}

type MeetingsCreateResponseSettingsBreakoutRoomRoom struct {
	Name         string   `json:"name,omitempty"`
	Participants []string `json:"participants,omitempty"`
}

type MeetingsCreateResponseSettingsCustomKey struct {
	Key   string `json:"key,omitempty"`
	Value string `json:"value,omitempty"`
}

type MeetingsCreateResponseSettingsGlobalDialInNumber struct {
	City        string `json:"city,omitempty"`
	Country     string `json:"country,omitempty"`
	CountryName string `json:"country_name,omitempty"`
	Number      string `json:"number,omitempty"`
	Type        string `json:"type,omitempty"`
}

type MeetingsCreateResponseSettingsLanguageInterpretation struct {
	Enable       bool                                                               `json:"enable,omitempty"`
	Interpreters []*MeetingsCreateResponseSettingsLanguageInterpretationInterpreter `json:"interpreters,omitempty"`
}

type MeetingsCreateResponseSettingsLanguageInterpretationInterpreter struct {
	Email     string `json:"email,omitempty"`
	Languages string `json:"languages,omitempty"`
}

type MeetingsCreateResponseTrackingField struct {
	Field   string `json:"field"`
	Value   string `json:"value,omitempty"`
	Visible bool   `json:"visible,omitempty"`
}

type MeetingsCreateResponse struct {
	Agenda          string                                 `json:"agenda"`
	AssistantID     string                                 `json:"assistant_id"`
	CreatedAt       time.Time                              `json:"created_at"`
	Duration        int                                    `json:"duration"`
	H323Password    string                                 `json:"h323_password"`
	HostEmail       string                                 `json:"host_email"`
	ID              int64                                  `json:"id"`
	JoinURL         string                                 `json:"join_url"`
	Occurrences     *MeetingsCreateResponseOccurances      `json:"occurrences"`
	Password        string                                 `json:"password"`
	Pmi             string                                 `json:"pmi"`
	PreSchedule     bool                                   `json:"pre_schedule"`
	Recurrence      *MeetingsCreateResponseRecurrance      `json:"recurrence"`
	RegistrationURL string                                 `json:"registration_url"`
	Settings        *MeetingCreateResponseSettings         `json:"settings"`
	StartTime       time.Time                              `json:"start_time"`
	StartURL        string                                 `json:"start_url"`
	Timezone        string                                 `json:"timezone"`
	Topic           string                                 `json:"topic"`
	TrackingFields  []*MeetingsCreateResponseTrackingField `json:"tracking_fields"`
	Type            int                                    `json:"type"`
}

func (m *MeetingsService) Create(ctx context.Context, userID string, opts *MeetingsCreateOptions) (*MeetingsCreateResponse, *http.Response, error) {
	out := &MeetingsCreateResponse{}

	res, err := m.client.request(ctx, http.MethodPost, "/users/"+url.QueryEscape(userID)+"/meetings", nil, opts, out)
	if err != nil {
		return nil, res, fmt.Errorf("making HTTP request: %w", err)
	}

	return out, res, nil
}

type MeetingsDeleteOptions struct {
	OccurrenceID          *string `url:"occurrence_id,omitempty"`
	ScheduleForReminder   *bool   `url:"schedule_for_reminder,omitempty"`
	CancelMeetingReminder *bool   `url:"cancel_meeting_reminder,omitempty"`
}

func (m *MeetingsService) Delete(ctx context.Context, meetingID int64, opts *MeetingsDeleteOptions) (*http.Response, error) {
	mID := strconv.Itoa(int(meetingID))
	res, err := m.client.request(ctx, http.MethodDelete, "/meetings/"+url.QueryEscape(mID), opts, nil, nil)
	if err != nil {
		return res, fmt.Errorf("making request: %w", err)
	}

	return res, nil
}
