package notion

// Databases
// reference: https://developers.notion.com/reference/database

// DatabaseQuery Notion database query result object
type DatabaseQuery struct {
	Object         string        `json:"object"`
	Results        []interface{} `json:"results"`
	NextCursor     string        `json:"next_cursor"`
	HasMore        bool          `json:"has_more"`
	Type           string        `json:"type"`
	Page           interface{}   `json:"page"`
	RequestStatus  int           `json:"status"` // Think of a better way to handle this
	RequestMessage string        `json:"message"`
}

// DatabaseObject Notion database object
type DatabaseObject struct {
	Object         string           `json:"object"`
	Id             string           `json:"id"`
	CreatedTime    string           `json:"created_time"`
	CreatedBy      interface{}      `json:"created_by"`
	LastEditedTime string           `json:"last_edited_time"`
	LastEditedBy   interface{}      `json:"last_edited_by"`
	Title          []RichTextObject `json:"title"`
	Description    []RichTextObject `json:"description"`
	Icon           interface{}      `json:"icon"`
	Cover          interface{}      `json:"cover"`
	Properties     interface{}      `json:"properties"`
	Parent         interface{}      `json:"parent"`
	Url            string           `json:"url"`
	Archived       bool             `json:"archived"`
	IsInline       bool             `json:"is_inline"`
	RequestStatus  int              `json:"status"` // Think of a better way to handle this
	RequestMessage string           `json:"message"`
}

// PropertyType Database property enum
type PropertyType string

const (
	TITLE            PropertyType = "TITLE"
	RICH_TEXT        PropertyType = "RICH_TEXT"
	NUMBER           PropertyType = "NUMBER"
	SELECT           PropertyType = "SELECT"
	MULTI_SELECT     PropertyType = "MULTI_SELECT"
	DATE             PropertyType = "DATE"
	PEOPLE           PropertyType = "PEOPLE"
	FILES            PropertyType = "FILES"
	CHECKBOX         PropertyType = "CHECKBOX"
	URL              PropertyType = "URL"
	EMAIL            PropertyType = "EMAIL"
	PHONE_NUMBER     PropertyType = "PHONE_NUMBER"
	FORMULA          PropertyType = "FORMULA"
	RELATION         PropertyType = "RELATION"
	ROLLUP           PropertyType = "ROLLUP"
	CREATED_TIME     PropertyType = "CREATED_TIME"
	CREATED_BY       PropertyType = "CREATED_BY"
	LAST_EDITED_TIME PropertyType = "LAST_EDITED_TIME"
	LAST_EDITED_BY   PropertyType = "LAST_EDITED_BY"
	STATUS           PropertyType = "STATUS"
)

// PropertyObject Notion property object (Typically nested within Databases)
type PropertyObject struct {
	Id   string       `json:"id"`
	Type PropertyType `json:"type"`
	Name string       `json:"name"`
}

type RichTextObject struct {
	PlainText   string      `json:"plain_text"`
	HRef        string      `json:"href,omitempty"`
	Annotations interface{} `json:"annotations"`
	Type        string      `json:"type"`
}

// DB filter objects
