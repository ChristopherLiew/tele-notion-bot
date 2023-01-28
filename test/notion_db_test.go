// Create more advanced test cases with a test DB on notion

package test

import (
	"testing"

	"github.com/spf13/viper"
	"go.uber.org/zap"

	"tele-notion-bot/config"
	"tele-notion-bot/logging"
	"tele-notion-bot/notion"
)

var cfg *viper.Viper
var logger *zap.Logger
var slogger *zap.SugaredLogger

func init() {
	cfg = config.GetConfig()
	logger = logging.GetLogger()
	slogger = logger.Sugar()
}

func TestGetNotionDatabase(t *testing.T) {

	res := notion.GetDatabase(
		cfg.GetString("NOTION.TEST_DB_ID"),
		cfg.GetString("NOTION.INTEGRATION_SECRET"),
		logger,
	)

	if res.RequestStatus != 0 { // Change to cover more error codes
		t.Errorf("Request returned status %d with messsage:\n%s", res.RequestStatus, res.RequestMessage)
	}
}

func TestGetNotionDatabaseWrongSecret(t *testing.T) {

	res := notion.GetDatabase(
		cfg.GetString("NOTION.TEST_DB_ID"),
		cfg.GetString("NOTION.INTEGRATION_SECRET")+"_wrong",
		logger,
	)

	if res.RequestStatus != 401 { // Change to cover more error codes
		t.Errorf("Request did not return 401 Token Invalid Error")
	}

}

func TestQueryDatabase(t *testing.T) {

	query := `{"page_size": 100, "sorts": ["property": "Date", "direction": "ascending"}]}`
	res := notion.QueryDatabase(
		cfg.GetString("NOTION.TEST_DB_ID"),
		query,
		cfg.GetString("NOTION.INTEGRATION_SECRET"),
		logger,
	)

	if res.RequestStatus != 0 { // Change to cover more error codes
		t.Errorf("Request returned status %d with messsage:\n%s", res.RequestStatus, res.RequestMessage)
	}
}

func TestQueryDatabaseWrongSecret(t *testing.T) {

	query := `{"page_size": 100, "sorts": ["property": "Date", "direction": "ascending"}]}`
	res := notion.QueryDatabase(
		cfg.GetString("NOTION.TEST_DB_ID"),
		query,
		cfg.GetString("NOTION.INTEGRATION_SECRET")+"_wrong",
		logger,
	)

	if res.RequestStatus != 401 {
		t.Error("Request did not return 401 Token Invalid Error")
	}

}

// func TestUpdateDatabase(t *testing.T)
