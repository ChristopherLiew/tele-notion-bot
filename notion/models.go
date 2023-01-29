package notion

import (
	"encoding/json"
	"fmt"
)

// Users

type UserObject struct {
	Object    string       `json:"object"`
	Id        string       `json:"id"`
	Type      string       `json:"type"`
	Name      string       `json:"name"`
	AvatarURL string       `json:"avatar_url"`
	Person    PersonObject `json:"person,omitempty"`
	Bot       BotObject    `json:"bot,omitempty"`
}

type PersonObject struct {
	Email string `json:"email"`
}

type BotObject struct {
	Owner         map[string]string `json:"owner"`
	WorkspaceName string            `json:"workspace_name"`
}

// Pages

type PageObject struct {
	Object         string                  `json:"object"`
	Id             string                  `json:"id"`
	CreatedTime    string                  `json:"created_time"`
	CreatedBy      UserObject              `json:"created_by"`
	LastEditedTime string                  `json:"last_edited_time"`
	LastEditedBy   UserObject              `json:"last_edited_by"`
	Archived       bool                    `json:"archived"`
	Icon           EmojiObject             `json:"icon,omitempty"`
	Cover          interface{}             `json:"cover,omitempty"`
	Properties     map[string]PageProperty `json:"properties"`
	Parent         ParentProperty          `json:"parent"`
	URL            string                  `json:"url"`
}

type EmojiObject struct {
	Type  string `json:"type"`
	Emoji string `json:"emoji"`
}

type PageProperty struct {
	Id           string      `json:"id"`
	Type         string      `json:"type"` // 'Enums' not useful since unchecked marshalling is allowed
	PropertyData interface{} // Can be any type of page property value (E.g. Multi select, Rich text etc.)
}

func (p *PageProperty) UnmarshalJSON(data []byte) (err error) {

	// Use new type to prevent infinite recursion
	type _PageProperty PageProperty
	var pageProperty _PageProperty

	// Unmarshal into new type before casting back to original
	if err = json.Unmarshal(data, &pageProperty); err == nil {
		*p = PageProperty(pageProperty)
	}

	// Handle dynamic nature of properties
	m := make(map[string]interface{})
	var val interface{}
	if err = json.Unmarshal(data, &m); err == nil {
		val = m[p.Type]
	}

	// simple types
	simpleType := map[string]bool{
		"checkbox":         true,
		"created_time":     true,
		"last_edited_time": true,
		"email":            true,
		"phone_number":     true,
		"url":              true,
		"number":           true,
	}

	if simpleType[p.Type] {
		if p.Type == "number" {
			p.PropertyData = val.(float32)
		} else {
			p.PropertyData = val
		}
		return
	}

	// complex types (disgusting code but laze, will refactor in future)
	switch p.Type {
	case "created_by", "last_edited_by":
		var uo UserObject
		if err = json.Unmarshal(data, &uo); err == nil {
			p.PropertyData = uo
		}
	case "date":
		var dt DateProperty
		if err = json.Unmarshal(data, &dt); err == nil {
			p.PropertyData = dt
		}
	case "files":
		var file FileProperty
		if err = json.Unmarshal(data, &file); err == nil {
			p.PropertyData = file
		}
	case "formula":
		var fml FormulaProperty
		if err = json.Unmarshal(data, &fml); err == nil {
			p.PropertyData = fml
		}
	case "select":
		var sel SelectProperty
		if err = json.Unmarshal(data, &sel); err == nil {
			p.PropertyData = sel
		}
	case "multi_select":
		var mulSel MultiSelectProperty
		if err = json.Unmarshal(data, &mulSel); err == nil {
			p.PropertyData = mulSel
		}
	case "people":
		var ppl PeopleProperty
		if err = json.Unmarshal(data, &ppl); err == nil {
			p.PropertyData = ppl
		}
	case "relation":
		var rel RelationProperty
		if err = json.Unmarshal(data, &rel); err == nil {
			p.PropertyData = rel
		}
	case "rollup":
		var roll RollupProperty
		if err = json.Unmarshal(data, &roll); err == nil {
			p.PropertyData = roll
		}
	case "rich_text":
		var rtp RichTextProperty
		if err = json.Unmarshal(data, &rtp); err == nil {
			p.PropertyData = rtp
		}
	case "status":
		var sta StatusProperty
		if err = json.Unmarshal(data, &sta); err == nil {
			p.PropertyData = sta
		}
	case "title":
		var ttl TitleProperty
		if err = json.Unmarshal(data, &ttl); err == nil {
			p.PropertyData = ttl
		}
	default:
		err = fmt.Errorf("unknown page property type: [%s]", p.Type)
		return
	}

	return
}

// Property types

type DateProperty struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type FileProperty struct {
	Files []FileObjects `json:"files"`
}

type FileObjects struct {
	Name     string                 `json:"name"`
	Type     string                 `json:"type"`
	External ExternalFileObject     `json:"external,omitempty"`
	File     NotionHostedFileObject `json:"file,omitempty"`
}

type ExternalFileObject struct {
	URL string `json:"url"`
}

type NotionHostedFileObject struct {
	URL        string `json:"url"`
	ExpiryTime string `json:"expiry_time"`
}

type FormulaProperty struct {
	Type    string  `json:"type"` // enum
	Boolean bool    `json:"boolean,omitempty"`
	Date    string  `json:"date,omitempty"`
	Number  float32 `json:"number,omitempty"`
	String  string  `json:"string,omitempty"`
}

type MultiSelectProperty struct {
	MultiSelectOptions []SelectProperty
}

type SelectProperty struct {
	Color string `json:"color"`
	Id    string `json:"id"`
	Name  string `json:"name"`
}

type StatusProperty struct {
	Color string `json:"color"`
	Id    string `json:"id"`
	Name  string `json:"name"`
}

type PeopleProperty struct {
	People []PersonObject `json:"people"`
}

type RollupProperty struct {
	Type        string `json:"type"`
	Function    string `json:"function"`
	Array       string `json:"array,omitempty"`
	Date        string `json:"date,omitempty"`
	Incomplete  string `json:"incomplete,omitempty"`
	Number      string `json:"number,omitempty"`
	Unsupported string `json:"unsupported,omitempty"`
}

type RelationProperty struct {
	HasMore  bool            `json:"has_more"`
	Relation []PageReference `json:"relation"`
}

type PageReference struct {
	Id string `json:"id"`
}

type TitleProperty struct {
	Title []RichTextObject `json:"title"`
}

type RichTextProperty struct {
	RichText []RichTextObject `json:"rich_text"`
}

type RichTextObject struct {
	PlainText   string           `json:"plain_text"`
	HRef        string           `json:"href,omitempty"`
	Annotations AnnotationObject `json:"annotations"`
	Type        string           `json:"type"`
	Text        TextObject       `json:"text,omitempty"`
	Mention     MentionObject    `json:"mention,omitempty"`
	Equation    EquationObject   `json:"equation,omitempty"`
}

type AnnotationObject struct {
	Bold          bool   `json:"bold"`
	Italic        bool   `json:"italic"`
	Strikethrough bool   `json:"strikethrough"`
	Underline     bool   `json:"underline"`
	Code          bool   `json:"code"`
	Color         string `json:"color"`
}

type EquationObject struct {
	Expression string `json:"expression"`
}

type MentionObject struct {
	Type        string                `json:"type"`
	Database    DatabaseMentionObject `json:"database,omitempty"`
	Date        DateProperty          `json:"date,omitempty"`
	LinkPreview LinkObject            `json:"link_preview,omitempty"`
	Page        PageReference         `json:"page,omitempty"`
	User        UserMentionObject     `json:"user,omitempty"`
}

type DatabaseMentionObject struct {
	Id string `json:"id"`
}

type LinkObject struct {
	URL string `json:"url"`
}

type TemplateMentionObject struct {
	Type                string `json:"type"`
	TemplateMentionDate string `json:"template_mention_date,omitempty"`
	TemplateMentionUser string `json:"template_mention_user,omitempty"`
}

type UserMentionObject struct {
	Object string `json:"object"`
	Id     string `json:"id"`
}

type TextObject struct {
	Content string     `json:"content"`
	Link    LinkObject `json:"link,omitempty"`
}

type ParentProperty struct {
	Type        string `json:"type"`
	DatabaseId  string `json:"database_id,omitempty"`
	PageId      string `json:"page_id,omitempty"`
	IsWorkspace bool   `json:"workspace,omitempty"`
	BlockId     string `json:"block_id,omitempty"`
}

// Databases

type DatabaseResponse struct {
	Object         string        `json:"object"`
	Results        []interface{} `json:"results"`
	NextCursor     string        `json:"next_cursor"`
	HasMore        bool          `json:"has_more"`
	Type           string        `json:"type"`
	Page           interface{}   `json:"page"`
	RequestStatus  int           `json:"status"` // Think of a better way to handle this
	RequestMessage string        `json:"message"`
}

type DatabaseObject struct {
	Object         string           `json:"object"`
	Id             string           `json:"id"`
	CreatedTime    string           `json:"created_time"`
	CreatedBy      interface{}      `json:"created_by"`
	LastEditedTime string           `json:"last_edited_time"`
	LastEditedBy   interface{}      `json:"last_edited_by"`
	Title          interface{}      `json:"title"`
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

// Search

type SearchSortObject struct {
	Direction string `json:"direction"` // Only supports "ascending" or "descending"
	Timestamp string `json:"timestamp"` // E.g. Last Edited Time
}

type SearchFilterObject struct {
	Value    string `json:"value"`    // Only supports "page" or "database"
	Property string `json:"property"` // Only supports "object"
}

type SearchParams struct {
	Query       *string             `json:"query,omitempty"`
	Sort        *SearchSortObject   `json:"sort,omitempty"`
	Filter      *SearchFilterObject `json:"filter,omitempty"`
	StartCursor *string             `json:"start_cursor,omitempty"`
	PageSize    int32               `json:"page_size"` // Max of 100
}

type PageSearchResponse struct {
	Object         string       `json:"object"`
	Results        []PageObject `json:"results"`
	NextCursor     string       `json:"next_cursor"`
	HasMore        bool         `json:"has_more"`
	Type           string       `json:"type"`
	PageOrDatabase interface{}  `json:"page_or_database"`
	Status         int          `json:"status,omitempty"`
	Code           string       `json:"code,omitempty"`
	Message        string       `json:"message,omitempty"`
}

type PageSnippet struct {
	Title string
	Icon  string // Uses a default icon if nothing
	URL   string
}
