# Investment Analysis Script

This script is designed to analyze the performance of a portfolio containing various stocks. It leverages the Yahoo Finance API to obtain historical stock data and calculate relevant metrics. The primary functionalities include:

- Retrieving historical stock data for a list of specified tickers.
- Calculating the market value and investment value of each stock in the portfolio.
- Estimating the ceiling price for each stock based on the average price-to-earnings ratio.

## Usage

### Libraries Required
Make sure you have the necessary Python libraries installed:

```bash
pip install numpy pandas yfinance
```

### Script Configuration
Adjust the script variables to fit your portfolio:
```bash
tickers = ['CMIG3.SA', 'TAEE3.SA', 'TRPL4.SA', 'PSSA3.SA', 'BBAS3.SA', 'SANB4.SA', 'KLBN4.SA']
quantities = [100, 100, 100, 100, 100, 100, 100]
```
### Run the Script
Execute the script in your preferred Python environment:
```bash
python investment_analysis.py
```

### Output
The script will provide detailed information for each stock in the portfolio, including the market quote, investment value, estimated ceiling price, and more.

### Note
This script is a basic tool and may require further customization based on your specific needs. Additionally, ensure you have an internet connection to fetch the latest stock data.

Feel free to explore and enhance the script according to your investment analysis requirements. If you encounter any issues or have suggestions for improvement, please let me know. Happy investing!