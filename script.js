document.addEventListener("DOMContentLoaded", function () {
    fetchDataFromBackend().then(data => {
        renderFinancialData(data);
    });
});

async function fetchDataFromBackend() {
    return fetch('http://localhost:8080/api/investment_data')
        .then(response => response.json())
        .then(data => {
            return data;
        })
        .catch(error => {
            console.error('Error fetching data from backend:', error);
            return null;
        });
}

function renderFinancialData(data) {
    const appDiv = document.getElementById("app");
    const title = document.createElement("h2");
    title.textContent = "Brazilian Stocks Apports";
    appDiv.appendChild(title);

    const stocksDataTable = createTable(
        ["Ticker", "Market Value", "Ceiling Price", "Graham Index"],
        data.map(stock => [stock.ticker, stock.current_market_value, stock.ceiling_price, stock.graham_index])
    );
    appDiv.appendChild(stocksDataTable);
}

function createTable(headers, rows) {
    const table = document.createElement("table");

    const headerRow = document.createElement("tr");
    headers.forEach(headerText => {
        const th = document.createElement("th");
        th.textContent = headerText;
        headerRow.appendChild(th);
    });
    table.appendChild(headerRow);

    rows.forEach(rowData => {
        const tr = document.createElement("tr");
        rowData.forEach(cellData => {
            const td = document.createElement("td");
            td.textContent = formatTableCell(cellData);
            tr.appendChild(td);
        });
        table.appendChild(tr);
    });

    return table;
}

function formatTableCell(data) {
    if (typeof data === "number") {
        return Number.isInteger(data) ? data.toFixed(0) : data.toFixed(2);
    }
    return data;
}
