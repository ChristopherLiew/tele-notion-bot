// TODO: Debug why cannot get auth token

package webserver

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"tele-notion-bot/internal/config"
	"tele-notion-bot/internal/database"
	"tele-notion-bot/internal/logging"
	"tele-notion-bot/internal/notion"
	"tele-notion-bot/internal/telegram"

	"go.uber.org/zap"
)

type NotionExchangeResponse struct {
	AccessToken          string `json:"access_token"`
	BotId                string `json:"bot_id"`
	DuplicatedTemplateId string `json:"duplicated_template_id"`
	Owner                string `json:"owner"`
	WorkspaceIcon        string `json:"workspace_icon"`
	WorkspaceId          string `json:"workspace_id"`
	WorkspaceName        string `json:"workspace_name"`
}

func AuthServer() (server *http.Server) {

	slogger := logging.GetLogger().Sugar()
	server = &http.Server{Addr: ":8080"}

	http.HandleFunc("/", handleMain)
	http.HandleFunc("/auth/callback", handleNotionAuthCallback)

	slogger.Info("Starting web server for notion oauth")

	go func() {
		if err := server.ListenAndServe(); err != nil {
			slogger.Error(err.Error())
		}
	}()

	// Handle server closure
	// server.Shutdown(context.Background())

	return
}

func handleMain(w http.ResponseWriter, r *http.Request) {

	var htmlIndex = `<html>
		<body>
			<p>Welcome to the Notion Telegram Bot :)</p>
		</body>
		</html>`

	fmt.Fprint(w, htmlIndex)
}

// TODO: Store token with teleusername to database
func handleNotionAuthCallback(w http.ResponseWriter, r *http.Request) {

	cfg := config.GetConfig()
	slogger := logging.GetLogger().Sugar()

	state := r.FormValue("state")
	code := r.FormValue("code")

	slogger.Infof("Obtained code from notion oauth callback: %s", code)

	accessToken, err := GetNotionAccessToken(
		state,
		code,
		cfg.GetString("NOTION.AUTH_CLIENT_ID"),
		cfg.GetString("NOTION.AUTH_CLIENT_SECRET"),
		slogger,
	)
	if err != nil {
		slogger.Error(err.Error())
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	slogger.Infof("Access token successfully retrieved: %v", *accessToken)
	database.AddNotionUser(telegram.TeleUserName, *accessToken)

	slogger.Infof("OAuth2 workflow completed. Redirecting user back to telegram")

	// redirect to tele notion bot
	http.Redirect(w, r, cfg.GetString("TELEGRAM.TEST_BOT_URL"), http.StatusSeeOther)
}

func GetNotionAccessToken(state string, code string, clientId string, clientSecret string, slogger *zap.SugaredLogger) (accessToken *string, err error) {

	slogger.Info("Retrieving access token")

	url := fmt.Sprintf("%s/oauth/token", notion.ApiRoot)
	payload := strings.NewReader(fmt.Sprintf(`{
		"grant_type": "authorization_code",
		"code": "%s",
		"redirect_uri": "%s"
	}`, code, "http://localhost:8080/auth/callback"))
	req, _ := http.NewRequest("POST", url, payload)
	authorization := fmt.Sprintf(`Basic "%s:%s"`, clientId, clientSecret)

	req.Header.Add("accept", "application/json")
	req.Header.Add("Notion-Version", notion.ApiVersion)
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", authorization)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		slogger.Errorf("Failed to retrieve access token: %s", err.Error())
		return nil, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		slogger.Error(err.Error())
		return nil, err
	}

	var accessResponse NotionExchangeResponse
	err = json.Unmarshal(body, &accessResponse)

	return &accessResponse.AccessToken, err
}
