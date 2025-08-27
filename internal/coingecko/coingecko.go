package coingecko

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type CoinInfo struct {
	ID     string `json:"id"`
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
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

func GetIDBySymbol(symbol string, coins []CoinInfo) ([]string, error) {
	var answer []string
	symbol = strings.ToLower(symbol)
	for _, coin := range coins {
		if coin.Symbol == symbol {
			answer = append(answer, coin.ID)
		}
	}

	if (len(answer) == 0) {
		return answer, fmt.Errorf("coin with symbol %s not found", symbol)
	}
	return answer, nil
}

func GetPriceByID(id string) (map[string]float64, error) {
	url := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=usd", id)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("error during getting coin: ", err)
		return nil, err
	}
	defer resp.Body.Close()

	var result map[string]map[string]float64
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("error during getting coin: ", err)
		return nil, err
	}
	return result[id], nil
}