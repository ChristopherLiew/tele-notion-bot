// Search contains all search related functionality for Notion
// objects like Pages and Databases
package notion

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

// GetAllPages retrieves all notion pages and returns a PageSearchResponse which mirrors the search api endpoint's response body.
func GetAllPages(notionIntToken string, slogger *zap.SugaredLogger, startCursor *string, pageSize int) (*PageSearchResponse, error) {

	// used struct to handle nullable fields (e.g. start cursor)
	searchParams := &SearchParams{
		Sort: &SearchSortObject{
			Direction: "descending",
			Timestamp: "last_edited_time",
		},
		Filter: &SearchFilterObject{
			Value:    "page",
			Property: "object",
		},
		StartCursor: startCursor,
		PageSize:    int32(pageSize),
	}

	params, err := json.Marshal(searchParams)
	if err != nil {
		slogger.Errorw(err.Error())
	}

	url := fmt.Sprintf("%s/search", ApiRoot)
	payload := strings.NewReader(string(params))
	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("accept", "application/json")
	req.Header.Add("Notion-Version", ApiVersion)
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

// GetPageSnippet retrieves a notion page's key information: Title, Icon (if available) and URL (if available).
func GetPageSnippets(notionPages []PageObject, slogger *zap.SugaredLogger) (snippets []PageSnippet) {

	var titleText string
	arr := make([]PageSnippet, 0)

	for _, page := range notionPages {
		for _, propertyValue := range page.Properties {
			if propertyValue.Type == "title" {

				titleProperty := propertyValue.
					PropertyData.(TitleProperty).
					Title

				// filter any blank pages
				if (len(titleProperty) > 0) && (page.URL != "") {
					titleText = titleProperty[0].PlainText
					snippet := PageSnippet{
						Title: titleText,
						URL:   page.URL,
						Icon:  page.Icon.Emoji,
					}
					arr = append(arr, snippet)
				}
			}
		}
	}
	return arr
}
