// Search contains all search related functionality for Notion
// objects like Pages and Databases
package notion

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"tele-notion-bot/config"

	"go.uber.org/zap"
)

var apiRoot string
var apiVersion string

func init() {
	cfg := config.GetConfig()
	apiRoot = cfg.GetString("NOTION.API_ROOT")
	apiVersion = cfg.GetString("NOTION.API_VERSION")
}

func GetAllPages(notionIntToken string, slogger *zap.SugaredLogger, startCursor string, pageSize int) (*PageSearchResponse, error) {

	// Convert this to string params and remove structs
	var searchParams *SearchParams
	var startCursorValue *string

	if startCursor == "" {
		startCursorValue = nil
	} else {
		startCursorValue = &startCursor
	}

	searchParams = &SearchParams{
		Sort: &SearchSortObject{
			Direction: "descending",
			Timestamp: "last_edited_time",
		},
		Filter: &SearchFilterObject{
			Value:    "page",
			Property: "object",
		},
		StartCursor: startCursorValue,
		PageSize:    int32(pageSize),
	}

	params, err := json.Marshal(searchParams)
	if err != nil {
		slogger.Errorw(err.Error())
	}

	url := fmt.Sprintf("%s/search", apiRoot)
	payload := strings.NewReader(string(params))
	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("accept", "application/json")
	req.Header.Add("Notion-Version", apiVersion)
	req.Header.Add("content-type", "application/json")
	req.Header.Add("Authorization", notionIntToken)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		slogger.Errorw(err.Error())
	}
	defer res.Body.Close()

	if res.StatusCode == 200 {

		var response PageSearchResponse

		body, err := io.ReadAll(res.Body)
		if err != nil {
			slogger.Errorw(err.Error())
		}

		if err := json.Unmarshal(body, &response); err != nil {
			slogger.Errorw(err.Error())
		}

		return &response, nil
	} else {
		return nil, fmt.Errorf("invalid search response with status: [%s]", res.Status)
	}
}

// // GetPageTitle extracts the title of a page from a SearchResult
// func GetPageTitle(page PageObject) string {
// 	for key, prop := range page.Properties {
// 		if prop.Type == "title" {
// 			// NEED TO MODEL THIS
// 			return prop["title"][0]["plain_text"]
// 		}
// 	}
// 	return "Untitled"
// }

// -> Iterate over properties -> look for "type": "title"

// // GetDatabases retrives all databases ordered by last edited time
// func GetDatabases()

// // SearchPages performs a text search on pages in one's workspace
// func SearchPages()
