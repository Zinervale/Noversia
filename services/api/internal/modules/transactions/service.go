package transactions

import (
	"context"
	"io"
)

type Service struct {
	repository *Repository
}

func NewService(repository *Repository) *Service {
	return &Service{repository: repository}
}

func (s *Service) List(ctx context.Context) ([]Transaction, error) {
	return s.repository.List(ctx)
}

func (s *Service) ImportCSV(ctx context.Context, reader io.Reader, filename string) (ImportReport, error) {
	report, err := ParseTransactionCSV(reader, filename)
	if err != nil {
		return ImportReport{}, err
	}

	if err := s.repository.PersistImportReport(ctx, &report); err != nil {
		return ImportReport{}, err
	}

	return report, nil
}

func (s *Service) GetImportReport(ctx context.Context, id string) (ImportReport, error) {
	return s.repository.GetImportReport(ctx, id)
}
