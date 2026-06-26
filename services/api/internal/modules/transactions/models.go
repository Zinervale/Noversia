package transactions

type Transaction struct {
	ID string `json:"id"`
	Label string `json:"label"`
	Amount float64 `json:"amount"`
	Currency string `json:"currency"`
	Date string `json:"date"`
	MerchantID string `json:"merchantId,omitempty"`
	MerchantName string `json:"merchantName,omitempty"`
	CategoryID string `json:"categoryId,omitempty"`
	CategoryName string `json:"categoryName,omitempty"`
}

type Category struct { ID string `json:"id"`; Name string `json:"name"` }

type CategorizationRule struct {
	ID string `json:"id"`
	Pattern string `json:"pattern"`
	MatchType string `json:"matchType"`
	CategoryID string `json:"categoryId"`
	CategoryName string `json:"categoryName,omitempty"`
	Priority int `json:"priority"`
	ConfidenceScore float64 `json:"confidenceScore"`
	Enabled bool `json:"enabled"`
}

type RuleSuggestion struct {
	Pattern string `json:"pattern"`
	CategoryID string `json:"categoryId"`
	CategoryName string `json:"categoryName"`
	Occurrences int `json:"occurrences"`
	Reason string `json:"reason"`
}

type ApplyRuleSuggestionRequest struct {
	Pattern string `json:"pattern"`
	CategoryID string `json:"categoryId"`
	Priority int `json:"priority"`
}

type UpdateCategoryRequest struct { CategoryID string `json:"categoryId"`; Reason string `json:"reason"` }

type ImportRowResult struct {
	Line int `json:"line"`
	Valid bool `json:"valid"`
	Date string `json:"date,omitempty"`
	Label string `json:"label,omitempty"`
	Amount float64 `json:"amount,omitempty"`
	Currency string `json:"currency,omitempty"`
	MerchantID string `json:"merchantId,omitempty"`
	MerchantName string `json:"merchantName,omitempty"`
	CategoryID string `json:"categoryId,omitempty"`
	CategoryName string `json:"categoryName,omitempty"`
	ConfidenceScore float64 `json:"confidenceScore,omitempty"`
	Errors []string `json:"errors,omitempty"`
	TransactionID string `json:"transactionId,omitempty"`
	Duplicate bool `json:"duplicate,omitempty"`
}

type ImportReport struct {
	ID string `json:"id"`
	Status string `json:"status"`
	Filename string `json:"filename"`
	DetectedRows int `json:"detectedRows"`
	ValidRows int `json:"validRows"`
	InvalidRows int `json:"invalidRows"`
	Rows []ImportRowResult `json:"rows"`
}
