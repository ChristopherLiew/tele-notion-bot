// Create more advanced test cases with a test DB on notion

package test

import (
	log "github.com/sirupsen/logrus"
	"testing"

	"github.com/spf13/viper"

	"github.com/christopherliew/travel-buddy-bot/config"
	"github.com/christopherliew/travel-buddy-bot/internal/notion"
)

var Cfg *viper.Viper

func init() {
	// Set up configs for testing
	config.LoadConfig()
	Cfg = config.GetConfig()
}

func TestGetNotionDatabase(t *testing.T) {

	res := notion.GetDatabase(
		Cfg.GetString("NOTION_TEST_DB_ID"),
		Cfg.GetString("NOTION_INTEGRATION_SECRET"),
	)

	if res.RequestStatus != 0 { // Change to cover more error codes
		t.Errorf("Request returned status %d with messsage:\n%s", res.RequestStatus, res.RequestMessage)
	}
}

func TestGetNotionDatabaseWrongSecret(t *testing.T) {

	res := notion.GetDatabase(
		Cfg.GetString("NOTION_TEST_DB_ID"),
		Cfg.GetString("NOTION_INTEGRATION_SECRET")+"_wrong",
	)

	if res.RequestStatus != 401 { // Change to cover more error codes
		t.Error("Request did not return 401 Token Invalid Error")
	}

	log.Infof("Error with code: %d and message %s", res.RequestStatus, res.RequestMessage)
}

func TestQueryDatabase(t *testing.T) {

	query := `{"page_size": 100, "sorts": ["property": "Date", "direction": "ascending"}]}`
	res := notion.QueryDatabase(
		Cfg.GetString("NOTION_TEST_DB_ID"),
		query,
		Cfg.GetString("NOTION_INTEGRATION_SECRET"),
	)

	if res.RequestStatus != 0 { // Change to cover more error codes
		t.Errorf("Request returned status %d with messsage:\n%s", res.RequestStatus, res.RequestMessage)
	}
}

func TestQueryDatabaseWrongSecret(t *testing.T) {

	query := `{"page_size": 100, "sorts": ["property": "Date", "direction": "ascending"}]}`
	res := notion.QueryDatabase(
		Cfg.GetString("NOTION_TEST_DB_ID"),
		query,
		Cfg.GetString("NOTION_INTEGRATION_SECRET")+"_wrong",
	)

	if res.RequestStatus != 401 {
		t.Error("Request did not return 401 Token Invalid Error")
	}

	log.Infof("Error with code: %d and message %s", res.RequestStatus, res.RequestMessage)
}

// func TestUpdateDatabase(t *testing.T)
