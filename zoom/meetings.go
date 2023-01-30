package zoom

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"github.com/eleanorhealth/go-common/pkg/errs"
)

type MeetingsServicer interface {
}

type MeetingsService struct {
	client *Client
}

var _ MeetingsServicer = (*MeetingsService)(nil)

type MeetingsListOptions struct {
	paginationOpts

	Type string `url:"type"`
}

type MeetingsListResponse struct {
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

	res, err := m.client.request(ctx, http.MethodGet, "/meetings/"+url.QueryEscape(userID), opts, nil, out)
	if err != nil {
		return nil, nil, errs.Wrap(err, "making HTTP request")
	}

	return out, res, nil
}

type MeetingsCreateOptions struct {
	DefaultPassword bool                           `json:"default_password"`
	Settings        *MeetingsCreateOptionsSettings `json:"settings"`
	StartTime       MeetingsCreateOptionsStartTime `json:"start_time"`
	Type            int                            `json:"type"`
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
	JBHTime        int  `json:"jbh_time"`
	JoinBeforeHost bool `json:"join_before_host"`
}

type MeetingsCreateResponseOccurances struct {
	Duration     int       `json:"duration"`
	OccurrenceID string    `json:"occurrence_id"`
	StartTime    time.Time `json:"start_time"`
	Status       string    `json:"status"`
}

type MeetingsCreateResponseRecurrance struct {
	EndDateTime    time.Time `json:"end_date_time"`
	EndTimes       int       `json:"end_times"`
	MonthlyDay     int       `json:"monthly_day"`
	MonthlyWeek    int       `json:"monthly_week"`
	MonthlyWeekDay int       `json:"monthly_week_day"`
	RepeatInterval int       `json:"repeat_interval"`
	Type           int       `json:"type"`
	WeeklyDays     string    `json:"weekly_days"`
}

type MeetingCreateResponseSettings struct {
	AllowMultipleDevices               bool                                                              `json:"allow_multiple_devices"`
	AlternativeHosts                   string                                                            `json:"alternative_hosts"`
	AlternativeHostsEmailNotification  bool                                                              `json:"alternative_hosts_email_notification"`
	AlternativeHostUpdatePolls         bool                                                              `json:"alternative_host_update_polls"`
	ApprovalType                       int                                                               `json:"approval_type"`
	ApprovedOrDeniedCountriesOrRegions *MeetingsCreateResponseSettingsApprovedOrDeniedCountriesOrRegions `json:"approved_or_denied_countries_or_regions"`
	Audio                              string                                                            `json:"audio"`
	AuthenticationDomains              string                                                            `json:"authentication_domains"`
	AuthenticationException            []*MeetingsCreateResponseSettingsAuthenticationException          `json:"authentication_exception"`
	AuthenticationName                 string                                                            `json:"authentication_name"`
	AuthenticationOption               string                                                            `json:"authentication_option"`
	AutoRecording                      string                                                            `json:"auto_recording"`
	BreakoutRoom                       *MeetingsCreateResponseSettingsBreakoutRoom                       `json:"breakout_room"`
	CalendarType                       int                                                               `json:"calendar_type"`
	CloseRegistration                  bool                                                              `json:"close_registration"`
	ContactEmail                       string                                                            `json:"contact_email"`
	ContactName                        string                                                            `json:"contact_name"`
	CustomKeys                         []*MeetingsCreateResponseSettingsCustomKey                        `json:"custom_keys"`
	EmailNotification                  bool                                                              `json:"email_notification"`
	EncryptionType                     string                                                            `json:"encryption_type"`
	FocusMode                          bool                                                              `json:"focus_mode"`
	GlobalDialInCountries              []string                                                          `json:"global_dial_in_countries"`
	GlobalDialInNumbers                []*MeetingsCreateResponseSettingsGlobalDialInNumber               `json:"global_dial_in_numbers"`
	HostSaveVideoOrder                 bool                                                              `json:"host_save_video_order"`
	HostVideo                          bool                                                              `json:"host_video"`
	JbhTime                            int                                                               `json:"jbh_time"`
	JoinBeforeHost                     bool                                                              `json:"join_before_host"`
	LanguageInterpretation             *MeetingsCreateResponseSettingsLanguageInterpretation             `json:"language_interpretation"`
	MeetingAuthentication              bool                                                              `json:"meeting_authentication"`
	MuteUponEntry                      bool                                                              `json:"mute_upon_entry"`
	ParticipantVideo                   bool                                                              `json:"participant_video"`
	PrivateMeeting                     bool                                                              `json:"private_meeting"`
	RegistrantsConfirmationEmail       bool                                                              `json:"registrants_confirmation_email"`
	RegistrantsEmailNotification       bool                                                              `json:"registrants_email_notification"`
	RegistrationType                   int                                                               `json:"registration_type"`
	ShowShareButton                    bool                                                              `json:"show_share_button"`
	UsePmi                             bool                                                              `json:"use_pmi"`
	WaitingRoom                        bool                                                              `json:"waiting_room"`
	Watermark                          bool                                                              `json:"watermark"`
}

type MeetingsCreateResponseSettingsApprovedOrDeniedCountriesOrRegions struct {
	ApprovedList []string `json:"approved_list"`
	DeniedList   []string `json:"denied_list"`
	Enable       bool     `json:"enable"`
	Method       string   `json:"method"`
}

type MeetingsCreateResponseSettingsAuthenticationException struct {
	Email   string `json:"email"`
	JoinURL string `json:"join_url"`
	Name    string `json:"name"`
}

type MeetingsCreateResponseSettingsBreakoutRoom struct {
	Enable bool                                              `json:"enable"`
	Rooms  []*MeetingsCreateResponseSettingsBreakoutRoomRoom `json:"rooms"`
}

type MeetingsCreateResponseSettingsBreakoutRoomRoom struct {
	Name         string   `json:"name"`
	Participants []string `json:"participants"`
}

type MeetingsCreateResponseSettingsCustomKey struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type MeetingsCreateResponseSettingsGlobalDialInNumber struct {
	City        string `json:"city"`
	Country     string `json:"country"`
	CountryName string `json:"country_name"`
	Number      string `json:"number"`
	Type        string `json:"type"`
}

type MeetingsCreateResponseSettingsLanguageInterpretation struct {
	Enable       bool                                                               `json:"enable"`
	Interpreters []*MeetingsCreateResponseSettingsLanguageInterpretationInterpreter `json:"interpreters"`
}

type MeetingsCreateResponseSettingsLanguageInterpretationInterpreter struct {
	Email     string `json:"email"`
	Languages string `json:"languages"`
}

type MeetingsCreateResponseTrackingField struct {
	Field   string `json:"field"`
	Value   string `json:"value"`
	Visible bool   `json:"visible"`
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
		return nil, res, errs.Wrap(err, "making HTTP request")
	}

	return out, res, nil
}
