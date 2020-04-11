package gmail

import (
	"context"
	"io/ioutil"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"golang.org/x/oauth2/jwt"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// GetConfig reads the service account json key file and returns an oauth2 configuration
func GetConfig(filepath string) (*jwt.Config, error) {
	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}

	return google.JWTConfigFromJSON(content, gmail.GmailReadonlyScope)
}

// NewService creates a new gmail service from the token source
func NewService(src oauth2.TokenSource) (*gmail.Service, error) {
	ctx := context.Background()
	return gmail.NewService(ctx, option.WithTokenSource(src))
}

// InitPushNotification makes a watch call to gmail api for it to start triggering push notifications
func InitPushNotification(srv *gmail.Service, userID string, watchRequest *gmail.WatchRequest) (*gmail.WatchResponse, error) {
	call := srv.Users.Watch(userID, watchRequest)
	return call.Do()
}

// GetUserToken performs an oauth2 verification. The user is directed to an authentication and authorization page
func GetUserToken(configPath string) (*oauth2.Token, error) {
	creds, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}
}
