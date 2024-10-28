package oauth2

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"time"

	"github.com/google/uuid"
)

var authorizationPort = 8080

var stateTruth = uuid.New().String()

type Option struct {
	RedirectUri    string
	ResponseType   string
	ApprovalPrompt string
	Scope          string
}

// Authorization Response
type AuthorizationResponse struct {
	TokenType    string `json:"token_type"`
	ExpiresAt    int    `json:"expires_at"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	AccessToken  string `json:"access_token"`
	// Athlete      map[string]interface{} `json:"athlete"`
}

func asError(message string) error {
	slog.Error(message)
	return errors.New(message)
}

func openWebBrowser(url string) {
	var err error

	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	default:
		err = fmt.Errorf("unsupported platform")
	}
	if err != nil {
		slog.Error(fmt.Sprintf("Failed to open web browser: %s", err))
	}
}

func printConsole(message string) {
	runnerPict := "🏃"
	fmt.Printf("%s %s\n", runnerPict, message)
}

func Authorize(authUrl string, tokenUrl string, clientId string, clientSecret string, option *Option) (*AuthorizationResponse, error) {
	if len(authUrl) == 0 || len(tokenUrl) == 0 || len(clientId) == 0 || len(clientSecret) == 0 {
		return nil, asError("Empty parameter found. authUrl, tokenUrl, clientId, clientSecret must be set.")
	}

	if option == nil {
		option = &Option{
			RedirectUri:    fmt.Sprintf("http://localhost:%d", authorizationPort),
			ResponseType:   "code",
			ApprovalPrompt: "auto",
			Scope:          "read,activity:read_all",
		}
	}

	fullUrl := fmt.Sprintf("%s?client_id=%s&redirect_uri=%s&response_type=%s&approval_prompt=%s&scope=%s&state=%s",
		authUrl, clientId, url.QueryEscape(option.RedirectUri), option.ResponseType, option.ApprovalPrompt, option.Scope, stateTruth)
	slog.Debug(fmt.Sprintf("Authorization URL: %s", fullUrl))

	printConsole(fmt.Sprintf("Please authorize this application in your web browser from the following URL\n%s", fullUrl))
	openWebBrowser(fullUrl)

	// channel to receive code
	ch := make(chan []string)
	// channel to wait chTimeout
	chTimeout := make(chan bool)

	// start http server to receive code
	server := &http.Server{
		Addr: fmt.Sprintf(":%d", authorizationPort),
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		slog.Debug(fmt.Sprintf("Request: %s", r.URL))
		w.Write([]byte("Authorization completed. Please close this window."))
		ch <- []string{r.URL.Query().Get("code"), r.URL.Query().Get("state")}
	})
	go server.ListenAndServe()

	// wait timeout
	go func() {
		time.Sleep(5 * time.Second)
		chTimeout <- true
	}()

	// wait for code or timeout
	var code string
	select {
	case <-chTimeout:
		close(ch)
		return nil, asError("Authorization timeout")
	case v := <-ch:
		close(chTimeout)
		code = v[0]
		state := v[1]
		if state != stateTruth {
			return nil, asError(fmt.Sprintf("Invalid state: %s", state))
		}
	}

	// request token
	// post request
	res, err := http.PostForm(tokenUrl, url.Values{
		"client_id":     {clientId},
		"client_secret": {clientSecret},
		"code":          {code},
		"grant_type":    {"authorization_code"},
	})
	if err != nil {
		return nil, asError(fmt.Sprintf("Token request failed: %s", err))
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, asError(fmt.Sprintf("Reading token response failed: %s", err))
	}

	rawToken := string(body)
	slog.Info(fmt.Sprintf("Token response: %s", rawToken))

	slog.Info(fmt.Sprintf("Token: %s", rawToken))
	go server.Shutdown(context.Background())

	var jsonToken AuthorizationResponse
	err = json.Unmarshal([]byte(rawToken), &jsonToken)
	if err != nil {
		return nil, asError(fmt.Sprintf("Failed to parse token response: %s", err))
	}
	slog.Debug(fmt.Sprintf("AuthorizationResponse: %v", jsonToken))
	slog.Debug(fmt.Sprintf("Access Token: %s", jsonToken.AccessToken))
	slog.Debug(fmt.Sprintf("Refresh Token: %s", jsonToken.RefreshToken))
	return &jsonToken, nil
}

// refresh token
func RefreshToken(tokenUrl string, clientId string, clientSecret string, refreshToken string) (*AuthorizationResponse, error) {
	if len(tokenUrl) == 0 || len(clientId) == 0 || len(clientSecret) == 0 || len(refreshToken) == 0 {
		return nil, asError("Empty parameter found. tokenUrl, clientId, clientSecret, refreshToken must be set.")
	}

	// post request
	res, err := http.PostForm(tokenUrl, url.Values{
		"client_id":     {clientId},
		"client_secret": {clientSecret},
		"refresh_token": {refreshToken},
		"grant_type":    {"refresh_token"},
	})
	if err != nil {
		return nil, asError(fmt.Sprintf("Token request failed: %s", err))
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, asError(fmt.Sprintf("Reading token response failed: %s", err))
	}

	rawToken := string(body)
	slog.Info(fmt.Sprintf("Token response: %s", rawToken))

	slog.Info(fmt.Sprintf("Token: %s", rawToken))

	var jsonToken AuthorizationResponse
	err = json.Unmarshal([]byte(rawToken), &jsonToken)
	if err != nil {
		return nil, asError(fmt.Sprintf("Failed to parse token response: %s", err))
	}
	slog.Debug(fmt.Sprintf("AuthorizationResponse: %v", jsonToken))
	slog.Debug(fmt.Sprintf("Access Token: %s", jsonToken.AccessToken))
	slog.Debug(fmt.Sprintf("Refresh Token: %s", jsonToken.RefreshToken))
	return &jsonToken, nil
}
