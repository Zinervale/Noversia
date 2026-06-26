package transactions

import (
	"strings"
	"testing"
)

func TestParseTransactionCSVValidRows(t *testing.T) {
	csv := `date,label,amount,currency
2026-06-25,CARREFOUR,-82.31,EUR
2026-06-24,SALAIRE,2450.00,EUR`

	report, err := ParseTransactionCSV(strings.NewReader(csv), "test.csv")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if report.DetectedRows != 2 {
		t.Fatalf("expected 2 rows, got %d", report.DetectedRows)
	}
	if report.ValidRows != 2 {
		t.Fatalf("expected 2 valid rows, got %d", report.ValidRows)
	}
	if report.InvalidRows != 0 {
		t.Fatalf("expected 0 invalid rows, got %d", report.InvalidRows)
	}
}

func TestParseTransactionCSVInvalidRows(t *testing.T) {
	csv := `date,label,amount,currency
bad-date,CARREFOUR,abc,EUR
2026-06-24,,10,EUR`

	report, err := ParseTransactionCSV(strings.NewReader(csv), "test.csv")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if report.ValidRows != 0 {
		t.Fatalf("expected 0 valid rows, got %d", report.ValidRows)
	}
	if report.InvalidRows != 2 {
		t.Fatalf("expected 2 invalid rows, got %d", report.InvalidRows)
	}
}

func TestParseTransactionCSVMissingColumn(t *testing.T) {
	csv := `date,label,amount
2026-06-25,CARREFOUR,-82.31`

	_, err := ParseTransactionCSV(strings.NewReader(csv), "test.csv")
	if err == nil {
		t.Fatal("expected error")
	}
}
