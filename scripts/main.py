
import numpy as np
import pandas as pd
import yfinance as yf
from datetime import datetime, timedelta

tickers =    ['CMIG3.SA', 'TAEE3.SA', 'TRPL4.SA', 'PSSA3.SA', 'BBAS3.SA', 'SANB4.SA', 'KLBN4.SA']
quantities = [ 100,        100,        100,        100,        100,        100,        100]

start_date = (datetime.now() - timedelta(days=1)).strftime('%Y-%m-%d')
end_date = datetime.now().strftime('%Y-%m-%d')
cp_start_date = (datetime.now() - timedelta(days=364)).strftime('%Y-%m-%d')

def calculate_graham_index(ticker):
    try:
        stock_data = yf.Ticker(ticker)
        earnings_per_share = stock_data.info['trailingEps']
        book_value_per_share = stock_data.info['bookValue']
        graham_index = np.sqrt(22.5 * earnings_per_share * book_value_per_share)
        return graham_index

    except Exception as e:
        print(f"Error calculating Graham Index for {ticker}: {e}")
        return None

def calculate_ceiling_price(ticker):
    data = yf.download(ticker, start=cp_start_date, end=datetime.now().strftime('%Y-%m-%d'))
    pe_ratio_average = data['Close'].mean() / data['Adj Close'].mean()
    ceiling_price = pe_ratio_average * data['Adj Close'].iloc[-1]
    return ceiling_price

def get_quote(ticker, start, end):
    try:
        stock_data = yf.download(ticker, start=start, end=end)
        return stock_data['Close'].iloc[0]
    except Exception as e:
        print(f"Error obtaining quote for {ticker}: {e}")
        return None

def calculate_investment_value(ticker, quantity, start, end):
    quote = get_quote(ticker, start, end)
    
    if quote is None:
        return None

    investment_value = quantity * quote
    return investment_value, quote

def get_previous_month_start_date():
    return (datetime.strptime(start_date, '%Y-%m-%d') - timedelta(days=30)).strftime('%Y-%m-%d')

def get_previous_month_end_date():
    return (datetime.strptime(end_date, '%Y-%m-%d') - timedelta(days=30)).strftime('%Y-%m-%d')

def get_current_market_value(ticker):
    try:
        ticker_data = yf.Ticker(ticker)
        current_price = ticker_data.history(period='1d')['Close'].iloc[0]
        return current_price
    except Exception as e:
        print(f"Error obtaining current market value for {ticker}: {e}")
        return None

total_current = 0
total_previous = 0

for i, ticker in enumerate(tickers):
    quantity = quantities[i]

    current_value, current_quote = calculate_investment_value(ticker, quantity, start_date, end_date)
    total_current += current_value
    current_market_value = get_current_market_value(ticker)

    previous_value, previous_quote = calculate_investment_value(ticker, quantity, get_previous_month_start_date(), get_previous_month_end_date())
    total_previous += previous_value
    ceiling_price_calculated = calculate_ceiling_price(ticker)

    graham_index_calculated = calculate_graham_index(ticker)

    print(f"Ticker: {ticker}")
    print(f"Quantity: {quantity}")

    if current_value is not None:
        formatted_value = "R$ {:.2f}".format(current_value)
        print("Market Quote: R${:.2f}".format(current_market_value))
        print(f"Investment value: {formatted_value}")
        print(f'The estimated ceiling price for {ticker} is: {ceiling_price_calculated:.2f}')
        print(f'The calculated Graham Index for {ticker} is: {graham_index_calculated:.2f}')
    else:
        print("Unable to obtain the quote for this ticker (Current).")

    if previous_value is not None:
        formatted_previous_value = "R$ {:.2f}".format(previous_value)
        print("Previous Market Quote: R${:.2f}".format(previous_quote))
        print(f"Previous Investment value: {formatted_previous_value}")
    else:
        print("Unable to obtain the quote for this ticker (Previous).")

    print("\n")

print("Total (Current): R${:.2f}".format(total_current))
