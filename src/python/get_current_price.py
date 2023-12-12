import yfinance as yf
import sys

cache = {}

def get_current_price(ticker):
    key = f"{ticker}"

    if key in cache:
        return cache[key]

    ticker = yf.Ticker(ticker)
    current_price = ticker.history(period='1d')['Close'].iloc[-1]
    cache[key] = current_price
    return current_price

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python get_current_price.py <ticker>")
        sys.exit(1)

    ticker = sys.argv[1]

    try:
        result = get_current_price(ticker)
        if result is not None:
            print(result)
        else:
            print("Error: Result is None")
            sys.exit(1)
    except Exception as e:
        print(f"Error running Python script for current price: {e}")
        sys.exit(1)
