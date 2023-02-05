// TODO:
// - Add test case for pagination

package test

import (
	"tele-notion-bot/internal/config"
	"tele-notion-bot/internal/logging"
	"tele-notion-bot/internal/notion"
	"testing"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

var Cfg *viper.Viper
var Logger *zap.Logger
var Slogger *zap.SugaredLogger

func init() {
	Cfg = config.GetConfig()
	Logger = logging.GetLogger()
	Slogger = Logger.Sugar()
}

var PageResults []notion.PageObject

func TestGetAllPages(t *testing.T) {

	res, err := notion.GetAllPages(
		Cfg.GetString("NOTION.INTEGRATION_SECRET"),
		Slogger,
		nil,
		10,
	)
	if err != nil {
		t.Errorf(err.Error())
	}

	if res.Object == "error" {
		t.Errorf("Code: %s\nMessage: %s\n", res.Code, res.Message)
	}

	// Store results for downstream test cases
	PageResults = res.Results
}

func TestGetAllPagesVarySize(t *testing.T) {

	pageSizes := []int{10, 20, 50, 100, 200, 400}

	for _, size := range pageSizes {
		res, err := notion.GetAllPages(
			Cfg.GetString("NOTION.INTEGRATION_SECRET"),
			Slogger,
			nil,
			size,
		)
		if err != nil {
			t.Errorf(err.Error())
		}
		if res.Object == "error" {
			t.Errorf("Code: %s\nMessage: %s\n", res.Code, res.Message)
		}
	}
}

func TestGetPagesSnippets(t *testing.T) {

	snippets := notion.GetPageSnippets(PageResults, Slogger)

	// Checks for issues in extracting information for snippets
	if (len(snippets) == 0) && (len(PageResults) > 0) {
		t.Error("no snippets found from page search results!")
	}

	// Checks for empty snippets (possible but uninformative) [P0 - Fix this]
	for _, snippet := range snippets {
		if snippet.Title == "" {
			t.Error("snippet with empty title found - snippets must have titles")
		} else if snippet.URL == "" {
			t.Error("snippet with empty url found - snippets must have URLs ")
		}
	}

}

func TestGetAllDatabases(t *testing.T) {

	res, err := notion.GetAllDBs(
		Cfg.GetString("NOTION.INTEGRATION_SECRET"),
		Slogger,
		nil,
		10,
	)
	if err != nil {
		t.Errorf(err.Error())
	}

	if res.Object == "error" {
		t.Errorf("Code: %s\nMessage: %s\n", res.Code, res.Message)
	}

}

func TestGetAllDBsVarySize(t *testing.T) {

	pageSizes := []int{10, 20, 50, 100, 200, 400}

	for _, size := range pageSizes {
		res, err := notion.GetAllDBs(
			Cfg.GetString("NOTION.INTEGRATION_SECRET"),
			Slogger,
			nil,
			size,
		)
		if err != nil {
			t.Errorf(err.Error())
		}
		if res.Object == "error" {
			t.Errorf("Code: %s\nMessage: %s\n", res.Code, res.Message)
		}
	}
}
