package transactions

import (
	"context"
	"io"
)

type Service struct { repository *Repository }
func NewService(repository *Repository) *Service { return &Service{repository: repository} }
func (s *Service) List(ctx context.Context) ([]Transaction,error) { return s.repository.List(ctx) }
func (s *Service) ListCategories(ctx context.Context) ([]Category,error) { return s.repository.ListCategories(ctx) }
func (s *Service) ListRules(ctx context.Context) ([]CategorizationRule,error) { return s.repository.ListRules(ctx) }
func (s *Service) CreateRule(ctx context.Context, rule CategorizationRule) (CategorizationRule,error) { return s.repository.CreateRule(ctx, rule) }
func (s *Service) ApplyRuleSuggestion(ctx context.Context, req ApplyRuleSuggestionRequest) (CategorizationRule,error) {
	return s.repository.CreateRule(ctx, CategorizationRule{Pattern:req.Pattern, CategoryID:req.CategoryID, MatchType:"contains", Priority:req.Priority, ConfidenceScore:0.90})
}
func (s *Service) UpdateCategory(ctx context.Context, txID, categoryID, reason string) (Transaction,error) { return s.repository.UpdateTransactionCategory(ctx, txID, categoryID, reason) }
func (s *Service) ListRuleSuggestions(ctx context.Context) ([]RuleSuggestion,error) { return s.repository.ListRuleSuggestions(ctx) }
func (s *Service) ImportCSV(ctx context.Context, reader io.Reader, filename string) (ImportReport,error) {
	report, err := ParseTransactionCSV(reader, filename); if err != nil { return ImportReport{}, err }
	rules, err := s.repository.ListRules(ctx); if err != nil { return ImportReport{}, err }
	categorizer := NewCategorizer(rules)
	for i := range report.Rows { if report.Rows[i].Valid { if rule, ok := categorizer.Categorize(report.Rows[i].Label); ok { report.Rows[i].CategoryID=rule.CategoryID; report.Rows[i].CategoryName=rule.CategoryName; report.Rows[i].ConfidenceScore=rule.ConfidenceScore } } }
	if err := s.repository.PersistImportReport(ctx, &report); err != nil { return ImportReport{}, err }
	return report,nil
}
func (s *Service) GetImportReport(ctx context.Context,id string) (ImportReport,error) { return s.repository.GetImportReport(ctx,id) }
