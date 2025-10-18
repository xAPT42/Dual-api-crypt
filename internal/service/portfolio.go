package service

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Portfolio représente un portefeuille crypto
type Portfolio struct {
	Assets map[string]float64 `json:"assets"`
}

// PriceData représente les données de prix d'une crypto
type PriceData struct {
	ID                           string  `json:"id"`
	Symbol                       string  `json:"symbol"`
	Name                         string  `json:"name"`
	CurrentPrice                 float64 `json:"current_price"`
	MarketCap                    float64 `json:"market_cap"`
	TotalVolume                  float64 `json:"total_volume"`
	High24h                      float64 `json:"high_24h"`
	Low24h                       float64 `json:"low_24h"`
	PriceChange24h               float64 `json:"price_change_24h"`
	PriceChangePercentage24h     float64 `json:"price_change_percentage_24h"`
	MarketCapChange24h           float64 `json:"market_cap_change_24h"`
	MarketCapChangePercentage24h float64 `json:"market_cap_change_percentage_24h"`
	LastUpdated                  string  `json:"last_updated"`
}

// PortfolioValue représente la valeur calculée du portefeuille
type PortfolioValue struct {
	Value     float64   `json:"value"`
	Currency  string    `json:"currency"`
	Timestamp time.Time `json:"timestamp"`
}

// GetDefaultPortfolio retourne le portefeuille par défaut
func GetDefaultPortfolio() Portfolio {
	return Portfolio{
		Assets: map[string]float64{
			"bitcoin":    0.5,  // 0.5 BTC
			"ethereum":   5.0,  // 5.0 ETH
			"cardano":    1000, // 1000 ADA
			"solana":     50,   // 50 SOL
			"chainlink":  100,  // 100 LINK
		},
	}
}

// CalculatePortfolioValue calcule la valeur totale du portefeuille
func CalculatePortfolioValue() (float64, string, time.Time, error) {
	portfolio := GetDefaultPortfolio()
	
	// Construire la liste des IDs des cryptos pour l'API
	var cryptoIDs []string
	for id := range portfolio.Assets {
		cryptoIDs = append(cryptoIDs, id)
	}
	
	// Construire l'URL de l'API CoinGecko
	apiURL := fmt.Sprintf("https://api.coingecko.com/api/v3/simple/price?ids=%s&vs_currencies=eur", 
		joinStrings(cryptoIDs, ","))
	
	// Faire la requête HTTP
	resp, err := http.Get(apiURL)
	if err != nil {
		return 0, "", time.Time{}, fmt.Errorf("erreur lors de la requête API: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return 0, "", time.Time{}, fmt.Errorf("API retourne un code d'erreur: %d", resp.StatusCode)
	}
	
	// Lire la réponse
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, "", time.Time{}, fmt.Errorf("erreur lors de la lecture de la réponse: %v", err)
	}
	
	// Parser la réponse JSON
	var priceData map[string]map[string]float64
	if err := json.Unmarshal(body, &priceData); err != nil {
		return 0, "", time.Time{}, fmt.Errorf("erreur lors du parsing JSON: %v", err)
	}
	
	// Calculer la valeur totale
	totalValue := 0.0
	for cryptoID, quantity := range portfolio.Assets {
		if priceInfo, exists := priceData[cryptoID]; exists {
			if eurPrice, exists := priceInfo["eur"]; exists {
				totalValue += quantity * eurPrice
			}
		}
	}
	
	return totalValue, "EUR", time.Now(), nil
}

// joinStrings joint une liste de strings avec un séparateur
func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	if len(strs) == 1 {
		return strs[0]
	}
	
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}
