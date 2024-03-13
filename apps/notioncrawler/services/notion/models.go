package notion

import "strings"

type Options struct {
	Token         string
	NotionSpaceId string
	DownloadDir   string
}

type ResourceType string

const (
	ResourceTypeSpace ResourceType = "space"
	ResourceTypeBlock ResourceType = "block"
)

type ExportType string

const (
	ExportTypeMarkdown ExportType = "markdown"
)

type ExportOptions struct {
	ResourceType   ResourceType
	BlockId        string
	ExportComments bool
	ExportType     ExportType

	// Only Active for ResourceType=Workspace
	FlattenExportFiletree bool

	// Only Active for ResourceType=Block
	ExportFiles bool
}

type LoadCachedPageChunkOptions struct {
	BlockId     string
	Limit       int
	ChunkNumber int
}

type Results struct {
	Results []Result `json:"results"`
}

type Result struct {
	ID        string `json:"id"`
	EventName string `json:"eventName"`
	Request   struct {
		SpaceID              string `json:"spaceId"`
		ShouldExportComments bool   `json:"shouldExportComments"`
		ExportOptions        struct {
			ExportType            string `json:"exportType"`
			FlattenExportFiletree bool   `json:"flattenExportFiletree"`
			TimeZone              string `json:"timeZone"`
			Locale                string `json:"locale"`
		} `json:"exportOptions"`
	} `json:"request"`
	Actor struct {
		Table string `json:"table"`
		ID    string `json:"id"`
	} `json:"actor"`
	RootRequest struct {
		EventName string `json:"eventName"`
		RequestID string `json:"requestId"`
	} `json:"rootRequest"`
	Headers struct {
		IP                 string `json:"ip"`
		CityFromIP         string `json:"cityFromIp"`
		CountryCodeFromIP  string `json:"countryCodeFromIp"`
		Subdivision1FromIP string `json:"subdivision1FromIp"`
	} `json:"headers"`
	EqueuedAt int64  `json:"equeuedAt"`
	State     string `json:"state"`
	Status    *struct {
		ExportUrl     string `json:"exportUrl"`
		PagesExported int    `json:"pagesExported"`
		Type          string `json:"type"`
	} `json:"status"`
	Error string `json:"error"`
}

type PageChunkBlock struct {
	Value struct {
		ID         string                 `mapstructure:"id"`
		Version    int                    `mapstructure:"version"`
		Type       string                 `mapstructure:"type"`
		Properties map[string]interface{} `mapstructure:"properties"`
		Content    []string               `mapstructure:"content"`
		Format     struct {
			PageIcon string `mapstructure:"page_icon"`
		} `mapstructure:"format"`
		CreatedTime       int64  `mapstructure:"created_time"`
		LastEditedTime    int64  `mapstructure:"last_edited_time"`
		ParentID          string `mapstructure:"parent_id"`
		ParentTable       string `mapstructure:"parent_table"`
		Alive             bool   `mapstructure:"alive"`
		CreatedByTable    string `mapstructure:"created_by_table"`
		CreatedByID       string `mapstructure:"created_by_id"`
		LastEditedByTable string `mapstructure:"last_edited_by_table"`
		LastEditedByID    string `mapstructure:"last_edited_by_id"`
		SpaceID           string `mapstructure:"space_id"`
	} `mapstructure:"value"`
	Role string `mapstructure:"role"`
}

func (r *Result) IsSuccess() bool {
	return strings.EqualFold(r.State, "success")
}

func (r *Result) IsFailure() bool {
	return strings.EqualFold(r.State, "failure")
}
