package integration

import (
	"context"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/fterrag/go-zoom/zoom"
	"github.com/stretchr/testify/assert"
)

const adminUserID = "oRmhqmsvTzeDPGSFLOUhiA"

var client *zoom.Client

func TestMain(m *testing.M) {
	httpClient := &http.Client{}
	client = zoom.NewClient(httpClient, os.Getenv("ZOOM_ACCOUNT_ID"), os.Getenv("ZOOM_CLIENT_ID"), os.Getenv("ZOOM_CLIENT_SECRET"), nil)

	os.Exit(m.Run())
}

func reset(t *testing.T) {
	usersListRes, _, err := client.Users.List(context.Background(), nil)
	if err != nil {
		t.Fatal(err)
	}

	for _, user := range usersListRes.Users {
		// Admin users cannot be deleted.
		if user.ID != adminUserID {
			_, err := client.Users.Delete(context.Background(), user.ID, &zoom.UsersDeleteOptions{
				Action: zoom.Ptr("delete"),
			})
			if err != nil {
				t.Fatal(err)
			}
		}

		meetingsListRes, _, err := client.Meetings.List(context.Background(), user.ID, &zoom.MeetingsListOptions{
			Type: zoom.Ptr("scheduled"),
		})
		if err != nil {
			t.Fatal(err)
		}

		for _, meeting := range meetingsListRes.Meetings {
			_, err := client.Meetings.Delete(context.Background(), meeting.ID, &zoom.MeetingsDeleteOptions{
				ScheduleForReminder: zoom.Ptr(false),
			})
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func TestUsersCreate(t *testing.T) {
	// Creating users requires a Pro plan.
	t.Skip()

	reset(t)
}

func TestMeetingsList(t *testing.T) {
	assert := assert.New(t)

	reset(t)

	createRes, _, err := client.Meetings.Create(context.Background(), adminUserID, &zoom.MeetingsCreateOptions{
		Type: zoom.Ptr(2),
	})
	if err != nil {
		t.Fatal(err)
	}

	listRes, _, err := client.Meetings.List(context.Background(), adminUserID, nil)
	if err != nil {
		t.Fatal(err)
	}

	assert.Len(listRes.Meetings, 1)
	assert.Equal(createRes.ID, listRes.Meetings[0].ID)
}

func TestMeetingsCreate(t *testing.T) {
	assert := assert.New(t)

	reset(t)

	defaultPassword := true
	duration := 30
	jbhTime := 0
	joinBeforeHost := true
	startTime := time.Now().UTC()
	meetingType := 2

	createRes, _, err := client.Meetings.Create(context.Background(), adminUserID, &zoom.MeetingsCreateOptions{
		DefaultPassword: zoom.Ptr(defaultPassword),
		Duration:        zoom.Ptr(duration),
		Settings: &zoom.MeetingsCreateOptionsSettings{
			JBHTime:        zoom.Ptr(jbhTime),
			JoinBeforeHost: zoom.Ptr(joinBeforeHost),
		},
		StartTime: zoom.Ptr(zoom.MeetingsCreateOptionsStartTime(startTime)),
		Type:      zoom.Ptr(meetingType),
	})
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(duration, createRes.Duration)
	assert.Equal(jbhTime, createRes.Settings.JbhTime)
	assert.Equal(joinBeforeHost, createRes.Settings.JoinBeforeHost)
	assert.Equal(startTime.Truncate(time.Minute), createRes.StartTime)
	assert.Equal(meetingType, createRes.Type)
}
