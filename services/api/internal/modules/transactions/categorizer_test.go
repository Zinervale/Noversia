package transactions

import "testing"

func TestCategorizerContains(t *testing.T) {
	c := NewCategorizer([]CategorizationRule{
		{Pattern: "CARREFOUR", MatchType: "contains", CategoryID: "cat_courses", CategoryName: "Courses", Enabled: true},
	})

	rule, ok := c.Categorize("CARREFOUR MARKET")
	if !ok {
		t.Fatal("expected match")
	}
	if rule.CategoryName != "Courses" {
		t.Fatalf("expected Courses, got %s", rule.CategoryName)
	}
}

func TestCategorizerDisabledRule(t *testing.T) {
	c := NewCategorizer([]CategorizationRule{
		{Pattern: "NETFLIX", MatchType: "contains", CategoryID: "cat_sub", Enabled: false},
	})

	_, ok := c.Categorize("NETFLIX")
	if ok {
		t.Fatal("expected no match")
	}
}
