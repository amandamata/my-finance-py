import yfinance as yf
import sys
import math

cache = {}

def calculate_graham_index(ticker):
    key = f"{ticker}"

    if key in cache:
        return cache[key]

    info = yf.Ticker(ticker).info    
    eps = info.get('trailingEps', None)
    bvps = info.get('bookValue', None)
    
    if eps is None or bvps is None:
        return None
    
    intrinsic_value = math.sqrt(22.5 * eps * bvps)

    cache[key] = intrinsic_value

    return intrinsic_value

if __name__ == "__main__":
    if len(sys.argv) != 2:
        print("Usage: python calculate_graham_index.py <ticker>")
        sys.exit(1)

    ticker = sys.argv[1]

    result = calculate_graham_index(ticker)
    if result is not None:
        print(result)
    else:
        sys.exit(1)
