package transactions

import (
	"encoding/csv"
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
)

func ParseTransactionCSV(reader io.Reader, filename string) (ImportReport, error) {
	csvReader := csv.NewReader(reader)
	csvReader.TrimLeadingSpace = true

	records, err := csvReader.ReadAll()
	if err != nil {
		return ImportReport{}, err
	}

	report := ImportReport{Status: "validated", Filename: filename, Rows: []ImportRowResult{}}
	if len(records) == 0 {
		return report, nil
	}

	header := map[string]int{}
	for i, col := range records[0] {
		header[strings.ToLower(strings.TrimSpace(col))] = i
	}

	for _, col := range []string{"date", "label", "amount", "currency"} {
		if _, ok := header[col]; !ok {
			return ImportReport{}, fmt.Errorf("colonne obligatoire manquante: %s", col)
		}
	}

	for idx, record := range records[1:] {
		line := idx + 2
		result := ImportRowResult{Line: line, Valid: true}

		dateValue := getCSVValue(record, header["date"])
		labelValue := getCSVValue(record, header["label"])
		amountValue := getCSVValue(record, header["amount"])
		currencyValue := strings.ToUpper(getCSVValue(record, header["currency"]))

		if dateValue == "" {
			result.Errors = append(result.Errors, "date obligatoire")
		} else if _, err := time.Parse("2006-01-02", dateValue); err != nil {
			result.Errors = append(result.Errors, "date invalide, format attendu YYYY-MM-DD")
		} else {
			result.Date = dateValue
		}

		if labelValue == "" {
			result.Errors = append(result.Errors, "label obligatoire")
		} else {
			result.Label = labelValue
		}

		amount, err := strconv.ParseFloat(strings.ReplaceAll(amountValue, ",", "."), 64)
		if amountValue == "" {
			result.Errors = append(result.Errors, "amount obligatoire")
		} else if err != nil {
			result.Errors = append(result.Errors, "amount invalide")
		} else {
			result.Amount = amount
		}

		if currencyValue == "" {
			result.Errors = append(result.Errors, "currency obligatoire")
		} else {
			result.Currency = currencyValue
		}

		if len(result.Errors) > 0 {
			result.Valid = false
			report.InvalidRows++
		} else {
			report.ValidRows++
		}

		report.Rows = append(report.Rows, result)
	}

	report.DetectedRows = len(records) - 1
	return report, nil
}

func getCSVValue(record []string, index int) string {
	if index < 0 || index >= len(record) {
		return ""
	}
	return strings.TrimSpace(record[index])
}
