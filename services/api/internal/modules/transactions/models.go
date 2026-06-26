package transactions

type Transaction struct {
	ID       string  `json:"id"`
	Label    string  `json:"label"`
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
	Date     string  `json:"date"`
}

type ImportRowResult struct {
	Line          int      `json:"line"`
	Valid         bool     `json:"valid"`
	Date          string   `json:"date,omitempty"`
	Label         string   `json:"label,omitempty"`
	Amount        float64  `json:"amount,omitempty"`
	Currency      string   `json:"currency,omitempty"`
	Errors        []string `json:"errors,omitempty"`
	TransactionID string   `json:"transactionId,omitempty"`
	Duplicate     bool     `json:"duplicate,omitempty"`
}

type ImportReport struct {
	ID           string            `json:"id"`
	Status       string            `json:"status"`
	Filename     string            `json:"filename"`
	DetectedRows int              `json:"detectedRows"`
	ValidRows    int              `json:"validRows"`
	InvalidRows  int              `json:"invalidRows"`
	Rows         []ImportRowResult `json:"rows"`
}
