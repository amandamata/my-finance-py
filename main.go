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
	tickers    = []Stock{{"CMIG3.SA", 100}, {"TAEE3.SA", 100}, {"TRPL4.SA", 100}, {"PSSA3.SA", 100}, {"BBAS3.SA", 100}, {"SANB4.SA", 100}, {"KLBN4.SA", 100}}
	startDate   = time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	endDate     = time.Now().Format("2006-01-02")
	cpStartDate = time.Now().AddDate(0, 0, -364).Format("2006-01-02")
)

type InvestmentData struct {
	Ticker             string  `json:"ticker"`
	CurrentMarketValue float64 `json:"current_market_value"`
	InvestmentValue    float64 `json:"investment_value"`
	CeilingPrice       float64 `json:"ceiling_price"`
	GrahamIndex        float64 `json:"graham_index"`
	Quantity           int     `json:"quantity"`
}

func calculateGrahamIndex(eps, bookValue, closePrice float64) float64 {
	return 22.5 * eps * bookValue / closePrice
}

func calculateCeilingPrice(closePrices []float64) float64 {
	total := 0.0
	for _, price := range closePrices {
		total += price
	}
	average := total / float64(len(closePrices))
	return average
}

func getQuote(ticker, start, end string) float64 {
	cmd := exec.Command("python3", "get_stock_price.py", ticker)
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

func getRealValues(ticker string) (float64, float64) {
	cmd := exec.Command("python3", "get_stock_price.py", ticker)
	output, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Printf("Error running Python script for %s: %v\n", ticker, err)
		return 0, 0
	}

	parts := strings.Split(strings.TrimSpace(string(output)), ",")
	eps, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		fmt.Printf("Error converting EPS for %s: %v\n", ticker, err)
		return 0, 0
	}

	bookValue, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		fmt.Printf("Error converting BookValue for %s: %v\n", ticker, err)
		return 0, 0
	}

	return eps, bookValue
}

func calculateInvestmentValue(ticker string, quantity int, start, end string, totalValue float64) InvestmentData {
	currentQuote := getQuote(ticker, start, end)
	eps, bookValue := getRealValues(ticker)
	investmentValue := totalValue / float64(len(tickers))

	ceilingPrice := calculateCeilingPrice([]float64{currentQuote})
	grahamIndex := calculateGrahamIndex(eps, bookValue, currentQuote)

	return InvestmentData{
		Ticker:             ticker,
		CurrentMarketValue: currentQuote,
		InvestmentValue:    investmentValue,
		CeilingPrice:       ceilingPrice,
		GrahamIndex:        grahamIndex,
		Quantity:           quantity,
	}
}

func getInvestmentData(w http.ResponseWriter, r *http.Request) {
	var investmentData []InvestmentData

	var totalValue float64
	for _, stock := range tickers {
		currentQuote := getQuote(stock.Ticker, startDate, endDate)
		totalValue += float64(stock.Quantity) * currentQuote
	}

	for _, stock := range tickers {
		data := calculateInvestmentValue(stock.Ticker, stock.Quantity, startDate, endDate, totalValue)
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
