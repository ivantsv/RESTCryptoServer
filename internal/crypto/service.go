package crypto

import (
	"RESTCryptoServer/internal/coingecko"
	"RESTCryptoServer/internal/db"
	"RESTCryptoServer/internal/redis"
	"fmt"
	"log"
	"math"
	"strings"
	"time"
	"errors"
)

var (
	ErrNameConflict error = errors.New("cryptocurrency already exists")
	ErrCryptoNotFound error = errors.New("cryptocurrency not found")
)

type CryptoResponse struct {
	Symbol       string    `json:"symbol"`
	Name         string    `json:"name"`
	CurrentPrice float64   `json:"current_price"`
	LastUpdated  time.Time `json:"last_updated"`
}

type CryptoResponseList struct {
	Cryptos []CryptoResponse `json:"cryptos"`
}

type CryptoHistoryResponse struct {
	Symbol  string                    `json:"symbol"`
	History []redis.PriceHistoryEntry `json:"history"`
}

type CryptoStats struct {
	MinPrice           float64 `json:"min_price"`
	MaxPrice           float64 `json:"max_price"`
	AvgPrice           float64 `json:"avg_price"`
	PriceChange        float64 `json:"price_change"`
	PriceChangePercent float64 `json:"price_change_percent"`
	RecordsCount       int     `json:"records_count"`
}

type CryptoStatsResponse struct {
	Symbol       string      `json:"symbol"`
	CurrentPrice float64     `json:"current_price"`
	Stats        CryptoStats `json:"stats"`
}

type CryptoService struct {
	cryptoDB    *db.CryptoDB
	redisClient *redis.RedisClient
}

func NewCryptoService(cryptoDB *db.CryptoDB, redisClient *redis.RedisClient) *CryptoService {
	return &CryptoService{
		cryptoDB:    cryptoDB,
		redisClient: redisClient,
	}
}

func (cs *CryptoService) AddCrypto(symbol string) (*CryptoResponse, error) {
	symbol = strings.ToLower(symbol)

	cnt, err := cs.redisClient.GetHistoryCount(symbol)
	if cnt != 0 {
		return nil, ErrNameConflict
	}
	if err != nil {
		return nil, fmt.Errorf("error during checking cryptocurrency: %s", err)
	}

	_, err = cs.cryptoDB.Get(symbol)
	if err == nil {
		return nil, ErrNameConflict
	}
	if err != db.ErrUnknownCoin {
		return nil, fmt.Errorf("database error: %w", err)
	}

	return cs.updateCoinPrice(symbol)
}

func (cs *CryptoService) GetAllCryptos() (*CryptoResponseList, error) {
	cryptos, err := cs.cryptoDB.GetAllSlice()
	if err != nil {
		return nil, fmt.Errorf("failed to get cryptocurrencies: %w", err)
	}
	
	response := &CryptoResponseList{
		Cryptos: make([]CryptoResponse, len(cryptos)),
	}
	
	for i, crypto := range cryptos {
		response.Cryptos[i] = CryptoResponse{
			Symbol:       crypto.Symbol,
			Name:         crypto.Name,
			CurrentPrice: crypto.CurrentPrice,
			LastUpdated:  crypto.LastUpdate,
		}
	}
	
	return response, nil
}

func (cs *CryptoService) GetCrypto(symbol string) (*CryptoResponse, error) {
	symbol = strings.ToLower(symbol)

	cnt, err := cs.redisClient.GetHistoryCount(symbol)
	if cnt == 0 {
		return nil, ErrCryptoNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("error during checking cryptocurrency: %s", err)
	}
	
	coinData, err := cs.cryptoDB.Get(symbol)
	if err != nil {
		if err == db.ErrUnknownCoin {
			return nil, ErrCryptoNotFound
		}
		return nil, fmt.Errorf("database error: %w", err)
	}
	
	return &CryptoResponse{
		Symbol:       symbol,
		Name:         coinData.Name,
		CurrentPrice: coinData.CurrentPrice,
		LastUpdated:  coinData.LastUpdate,
	}, nil
}

func (cs *CryptoService) RefreshCrypto(symbol string) (*CryptoResponse, error) {
	symbol = strings.ToLower(symbol)
	
	cnt, err := cs.redisClient.GetHistoryCount(symbol)
	if cnt == 0 {
		return nil, ErrCryptoNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("error during checking cryptocurrency: %s", err)
	}

	_, err = cs.cryptoDB.Get(symbol)
	if err != nil {
		if err == db.ErrUnknownCoin {
			return nil, ErrCryptoNotFound
		}
		return nil, fmt.Errorf("database error: %w", err)
	}
	
	return cs.updateCoinPrice(symbol)
}

func (cs *CryptoService) GetCryptoHistory(symbol string) (*CryptoHistoryResponse, error) {
	symbol = strings.ToLower(symbol)

	cnt, err := cs.redisClient.GetHistoryCount(symbol)
	if cnt == 0 {
		return nil, ErrCryptoNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("error during checking cryptocurrency: %s", err)
	}

	history, err := cs.redisClient.GetPriceHistory(symbol, 100)
	if err != nil {
		log.Printf("Warning: failed to get price history from Redis: %v", err)
		history = []redis.PriceHistoryEntry{}
	}
	
	return &CryptoHistoryResponse{
		Symbol:  symbol,
		History: history,
	}, nil
}

func (cs *CryptoService) GetCryptoStats(symbol string) (*CryptoStatsResponse, error) {
	symbol = strings.ToLower(symbol)
	
	cnt, err := cs.redisClient.GetHistoryCount(symbol)
	if cnt == 0 {
		return nil, ErrCryptoNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("error during checking cryptocurrency: %s", err)
	}

	coinData, err := cs.cryptoDB.Get(symbol)
	if err != nil {
		if err == db.ErrUnknownCoin {
			return nil, ErrCryptoNotFound
		}
		return nil, fmt.Errorf("database error: %w", err)
	}

	history, err := cs.redisClient.GetPriceHistory(symbol, 100)
	if err != nil {
		log.Printf("Warning: failed to get price history from Redis: %v", err)
		return &CryptoStatsResponse{
			Symbol:       symbol,
			CurrentPrice: coinData.CurrentPrice,
			Stats: CryptoStats{
				MinPrice:           coinData.CurrentPrice,
				MaxPrice:           coinData.CurrentPrice,
				AvgPrice:           coinData.CurrentPrice,
				PriceChange:        0,
				PriceChangePercent: 0,
				RecordsCount:       1,
			},
		}, nil
	}
	
	stats := cs.calculateStats(history, coinData.CurrentPrice)
	
	return &CryptoStatsResponse{
		Symbol:       symbol,
		CurrentPrice: coinData.CurrentPrice,
		Stats:        stats,
	}, nil
}

func (cs *CryptoService) DeleteCrypto(symbol string) error {
	symbol = strings.ToLower(symbol)
	
	err := cs.cryptoDB.Delete(symbol)
	if err != nil {
		if err == db.ErrUnknownCoin {
			return ErrCryptoNotFound
		}
		return fmt.Errorf("database error: %w", err)
	}

	err = cs.redisClient.DeletePriceHistory(symbol)
	if err != nil {
		log.Printf("Warning: failed to delete price history from Redis: %v", err)
	}
	
	return nil
}

func (cs *CryptoService) UpdateAllCryptos() (int, error) {
	cryptos, err := cs.cryptoDB.GetAllSlice()
	if err != nil {
		return 0, fmt.Errorf("failed to get cryptocurrencies: %w", err)
	}
	
	updated := 0
	for _, crypto := range cryptos {
		_, err := cs.updateCoinPrice(crypto.Symbol)
		if err != nil {
			log.Printf("Failed to update %s: %v", crypto.Symbol, err)
			continue
		}
		updated++
	}
	
	return updated, nil
}

func (cs *CryptoService) updateCoinPrice(symbol string) (*CryptoResponse, error) {
	coins, err := coingecko.GetCoinList()
	if err != nil {
		return nil, fmt.Errorf("failed to get coin list: %w", err)
	}

	ids, err := coingecko.GetIDBySymbol(symbol, coins)
	if err != nil {
		return nil, fmt.Errorf("cryptocurrency with symbol %s not found on CoinGecko", symbol)
	}

	if len(ids) == 0 {
		return nil, fmt.Errorf("cryptocurrency with symbol %s not found on CoinGecko", symbol)
	}

	coinID := ids[0]

	priceData, err := coingecko.GetPriceByID(coinID)
	if err != nil {
		return nil, fmt.Errorf("failed to get price for %s: %w", symbol, err)
	}

	price, exists := priceData["usd"]
	if !exists {
		return nil, fmt.Errorf("USD price not available for %s", symbol)
	}

	var coinName string
	for _, coin := range coins {
		if coin.ID == coinID {
			coinName = coin.Name
			break
		}
	}

	if coinName == "" {
		coinName = strings.ToUpper(symbol)
	}

	coinData := db.CoinData{
		Name:         coinName,
		CurrentPrice: price,
		LastUpdate:   time.Now().UTC(),
	}

	err = cs.cryptoDB.Insert(symbol, coinData)
	if err != nil {
		return nil, fmt.Errorf("failed to update in database: %w", err)
	}

	err = cs.redisClient.AddPriceHistory(symbol, price)
	if err != nil {
		log.Printf("Warning: failed to add to Redis history for %s: %v", symbol, err)
	}

	return &CryptoResponse{
		Symbol:       symbol,
		Name:         coinName,
		CurrentPrice: price,
		LastUpdated:  coinData.LastUpdate,
	}, nil
}

func (cs *CryptoService) calculateStats(history []redis.PriceHistoryEntry, currentPrice float64) CryptoStats {
	if len(history) == 0 {
		return CryptoStats{
			MinPrice:           currentPrice,
			MaxPrice:           currentPrice,
			AvgPrice:           currentPrice,
			PriceChange:        0,
			PriceChangePercent: 0,
			RecordsCount:       1,
		}
	}
	
	var minPrice, maxPrice, sum float64
	minPrice = history[0].Price
	maxPrice = history[0].Price
	
	for _, entry := range history {
		if entry.Price < minPrice {
			minPrice = entry.Price
		}
		if entry.Price > maxPrice {
			maxPrice = entry.Price
		}
		sum += entry.Price
	}
	
	avgPrice := sum / float64(len(history))

	var priceChange, priceChangePercent float64
	if len(history) > 1 {
		oldestPrice := history[len(history)-1].Price
		priceChange = currentPrice - oldestPrice
		if oldestPrice > 0 {
			priceChangePercent = (priceChange / oldestPrice) * 100
		}
	}

	return CryptoStats{
		MinPrice:           math.Round(minPrice*100) / 100,
		MaxPrice:           math.Round(maxPrice*100) / 100,
		AvgPrice:           math.Round(avgPrice*100) / 100,
		PriceChange:        math.Round(priceChange*100) / 100,
		PriceChangePercent: math.Round(priceChangePercent*100) / 100,
		RecordsCount:       len(history),
	}
}