package strava

import (
	"encoding/json"
	"log/slog"
	"os"

	"github.com/shtamura/oauth2"
)

// Strava Authorization settings
// https://developers.strava.com/docs/authentication/
type AuthConfig struct {
	ClientId     string
	ClientSecret string
}

func NewConfig() *AuthConfig {
	return &AuthConfig{
		ClientId:     os.Getenv("STRAVA_CLIENT_ID"),
		ClientSecret: os.Getenv("STRAVA_CLIENT_SECRET"),
	}
}

// Authorize to Strava
func Authorize(config *AuthConfig) error {
	res, err := oauth2.Authorize(
		"https://www.strava.com/oauth/authorize",
		"https://www.strava.com/oauth/token",
		config.ClientId,
		config.ClientSecret,
		nil)
	if err != nil {
		return err
	}
	return saveCredential(res)
}

func RefreshToken(config *AuthConfig) error {
	credential, err := GetCredential()
	if err != nil {
		return err
	}

	res, err := oauth2.RefreshToken(
		"https://www.strava.com/oauth/token",
		config.ClientId,
		config.ClientSecret,
		credential.RefreshToken,
	)
	if err != nil {
		return err
	}
	return saveCredential(res)
}

// save credential to home directory(.strava-cli/credential.json)
func saveCredential(credential *oauth2.AuthorizationResponse) error {
	// save
	home, err := os.UserHomeDir()
	if err != nil {
		slog.Error("Failed to get home directory: %v", err)
		return err
	}
	// save to .strava-cli/credential.json
	jsonCred, err := json.Marshal(credential)
	if err != nil {
		slog.Error("Failed to marshal credential: %v", err)
		return err
	}
	err = os.MkdirAll(home+"/.strava-cli", 0755)
	if err != nil {
		slog.Error("Failed to create directory: %v", err)
		return err
	}
	err = os.WriteFile(home+"/.strava-cli/credential.json", jsonCred, 0644)
	if err != nil {
		slog.Error("Failed to write credential: %v", err)
		return err
	}
	return nil
}

// get credential from home directory(.strava-cli/credential.json)
func GetCredential() (*oauth2.AuthorizationResponse, error) {
	// get
	home, err := os.UserHomeDir()
	if err != nil {
		slog.Error("Failed to get home directory: %v", err)
		return nil, err
	}
	// get from .strava-cli/credential.json
	jsonCred, err := os.ReadFile(home + "/.strava-cli/credential.json")
	if os.IsNotExist(err) {
		slog.Warn("Credential not found")
		return nil, nil
	}
	if err != nil {
		slog.Error("Failed to read credential: %v", err)
		return nil, err
	}
	credential := &oauth2.AuthorizationResponse{}
	err = json.Unmarshal(jsonCred, credential)
	if err != nil {
		slog.Error("Failed to unmarshal credential: %v", err)
		return nil, err
	}
	return credential, nil
}
