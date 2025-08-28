package crypto

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type SymbolJSON struct {
	Symbol string `json:"symbol"`
}

func GETCryptosHandler(cs *CryptoService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cryptos, err := cs.GetAllCryptos()
		if err != nil {
			log.Println("error during getting all cryptos: ", err)
			http.Error(w, `Server error`, http.StatusInternalServerError)	
			return
		}
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		
		json.NewEncoder(w).Encode(cryptos)
	}
}

func POSTCryptoHandler(cs *CryptoService) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		var symbolJSON SymbolJSON
		
		err := json.NewDecoder(r.Body).Decode(&symbolJSON)
		if err != nil {
			log.Println("Error during JSON parsing: ", err)
			http.Error(w, `Bad Request`, http.StatusBadRequest)
			return
		}

		resp, err := cs.AddCrypto(symbolJSON.Symbol)
		if err == ErrNameConflict {
			log.Println(err)
			http.Error(w, `Name conflict`, http.StatusConflict)
			return
		}
		if err != nil {
			log.Println("adding cryptocurrency error: ", err)
			http.Error(w, `Server error`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(resp)
	}
}

func GETCryptoSymbolHandler(cs *CryptoService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		symbol := chi.URLParam(r, "symbol")

		resp, err := cs.GetCrypto(symbol)
		if err != nil {
			log.Println("error during getting crypto from postgres: ", err)
			http.Error(w, `Bad Request`, http.StatusBadRequest)
			return
		}
	
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(resp)
	}
}

func PUTCryptoSymbolRefreshHandler(cs *CryptoService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		symbol := chi.URLParam(r, "symbol")

		resp, err := cs.RefreshCrypto(symbol)
		if err == ErrCryptoNotFound {
			log.Println(err)
			http.Error(w, `Bad Request`, http.StatusBadRequest)
			return
		}
		if err != nil {
			log.Println("updating cryptocurrency error: ", err)
			http.Error(w, `Server error`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(resp)
	}
}

func GETCryptoHistoryHandler(cs *CryptoService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		symbol := chi.URLParam(r, "symbol")

		resp, _ := cs.GetCryptoHistory(symbol)
		
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		
		json.NewEncoder(w).Encode(resp)
	}
}

func GETCryptoStatsHandler(cs *CryptoService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		symbol := chi.URLParam(r, "symbol")

		resp, _ := cs.GetCryptoStats(symbol)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		
		json.NewEncoder(w).Encode(resp)
	}
}

func DELETECryptoSymbolHandler(cs *CryptoService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		symbol := chi.URLParam(r, "symbol")

		err := cs.DeleteCrypto(symbol)
		if err != nil {
			log.Println(err)
			http.Error(w, `Bad Request`, http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{}`))
	}
}