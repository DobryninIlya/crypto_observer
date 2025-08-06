package sqlstore

import "database/sql"

type Store struct {
	db                 *sql.DB
	currencyRepository CurrencyInterface
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s *Store) Currency() CurrencyInterface {
	if s.currencyRepository != nil {
		return s.currencyRepository
	}

	s.currencyRepository = &CurrencyRepository{
		store: s,
	}

	return s.currencyRepository
}
