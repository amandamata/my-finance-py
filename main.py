# import yfinance as yf
# from datetime import datetime, timedelta


# for ticker in tickers:
#     try:
#         stock_data = yf.download(ticker, start=start_date, end=end_date)
#         print(f"Data for {ticker} successfully obtained:")
#         print(stock_data.head())
#         print("\n")
#     except Exception as e:
#         print(f"Error obtaining data for {ticker}: {e}\n")
import yfinance as yf
from datetime import datetime, timedelta

tickers = ['CMIG3.SA', 'TAEE3.SA', 'TRPL4.SA', 'PSSA3.SA', 'BBAS3.SA', 'SANB4.SA', 'KLBN4.SA']
start_date = (datetime.now() - timedelta(days=1)).strftime('%Y-%m-%d')
end_date = datetime.now().strftime('%Y-%m-%d')

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

tickers = ['CMIG3.SA', 'TAEE3.SA', 'TRPL4.SA', 'PSSA3.SA', 'BBAS3.SA', 'SANB4.SA', 'KLBN4.SA']
end_date = datetime.now().strftime('%Y-%m-%d')

quantities = [12, 12, 12, 12, 12, 12, 12, 12]
total_current = 0
total_previous = 0

for i, ticker in enumerate(tickers):
    quantity = quantities[i]

    current_value, current_quote = calculate_investment_value(ticker, quantity, start_date, end_date)
    total_current += current_value

    previous_value, previous_quote = calculate_investment_value(ticker, quantity, get_previous_month_start_date(), get_previous_month_end_date())
    total_previous += previous_value

    print(f"Ticker: {ticker}")
    print(f"Quantity: {quantity}")

    if current_value is not None:
        formatted_value = "R$ {:.2f}".format(current_value)
        print(f"Investment value (Current): {formatted_value}")
        print("Quote: R${:.2f}".format(current_quote))
    else:
        print("Unable to obtain the quote for this ticker (Current).")

    if previous_value is not None:
        formatted_previous_value = "R$ {:.2f}".format(previous_value)
        print(f"Investment value (Previous): {formatted_previous_value}")
        print("Quote: R${:.2f}".format(previous_quote))
    else:
        print("Unable to obtain the quote for this ticker (Previous).")

    print("\n")

print("Total (Current): R${:.2f}".format(total_current))
print("Total (Previous): R${:.2f}".format(total_previous))
