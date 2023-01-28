package test

import (
	"encoding/json"
	"io"
	"os"
	"tele-notion-bot/notion"
	"testing"
)

func TestUnmarshalPageSearchResponse(t *testing.T) {

	var response notion.PageSearchResponse

	searchTestDataJSON, err := os.Open("./data/search_pages_test.json")
	if err != nil {
		t.Errorf(err.Error())
	}
	defer searchTestDataJSON.Close()

	data, err := io.ReadAll(searchTestDataJSON)
	if err != nil {
		t.Error(err.Error())
	}

	if err := json.Unmarshal(data, &response); err != nil {
		t.Error(err.Error())
	}

	propData := response.Results[0].Properties["Description & Notes"].PropertyData
	_, ok := propData.(notion.RichTextProperty)
	if ok == false {
		t.Error("Supposed to be unmarshalled to type of Rich Text Property")
	}
}
