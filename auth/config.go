package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

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
// If the tokenFile exist, we fetch it from there
func GetUserToken(configPath string, tokenFile string) (*oauth2.Token, error) {
	// Try out an existing token file if provided
	if len(tokenFile) > 0 {
		if token, err := TokenFromFile(tokenFile); err == nil {
			return token, nil
		}
	}

	creds, err := ioutil.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	cnf, err := google.ConfigFromJSON(creds, gmail.GmailReadonlyScope)
	if err != nil {
		return nil, err
	}

	// TODO: have a function to generate random state
	url := cnf.AuthCodeURL("some-random-state", oauth2.AccessTypeOffline)
	fmt.Printf("Copy and paste this URL to your browser and then paste the code received after authorization\n%s\n", url)

	var code string
	fmt.Printf("Code: ")
	_, err = fmt.Scan(&code)
	if err != nil {
		return nil, err
	}

	// Let's exchange the code for an access token
	ctx := context.Background()
	token, err := cnf.Exchange(ctx, code, oauth2.AccessTypeOffline)
	if err != nil {
		return nil, err
	}
	return token, err
}

// TokenToFile saves the given token to file
func TokenToFile(filepath string, token *oauth2.Token) error {
	f, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	return json.NewEncoder(f).Encode(token)
}

// TokenFromFile reads a token from the provided file
func TokenFromFile(filepath string) (*oauth2.Token, error) {
	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	token := &oauth2.Token{}
	if err := json.NewDecoder(f).Decode(token); err != nil {
		return nil, err
	}

	return token, nil
}
