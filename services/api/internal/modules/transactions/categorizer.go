package transactions

import "strings"

type Categorizer struct { rules []CategorizationRule }
func NewCategorizer(rules []CategorizationRule) *Categorizer { return &Categorizer{rules: rules} }
func (c *Categorizer) Categorize(label string) (CategorizationRule, bool) {
	normalized := strings.ToUpper(strings.TrimSpace(label))
	for _, rule := range c.rules {
		if !rule.Enabled { continue }
		pattern := strings.ToUpper(strings.TrimSpace(rule.Pattern))
		if rule.MatchType == "contains" && strings.Contains(normalized, pattern) { return rule, true }
		if rule.MatchType == "equals" && normalized == pattern { return rule, true }
	}
	return CategorizationRule{}, false
}

func DetectMerchant(label string) string {
	normalized := strings.ToUpper(strings.TrimSpace(label))
	if normalized == "" { return "" }
	parts := strings.FieldsFunc(normalized, func(r rune) bool { return r == ' ' || r == '-' || r == '_' || r == '.' || r == '*' })
	if len(parts) == 0 { return normalized }
	first := parts[0]
	if first == "CB" || first == "PAIEMENT" || first == "VIR" || first == "PRLV" {
		if len(parts) > 1 { return parts[1] }
	}
	return first
}
