// Create more advanced test cases with a test DB on notion

package test

import (
	"testing"

	"tele-notion-bot/notion"
)

func TestGetNotionDatabase(t *testing.T) {

	res := notion.GetDatabase(
		Cfg.GetString("NOTION.TEST_DB_ID"),
		Cfg.GetString("NOTION.INTEGRATION_SECRET"),
		Logger,
	)

	if res.RequestStatus != 0 { // Change to cover more error codes
		t.Errorf("Request returned status %d with messsage:\n%s", res.RequestStatus, res.RequestMessage)
	}
}

func TestGetNotionDatabaseWrongSecret(t *testing.T) {

	res := notion.GetDatabase(
		Cfg.GetString("NOTION.TEST_DB_ID"),
		Cfg.GetString("NOTION.INTEGRATION_SECRET")+"_wrong",
		Logger,
	)

	if res.RequestStatus != 401 { // Change to cover more error codes
		t.Errorf("Request did not return 401 Token Invalid Error")
	}

}

func TestQueryDatabase(t *testing.T) {

	query := `{"page_size": 100, "sorts": ["property": "Date", "direction": "ascending"}]}`
	res := notion.QueryDatabase(
		Cfg.GetString("NOTION.TEST_DB_ID"),
		query,
		Cfg.GetString("NOTION.INTEGRATION_SECRET"),
		Logger,
	)

	if res.RequestStatus != 0 { // Change to cover more error codes
		t.Errorf("Request returned status %d with messsage:\n%s", res.RequestStatus, res.RequestMessage)
	}
}

func TestQueryDatabaseWrongSecret(t *testing.T) {

	query := `{"page_size": 100, "sorts": ["property": "Date", "direction": "ascending"}]}`
	res := notion.QueryDatabase(
		Cfg.GetString("NOTION.TEST_DB_ID"),
		query,
		Cfg.GetString("NOTION.INTEGRATION_SECRET")+"_wrong",
		Logger,
	)

	if res.RequestStatus != 401 {
		t.Error("Request did not return 401 Token Invalid Error")
	}

}

// func TestUpdateDatabase(t *testing.T)
