package webserver

import (
	"context"
	"fmt"
	"net/http"
	"tele-notion-bot/internal/config"
	"tele-notion-bot/internal/database"
	"tele-notion-bot/internal/logging"
	"tele-notion-bot/internal/telegram"

	"golang.org/x/oauth2"
)

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

func handleNotionAuthCallback(w http.ResponseWriter, r *http.Request) {

	cfg := config.GetConfig()
	slogger := logging.GetLogger().Sugar()

	code := r.FormValue("code")
	conf := &oauth2.Config{
		ClientID:     cfg.GetString("NOTION.AUTH_CLIENT_ID"),
		ClientSecret: cfg.GetString("NOTION.AUTH_CLIENT_SECRET"),
		RedirectURL:  cfg.GetString("NOTION.PUB_INTEGRATION_REDIRECT_URI"),
		Scopes:       []string{"all"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://api.notion.com/v1/oauth/authorize",
			TokenURL: "https://api.notion.com/v1/oauth/token",
		},
	}

	token, err := conf.Exchange(context.Background(), code)
	if err != nil {
		slogger.Errorf("Unable to exchange authorization code for token: %v", err)
		// pass error as a message to user
		http.Redirect(w, r, cfg.GetString("TELEGRAM.TEST_BOT_URL"), http.StatusTemporaryRedirect)
		return
	}

	database.AddNotionUser(telegram.TeleUserName, token.AccessToken)
	slogger.Infof("OAuth2 workflow completed. Redirecting user back to telegram")
	http.Redirect(w, r, cfg.GetString("TELEGRAM.TEST_BOT_URL"), http.StatusSeeOther)
}
