package db

import (
    "database/sql"
    "errors"
    "log"
    "os"
	"time"
)

var ErrUnknownCoin = errors.New("unknown coin name")

type CoinData struct {
	Name string `json:"name"`
	CurrentPrice float64 `json:"current_price"`
	LastUpdate time.Time `json:"last_updated"`
}

type CoinDataWithSymbol struct {
	Symbol       string    `json:"symbol"`
	Name         string    `json:"name"`
	CurrentPrice float64   `json:"current_price"`
	LastUpdate   time.Time `json:"last_updated"`
}

type CryptoDB struct {
	conn *sql.DB
}

func NewCryptoDB() (*CryptoDB, error) {
	dsn := os.Getenv("DB_DSN")
	if dsn == "" {
		return nil, errors.New("DB_DSN environment variable is required")
	}

	db, err := sql.Open("postgres", dsn)
    if err != nil {
        return nil, err
    }

	if err := db.Ping(); err != nil {
        return nil, err
    }

    if err := runMigrations(db); err != nil {
        return nil, err
    }

    return &CryptoDB{conn: db}, nil
}

func (cdb *CryptoDB) Insert(symbol string, data CoinData) error {
	_, err := cdb.conn.Exec(`
		INSERT INTO crypto (symbol, name, current_price, last_update)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (symbol) DO UPDATE 
		    SET current_price = EXCLUDED.current_price,
		        last_update   = EXCLUDED.last_update,
		        name          = EXCLUDED.name
	`, symbol, data.Name, data.CurrentPrice, data.LastUpdate)
	
	if err != nil {
		log.Printf("Failed to insert/update crypto %s: %v", symbol, err)
		return err
	}

	log.Printf("Successfully updated %s (%s) price: $%.2f", symbol, data.Name, data.CurrentPrice)
	return nil
}

func (cdb *CryptoDB) Get(symbol string) (CoinData, error) {
	var data CoinData
	err := cdb.conn.QueryRow(`
		SELECT name, current_price, last_update 
		FROM crypto WHERE symbol = $1
	`, symbol).Scan(&data.Name, &data.CurrentPrice, &data.LastUpdate)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return CoinData{}, ErrUnknownCoin
		}
		return CoinData{}, err
	}

	return data, nil
}

func (cdb *CryptoDB) GetAll() (map[string]CoinData, error) {
	rows, err := cdb.conn.Query(`
		SELECT symbol, name, current_price, last_update 
		FROM crypto 
		ORDER BY symbol
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cryptos := make(map[string]CoinData)
	for rows.Next() {
		var symbol string
		var data CoinData
		
		err := rows.Scan(&symbol, &data.Name, &data.CurrentPrice, &data.LastUpdate)
		if err != nil {
			return nil, err
		}
		
		cryptos[symbol] = data
	}

	return cryptos, nil
}

func (cdb *CryptoDB) GetAllSlice() ([]CoinDataWithSymbol, error) {
	rows, err := cdb.conn.Query(`
		SELECT symbol, name, current_price, last_update 
		FROM crypto 
		ORDER BY symbol
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cryptos []CoinDataWithSymbol
	for rows.Next() {
		var crypto CoinDataWithSymbol
		
		err := rows.Scan(&crypto.Symbol, &crypto.Name, &crypto.CurrentPrice, &crypto.LastUpdate)
		if err != nil {
			return nil, err
		}
		
		cryptos = append(cryptos, crypto)
	}

	return cryptos, nil
}

func (cdb *CryptoDB) Delete(symbol string) error {
	res, err := cdb.conn.Exec(`DELETE FROM crypto WHERE symbol = $1`, symbol)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	
	if rows == 0 {
		return ErrUnknownCoin
	}

	log.Printf("Successfully deleted crypto: %s", symbol)
	return nil
}

func (cdb *CryptoDB) UpdatePrice(symbol string, newPrice float64) error {
	res, err := cdb.conn.Exec(`
		UPDATE crypto 
		SET current_price = $1, last_update = NOW() 
		WHERE symbol = $2
	`, newPrice, symbol)
	
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	
	if rows == 0 {
		return ErrUnknownCoin
	}

	log.Printf("Updated price for %s: $%.2f", symbol, newPrice)
	return nil
}

func (cdb *CryptoDB) GetCount() (int, error) {
	var count int
	err := cdb.conn.QueryRow(`SELECT COUNT(*) FROM crypto`).Scan(&count)
	return count, err
}

func (cdb *CryptoDB) Close() error {
	if cdb.conn != nil {
		return cdb.conn.Close()
	}
	return nil
}

func (cdb *CryptoDB) Ping() error {
	return cdb.conn.Ping()
}