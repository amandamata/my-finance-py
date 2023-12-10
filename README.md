# My Finance Project

This project aims to provide financial information about a stock portfolio. Currently, the information is displayed in a table, showing the stock ticker, current market value, investment value, ceiling price, and Graham index.

## Project Setup

### Backend (Go)

Make sure you have Go installed on your machine.

1. Navigate to the project directory:

   ```bash
   cd path/to/your/project
   ```

2. Run the Go server:
   ```bash
   go run main.go
   ```

The server will be running at http://localhost:8080.

### Frontend (HTML, CSS, JavaScript)

Open the index.html file in a web browser. The frontend will make a request to the Go server to fetch financial data and display it in a table.

Keep in mind that, for production environments, it's important to configure CORS permissions properly on the server.

### Dependencies

Go (backend)
Python 3 (used by the script to fetch stock prices)
Go libraries: github.com/gorilla/mux and github.com/gorilla/handlers (installed automatically when running the server for the first time)
Make sure to have all dependencies installed before starting the project.

