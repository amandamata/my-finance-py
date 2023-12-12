package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

type Stock struct {
	Ticker   string `json:"ticker"`
}

var (
	startDate          = time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	endDate            = time.Now().Format("2006-01-02")
	cpStartDate        = time.Now().AddDate(0, 0, -364).Format("2006-01-02")
	stockDataCache     = make(map[string]map[string]float64)
	financialsFetched  = false
	financialsFetchDate time.Time
)

type InvestmentData struct {
	Ticker        string  `json:"ticker"`
	MarketValue   float64 `json:"marketValue"`
	CeilingPrice  float64 `json:"ceilingPrice"`
	GrahamIndex   float64 `json:"grahamIndex"`
	Quantity      int     `json:"quantity"`
	Invest        float64 `json:"invest"`
}

func getCurrentPrice(ticker string) float64 {
	cmd := exec.Command("python3", "python/get_current_price.py", ticker, startDate, endDate)
	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Printf("Error running Python script for current price: %v\n", err)
		fmt.Printf("Python script output: %s\n", output)
		return 0
	}

	currentPrice, err := strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
	if err != nil {
		fmt.Printf("Error parsing Python script output for current price: %v\n", err)
		return 0
	}

	return currentPrice
}

func calculateGrahamIndex(ticker string) float64 {
	cmd := exec.Command("python3", "python/calculate_graham_index.py", ticker)
	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Printf("Error running Python script: %v\n", err)
		fmt.Printf("Python script output: %s\n", output)
		return 0
	}

	graham, err := strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
	if err != nil {
		fmt.Printf("Error parsing Python script output: %v\n", err)
		return 0
	}

	return graham
}

func calculateCeilingPrice(ticker, start, end string) float64 {
	cmd := exec.Command("python3", "python/calculate_ceiling_price.py", ticker, start, end)
	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Printf("Error running Python script for ceiling price of %s: %v\n", ticker, err)
		fmt.Printf("Python script output: %s\n", output)
		return 0
	}

	ceilingPrice, err := strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
	if err != nil {
		fmt.Printf("Error parsing Python script output for ceiling price of %s: %v\n", ticker, err)
		return 0
	}

	return ceilingPrice
}

func calculateInvestmentValue(selectedStocks []string, availableBalance float64) []InvestmentData {
	var investmentData []InvestmentData

	for _, stock := range selectedStocks {
		ceilingPrice := calculateCeilingPrice(stock, cpStartDate, endDate)
		availableForInvest := availableBalance / float64(len(selectedStocks))
		graham := calculateGrahamIndex(stock)
		currentQuote := getCurrentPrice(stock)
		quantity := availableForInvest / currentQuote

		data := InvestmentData{
			Ticker:       stock,
			MarketValue:  currentQuote,
			CeilingPrice: ceilingPrice,
			GrahamIndex:  graham,
			Quantity:     int(quantity),
			Invest:       availableForInvest,
		}
		investmentData = append(investmentData, data)
	}

	return investmentData
}

func getInvestmentData(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Processing getInvestmentData...")

    selectedStocksParam := r.URL.Query().Get("selectedStocks")
    selectedStocks := strings.Split(selectedStocksParam, ",")
    availableBalanceParam := r.URL.Query().Get("availableBalance")
    availableBalance, err := strconv.ParseFloat(availableBalanceParam, 64)
    if err != nil {
        http.Error(w, "Error parsing availableBalance: "+err.Error(), http.StatusBadRequest)
        return
    }

    data := calculateInvestmentValue(selectedStocks, availableBalance)

    w.Header().Set("Content-Type", "application/json")
    if err := json.NewEncoder(w).Encode(data); err != nil {
        fmt.Printf("Error encoding JSON response: %v\n", err)
        http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
    }
}


func main() {
	r := mux.NewRouter()
	r.HandleFunc("/api/investment_data", getInvestmentData).Methods("GET")

	fmt.Println("Server is running on :8080")

	headers := handlers.AllowedHeaders([]string{"Content-Type"})
	origins := handlers.AllowedOrigins([]string{"*"})
	methods := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	http.Handle("/", handlers.CORS(headers, origins, methods)(r))

	http.ListenAndServe(":8080", nil)
}
