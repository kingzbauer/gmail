package pubsub

import (
	"google.golang.org/api/gmail/v1"
)

// InitPushNotification makes a watch call to gmail api for it to start triggering push notifications
func InitPushNotification(srv *gmail.Service, userID string, watchRequest *gmail.WatchRequest) (*gmail.WatchResponse, error) {
	call := srv.Users.Watch(userID, watchRequest)
	return call.Do()
}
