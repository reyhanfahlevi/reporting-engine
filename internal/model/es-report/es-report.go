package es_report

// ParamReporting param reporting for elastic data
type ParamReporting struct {
	Index   string
	Mapping string
	Data    interface{}
}

// ParamGetReports get report using params
type ParamGetReports struct {
	Index       string
	Filter      map[string]interface{}
	RangeFilter map[string]RangeFilter
	From        int
	Size        int
}

// RangeFilter struct
type RangeFilter struct {
	From interface{} `json:"from"`
	To   interface{} `json:"to"`
}
