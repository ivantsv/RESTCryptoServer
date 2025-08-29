package coingecko

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type CoinInfo struct {
	ID     string `json:"id"`
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}

var popularCryptoMap = map[string]string{
	"btc":  "bitcoin",
	"eth":  "ethereum", 
	"usdt": "tether",
	"bnb":  "binancecoin",
	"sol":  "solana",
	"usdc": "usd-coin",
	"xrp":  "ripple",
	"doge": "dogecoin",
	"ada":  "cardano",
	"trx":  "tron",
	"avax": "avalanche-2",
	"shib": "shiba-inu",
	"dot":  "polkadot",
	"link": "chainlink",
	"ltc":  "litecoin",
}

func GetCoinList() ([]CoinInfo, error) {
	resp, err := http.Get("https://api.coingecko.com/api/v3/coins/list")
	if err != nil {
		log.Println("error during getting coin list: ", err)
		return nil, err
	}
	defer resp.Body.Close()

	var coins []CoinInfo
	if err := json.NewDecoder(resp.Body).Decode(&coins); err != nil {
		return nil, err
	}
	return coins, nil
}

func GetIDBySymbol(symbol string, coins []CoinInfo) (string, error) {
	symbol = strings.ToLower(symbol)

	if popularID, exists := popularCryptoMap[symbol]; exists {
		log.Printf("Using predefined mapping for %s -> %s", symbol, popularID)
		return popularID, nil
	}

	var matches []CoinInfo
	for _, coin := range coins {
		if coin.Symbol == symbol {
			matches = append(matches, coin)
		}
	}

	if len(matches) == 0 {
		return "", fmt.Errorf("coin with symbol %s not found", symbol)
	}

	if len(matches) > 1 {
		log.Printf("Found %d coins with symbol %s:", len(matches), symbol)
		for _, match := range matches {
			log.Printf("  - ID: %s, Name: %s", match.ID, match.Name)
		}

		for _, match := range matches {
			lowerName := strings.ToLower(match.Name)
			if symbol == "btc" && strings.Contains(lowerName, "bitcoin") {
				log.Printf("Selected Bitcoin: %s", match.ID)
				return match.ID, nil
			}
			if symbol == "eth" && strings.Contains(lowerName, "ethereum") {
				log.Printf("Selected Ethereum: %s", match.ID)
				return match.ID, nil
			}
		}
	}

	log.Printf("Using first match for %s: %s (%s)", symbol, matches[0].ID, matches[0].Name)
	return matches[0].ID, nil
}

func parsePrice(value any) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case string:
		return strconv.ParseFloat(v, 64)
	case int:
		return float64(v), nil
	case int64:
		return float64(v), nil
	default:
		return 0, fmt.Errorf("cannot convert %T to float64", value)
	}
}

func GetPriceByID(id string) (map[string]float64, error) {
	url := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=usd", id)
	
	log.Printf("Fetching price from: %s", url)
	
	resp, err := http.Get(url)
	if err != nil {
		log.Println("error during getting coin: ", err)
		return nil, err
	}
	defer resp.Body.Close()

	var rawResult map[string]map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&rawResult); err != nil {
		log.Println("error during decoding response: ", err)
		return nil, err
	}
	
	log.Printf("Raw API response: %+v", rawResult)

	coinData, exists := rawResult[id]
	if !exists {
		return nil, fmt.Errorf("no price data found for ID: %s", id)
	}

	result := make(map[string]float64)
	for currency, priceValue := range coinData {
		price, err := parsePrice(priceValue)
		if err != nil {
			log.Printf("Warning: could not parse price for %s/%s: %v", id, currency, err)
			continue
		}
		result[currency] = price
	}
	
	if len(result) == 0 {
		return nil, fmt.Errorf("no valid price data found for %s", id)
	}
	
	log.Printf("Parsed prices for %s: %+v", id, result)
	return result, nil
}