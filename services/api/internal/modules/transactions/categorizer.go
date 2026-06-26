package transactions

import "strings"

type Categorizer struct {
	rules []CategorizationRule
}

func NewCategorizer(rules []CategorizationRule) *Categorizer {
	return &Categorizer{rules: rules}
}

func (c *Categorizer) Categorize(label string) (CategorizationRule, bool) {
	normalized := strings.ToUpper(strings.TrimSpace(label))
	for _, rule := range c.rules {
		if !rule.Enabled {
			continue
		}
		pattern := strings.ToUpper(strings.TrimSpace(rule.Pattern))
		switch rule.MatchType {
		case "contains":
			if strings.Contains(normalized, pattern) {
				return rule, true
			}
		case "equals":
			if normalized == pattern {
				return rule, true
			}
		}
	}
	return CategorizationRule{}, false
}
