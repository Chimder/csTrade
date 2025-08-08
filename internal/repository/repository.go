package repository

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	Offer       OfferRepository
	User        UserRepository
	Transaction TransactionRepository
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{
		Offer:       NewOfferRepo(db),
		User:        NewUserRepository(db),
		Transaction: NewTransactionRepo(db),
	}
}
