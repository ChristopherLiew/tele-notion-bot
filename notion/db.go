// Refactor to standardise with Search

package notion

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"go.uber.org/zap"
)

// CreateDatabase creates a new Notion database on a given page within a given workspace
func CreateDatabase(notionSecret string, logger *zap.Logger) (output DatabaseObject) {

	sugar := logger.Sugar()
	url := fmt.Sprintf("%s/databases", ApiRoot)
	req, _ := http.NewRequest("POST", url, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("Notion-Version", ApiVersion)
	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", notionSecret)

	res, _ := http.DefaultClient.Do(req)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			sugar.Errorf("Unable to close Reader %s", err)
		}
	}(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		sugar.Errorw(err.Error())
	}

	if err := json.Unmarshal(body, &output); err != nil {
		sugar.Errorw(err.Error())
	}

	return
}

// GetDatabase retrives a given database's metadata and properties
func GetDatabase(databaseID string, notionSecret string, logger *zap.Logger) (output DatabaseObject) {

	sugar := logger.Sugar()
	url := fmt.Sprintf("%s/databases/%s", ApiRoot, databaseID)
	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")
	req.Header.Add("Notion-Version", ApiVersion)
	req.Header.Add("authorization", notionSecret)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		sugar.Errorf("Unable to close Reader %s", err)
	}

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			sugar.Errorw(err.Error())
		}
	}(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		sugar.Errorw(err.Error())
	}

	if err := json.Unmarshal(body, &output); err != nil {
		sugar.Errorw(err.Error())
	}

	return
}

// QueryDatabase pulls all data found in a given database based on a user defined query
func QueryDatabase(databaseId string, query string, notionSecret string, logger *zap.Logger) (output DatabaseResponse) {

	sugar := logger.Sugar()
	var inputQuery = []byte(query)
	url := fmt.Sprintf("%s/databases/%s/query", ApiRoot, databaseId)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(inputQuery))

	req.Header.Add("accept", "application/json")
	req.Header.Add("Notion-Version", ApiVersion)
	req.Header.Add("authorization", notionSecret)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		sugar.Errorw(err.Error())
	}

	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			sugar.Errorw(err.Error())
		}
	}(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		sugar.Errorw(err.Error())
	}

	if err := json.Unmarshal(body, &output); err != nil {
		sugar.Errorw(err.Error())
	}

	return
}

// UpdateDatabase updates an existing database given a user defined update query
func UpdateDatabase(databaseId string, update string, notionSecret string, logger *zap.Logger) (output DatabaseObject) {

	sugar := logger.Sugar()
	var updateQuery = []byte(update)
	url := fmt.Sprintf("%s/databases/%s/query", ApiRoot, databaseId)
	req, _ := http.NewRequest("PATCH", url, bytes.NewBuffer(updateQuery))

	req.Header.Add("accept", "application/json")
	req.Header.Add("Notion-Version", ApiVersion)
	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", notionSecret)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		sugar.Errorw(err.Error())
	}

	defer func(Body io.ReadCloser) {
		if err := Body.Close(); err != nil {
			sugar.Errorw(err.Error())
		}
	}(res.Body)

	body, err := io.ReadAll(res.Body)
	if err != nil {
		sugar.Errorw(err.Error())
	}

	if err := json.Unmarshal(body, &output); err != nil {
		sugar.Errorw(err.Error())
	}

	return
}
