package virustotal

const (
	//BaseURL for virustotal
	BaseURL = "https://www.virustotal.com/vtapi/v2/file"
)

//Response from request
type Response struct {
	Scans        map[string]Data `json:"scans"`
	ResponseCode int             `json:"response_code"`
	VerboseMSG   string          `json:"verbose_msg"`
	Resource     string          `json:"resource"`
	ScanID       string          `json:"scan_id"`
	MD5          string          `json:"md5"`
	SHA1         string          `json:"sha1"`
	SHA256       string          `json:"sha256"`
	ScanDate     string          `json:"scan_date"`
	Permalink    string          `json:"permalink"`
	Positives    int             `json:"positives"`
	Total        int             `json:"total"`
}

//Data of scan results
type Data struct {
	Detected bool   `json:"detected"`
	Version  string `json:"version"`
	Result   string `json:"result"`
	Update   string `json:"update"`
}
