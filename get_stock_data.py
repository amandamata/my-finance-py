import yfinance as yf
import sys

cache = {}

def get_stock_data(ticker):
    key = f"{ticker}"

    if key in cache:
        return cache[key]

    try:
        stock_data = yf.Ticker(ticker)
        eps = stock_data.info.get('trailingEps', None)
        book_value = stock_data.info.get('bookValue', None)

        try:
            historical_data = stock_data.history(period='1d')
            close_price = historical_data['Close'].mean()
            adj_close = historical_data['Adj Close'].mean()
        except KeyError:
            close_price = stock_data.history(period='1d')['Close'].mean()
            adj_close = close_price

        result = f"{eps},{book_value},{close_price},{adj_close}"

        cache[key] = result

        return result
    except Exception as e:
        print(f"Error obtaining data for {ticker}: {e}", file=sys.stderr)
        return None

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python get_stock_data.py <ticker>")
        sys.exit(1)

    ticker = sys.argv[1]
    result = get_stock_data(ticker)
    if result:
        print(result)
    else:
        sys.exit(1)
