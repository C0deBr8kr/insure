package models

type VirusScanResponse struct {
    ResponseCode   int      `json:"response_code"`
    VerboseMessage string `json:"verbose_msg"`
    Resource string `json:"resource"`
    ScanId string `json:"scan_id"`
    PermaLink string `json:"permalink"`
    Positives int `json:"positives"`
}