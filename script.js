document.addEventListener("DOMContentLoaded", function () {
    const availableBalanceInput = document.getElementById("availableBalance");
    const updateButton = document.getElementById("updateButton");

    updateButton.addEventListener("click", function () {
        const availableBalance = availableBalanceInput.value;
        
        fetchDataFromBackend(availableBalance).then(data => {
            renderFinancialData(data);
        });
    });
});

async function fetchDataFromBackend(availableBalance) {
    const url = new URL('http://localhost:8080/api/investment_data');
    url.searchParams.append('availableBalance', availableBalance);

    return fetch(url)
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
        ["Ticker", "Market Value", "Ceiling Price", "Graham Index", "Quantity", "Investment"],
        data.map(stock => [
            stock.ticker,
            stock.current_market_value,
            stock.ceiling_price,
            stock.graham_index,
            stock.quantity,
            stock.investment_value
        ])
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

async function updateValues() {
    const availableBalance = parseFloat(document.getElementById("availableBalance").value) || 0;
    const quantityPerStock = parseInt(document.getElementById("quantityPerStock").value) || 0;

    const response = await fetch(`http://localhost:8080/api/investment_data?availableBalance=${availableBalance}&quantityPerStock=${quantityPerStock}`);
    const data = await response.json();

    renderFinancialData(data);
}
