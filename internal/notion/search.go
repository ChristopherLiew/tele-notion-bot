package notion

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"go.uber.org/zap"
)

func searchAPICall(notionIntToken string, slogger *zap.SugaredLogger, pageOrDB string, startCursor *string, pageSize int) []byte {

	searchParams := &SearchParams{
		Sort: &SearchSortObject{
			Direction: "descending",
			Timestamp: "last_edited_time",
		},
		Filter: &SearchFilterObject{
			Value:    pageOrDB,
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
	isFailedRequest := res.StatusCode != 200
	if (err != nil) || isFailedRequest {
		slogger.Errorw(err.Error())
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		slogger.Errorw(err.Error())
	}

	return body
}

func GetAllPages(notionIntToken string, slogger *zap.SugaredLogger, startCursor *string, pageSize int) (*PageSearchResponse, error) {

	var response PageSearchResponse

	body := searchAPICall(
		notionIntToken,
		slogger,
		"page",
		startCursor,
		pageSize,
	)
	err := json.Unmarshal(body, &response)

	return &response, err
}

// Turn this into go routine
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
				hasTitle := len(titleProperty) > 0
				hasURL := page.URL != ""

				if hasTitle && hasURL {
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

func GetAllDBs(notionIntToken string, slogger *zap.SugaredLogger, startCursor *string, pageSize int) (*DBSearchResponse, error) {

	var response DBSearchResponse

	body := searchAPICall(
		notionIntToken,
		slogger,
		"database",
		startCursor,
		pageSize,
	)
	err := json.Unmarshal(body, &response)

	return &response, err

}
