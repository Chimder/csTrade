package service

import (
	"context"
	"csTrade/internal/domain/transaction"
	"csTrade/internal/repository"
)

type TransactionService struct {
	repo *repository.Repository
}

func NewTransactionService(repo *repository.Repository) *TransactionService {
	return &TransactionService{repo: repo}
}

func (trs *TransactionService) GetTransactionByBuyerID(ctx context.Context, id string) ([]transaction.TransactionDB, error) {
	return trs.repo.Transaction.GetTransactionByBuyerID(ctx, id)
}

// func (trs *TransactionService) GetTransactionByID(ctx context.Context, ID string) (*transaction.TransactionDB, error) {
// 	return trs.repo.Transaction.GetTransactionByID(ctx, ID)
// }



// func (trs *TransactionService) GetTransactionBySellerID(ctx context.Context, id string) ([]transaction.TransactionDB, error) {
// 	return trs.repo.Transaction.GetTransactionBySellerID(ctx, id)
// }
