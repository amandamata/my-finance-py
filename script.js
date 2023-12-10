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
    const tableBody = document.querySelector("#investmentTable tbody");

    tableBody.innerHTML = "";

    data.forEach(stock => {
        const row = tableBody.insertRow();
        const cells = [
            stock.ticker,
            stock.market_value.toFixed(2),
            stock.quantity,
            stock.apport.toFixed(2),
            stock.ceiling_price.toFixed(2),
            stock.graham_index.toFixed(2),
        ];

        cells.forEach((cell, index) => {
            const td = row.insertCell(index);
            td.textContent = formatTableCell(cell);
        });
    });
}

function formatTableCell(data) {
    if (typeof data === "number") {
        return Number.isInteger(data) ? data.toFixed(0) : data.toFixed(2);
    }
    return data;
}
