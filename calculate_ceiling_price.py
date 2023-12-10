import yfinance as yf
import sys

def calculate_ceiling_price(ticker, start, end):
    try:
        stock_data = yf.Ticker(ticker)
        dividends = stock_data.dividends.loc[start:end]
        total_dividends = dividends.sum()

        dividend_yield = 0.06
        ceiling_price = total_dividends / dividend_yield

        return ceiling_price
    except Exception as e:
        print(f"Error calculating ceiling price for {ticker}: {e}")
        return None

if __name__ == "__main__":
    if len(sys.argv) != 4:
        print("Usage: python calculate_ceiling_price.py <ticker> <start_date> <end_date>")
        sys.exit(1)

    ticker = sys.argv[1]
    start_date = sys.argv[2]
    end_date = sys.argv[3]

    result = calculate_ceiling_price(ticker, start_date, end_date)
    if result is not None:
        print(result)
    else:
        sys.exit(1)
