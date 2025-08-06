package sqlstore

import (
	"cryptoObserver/internal/app/model"
	"cryptoObserver/internal/app/store/sqlstore/utils"
	"database/sql"
)

type CurrencyInterface interface {
	AddCurrency(currency string) error
	RemoveCurrency(currency string) error
	GetPrice(coin string, timestamp int64) (model.Decimal, error)
	GetCurrencyList() ([]string, error)
	UpdatePrice(coin string, price model.Decimal, timestamp int64) error
}

type CurrencyRepository struct {
	store *Store
}

func (r *CurrencyRepository) AddCurrency(currency string) error {
	_, err := r.store.db.Exec(
		"INSERT INTO currencies (symbol) VALUES ($1) ON CONFLICT (symbol) DO NOTHING",
		currency,
	)
	return err
}

func (r *CurrencyRepository) RemoveCurrency(currency string) error {
	_, err := r.store.db.Exec(
		"DELETE FROM currencies WHERE symbol = $1",
		currency,
	)
	return err
}

func (r *CurrencyRepository) GetPrice(coin string, timestamp int64) (model.Decimal, error) {
	var price string
	err := r.store.db.QueryRow(
		`SELECT cp.price
		 FROM currency_prices cp
		 JOIN currencies c ON cp.currency_id = c.id
		 WHERE c.symbol = $1
		 ORDER BY ABS(cp.timestamp - $2)
		 LIMIT 1`,
		coin, timestamp,
	).Scan(&price)
	if err == sql.ErrNoRows {
		return model.Decimal{}, nil
	}
	decimal, err := utils.ParseDecimal(price)
	if err != nil {
		return model.Decimal{}, err
	}
	return decimal, nil
}

func (r *CurrencyRepository) GetCurrencyList() ([]string, error) {
	var currencies []string
	rows, err := r.store.db.Query("SELECT symbol FROM currencies")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var currency string
		if err := rows.Scan(&currency); err != nil {
			return nil, err
		}
		currencies = append(currencies, currency)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return currencies, nil
}

func (r *CurrencyRepository) UpdatePrice(coin string, price model.Decimal, timestamp int64) error {
	// Получаем id валюты по символу
	var currencyID int
	err := r.store.db.QueryRow(
		"SELECT id FROM currencies WHERE symbol = $1",
		coin,
	).Scan(&currencyID)
	if err != nil {
		return err
	}

	// Преобразуем Decimal в строку
	priceStr := utils.DecimalToString(price)

	// Вставляем или обновляем цену
	_, err = r.store.db.Exec(
		`INSERT INTO currency_prices (currency_id, price, timestamp)
   VALUES ($1, $2, $3)
   ON CONFLICT (currency_id, timestamp) DO UPDATE SET price = EXCLUDED.price`,
		currencyID, priceStr, timestamp,
	)
	return err
}
