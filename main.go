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
	Quantity int    `json:"quantity"`
}

var (
	tickers      = []Stock{{"CMIG3.SA", 100}, {"TAEE3.SA", 100}, {"TRPL4.SA", 100}, {"PSSA3.SA", 100}, {"BBAS3.SA", 100}, {"SANB4.SA", 100}, {"KLBN4.SA", 100}}
	startDate    = time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	endDate      = time.Now().Format("2006-01-02")
	cpStartDate  = time.Now().AddDate(0, 0, -364).Format("2006-01-02")
	tickersCount = float64(len(tickers))
)

type InvestmentData struct {
	Ticker             string  `json:"ticker"`
	MarketValue 	   float64 `json:"market_value"`
	Quantity           int     `json:"quantity"`
	Apport    		   float64 `json:"apport"`
	CeilingPrice       float64 `json:"ceiling_price"`
	GrahamIndex        float64 `json:"graham_index"`
}

func calculateGrahamIndex(ticker string) float64 {
	cmd := exec.Command("python3", "calculate_graham_index.py", ticker)
	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Printf("Error running Python script for graham index of %s: %v\n", ticker, err)
		fmt.Printf("Python script output: %s\n", output)
		return 0
	}

	graham, err := strconv.ParseFloat(strings.TrimSpace(string(output)), 64)
	if err != nil {
		fmt.Printf("Error parsing Python script output for graham index of %s: %v\n", ticker, err)
		return 0
	}

	return graham
}

func calculateCeilingPrice(ticker, start, end string) float64 {
	cmd := exec.Command("python3", "calculate_ceiling_price.py", ticker, start, end)
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

func getQuote(ticker, start, end string) float64 {
	cmd := exec.Command("python3", "get_stock_data.py", ticker)
	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Printf("Error running Python script: %v\n", err)
		fmt.Printf("Python script output: %s\n", output)
		return 0
	}

	parts := strings.Split(strings.TrimSpace(string(output)), ",")
	value, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		fmt.Printf("Error parsing Python script output: %v\n", err)
		return 0
	}

	return value
}

func calculateInvestmentValue(ticker string, start, end string, availableBalance float64) InvestmentData {
	currentQuote := getQuote(ticker, start, end)
	ceilingPrice := calculateCeilingPrice(ticker, cpStartDate, end)
	grahamIndex := calculateGrahamIndex(ticker)
	avaiableForApport := availableBalance / float64(tickersCount)
	quantity := avaiableForApport / currentQuote

	return InvestmentData{
		Ticker:             ticker,
		MarketValue: 		currentQuote,
		Quantity:           int(quantity),
		Apport:    			avaiableForApport,
		CeilingPrice:       ceilingPrice,
		GrahamIndex:        grahamIndex,
	}
}

func getInvestmentData(w http.ResponseWriter, r *http.Request) {
	var investmentData []InvestmentData
	availableBalanceParam := r.URL.Query().Get("availableBalance")
	availableBalance, err := strconv.ParseFloat(availableBalanceParam, 64)
	if err != nil {
		http.Error(w, "Error parsing availableBalance: "+err.Error(), http.StatusBadRequest)
		return
	}
	
	for _, stock := range tickers {
		data := calculateInvestmentValue(stock.Ticker, startDate, endDate, availableBalance)
		investmentData = append(investmentData, data)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(investmentData)
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
