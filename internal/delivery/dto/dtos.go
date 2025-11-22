package dto

// CheckLinksRequest link verification request
type CheckLinksRequest struct {
	Links []string `json:"links"`
}

// CheckLinksResponse response with verification results
type CheckLinksResponse struct {
	Links    map[string]string `json:"links"`
	LinksNum int               `json:"links_num"`
}

// ReportRequest request for report generation
type ReportRequest struct {
	LinksList []int `json:"links_list"`
}
