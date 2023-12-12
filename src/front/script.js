document.addEventListener("DOMContentLoaded", async function () {
    const availableBalanceInput = document.getElementById("availableBalance");
    const updateButton = document.getElementById("updateButton");
    const appDiv = document.getElementById("app");

    await renderAvailableStocksSelect();

    updateButton.addEventListener("click", async function () {
        const availableBalance = availableBalanceInput.value;
        const selectedStocks = getSelectedStocks();
        const data = await fetchDataFromBackend(selectedStocks, availableBalance);

        clearAppDiv();
        renderFinancialData(data);
    });

    async function fetchDataFromBackend(selectedStocks, availableBalance) {
        const url = new URL('http://localhost:8080/api/investment_data');
        url.searchParams.append('selectedStocks', selectedStocks.join(','));
        url.searchParams.append('availableBalance', availableBalance);

        try {
            const response = await fetch(url);

            if (response.ok) {
                return await response.json();
            } else {
                console.error('Error fetching data from backend:', response.statusText);
                return null;
            }
        } catch (error) {
            console.error('Error fetching data from backend:', error);
            return null;
        }
    }

    async function renderAvailableStocksSelect() {
        const data = await fetchAvailableStocks();

        const stocksContainer = document.getElementById("stocksContainer");

        if (!Array.isArray(data)) {
            console.error('Invalid data for rendering stocks select:', data);
            return;
        }

        data.forEach(stock => {
            const checkbox = document.createElement("input");
            checkbox.type = "checkbox";
            checkbox.value = stock;
            checkbox.id = stock;
            checkbox.className = "stockCheckbox";

            const label = document.createElement("label");
            label.htmlFor = stock;
            label.appendChild(document.createTextNode(stock));

            stocksContainer.appendChild(checkbox);
            stocksContainer.appendChild(label);
            stocksContainer.appendChild(document.createElement("br"));
        });
    }

    async function fetchAvailableStocks() {
        try {
            const response = await fetch('http://localhost:8080/api/available_stocks');

            if (response.ok) {
                return await response.json() || [];
            } else {
                console.error('Error fetching available stocks:', response.statusText);
                return [];
            }
        } catch (error) {
            console.error('Error fetching available stocks:', error);
            return [];
        }
    }

    function clearAppDiv() {
        while (appDiv.firstChild) {
            appDiv.removeChild(appDiv.firstChild);
        }
    }

    function getSelectedStocks() {
        const checkboxes = document.getElementsByClassName("stockCheckbox");
        return Array.from(checkboxes).filter(checkbox => checkbox.checked).map(checkbox => checkbox.value);
    }

    function formatCurrency(value) {
        return new Intl.NumberFormat('pt-BR', { style: 'currency', currency: 'BRL' }).format(value);
    }

    function renderFinancialData(data) {
        clearAppDiv();

        const table = document.createElement("table");
        table.border = "1";

        const headerRow = table.insertRow();
        for (const key in data[0]) {
            const headerCell = document.createElement("th");
            headerCell.textContent = key.charAt(0).toUpperCase() + key.slice(1);
            headerRow.appendChild(headerCell);
        }

        data.forEach(item => {
            const row = table.insertRow();
            for (const key in item) {
                const cell = row.insertCell();
                const formattedValue = key === "marketValue" || key === "ceilingPrice" || key === "invest" || key === "grahamIndex"
                    ? formatCurrency(Number(item[key]))
                    : item[key];
                cell.textContent = formattedValue;
            }
        });

        appDiv.appendChild(table);

        const backButton = document.createElement("button");
        backButton.textContent = "Back";
        backButton.addEventListener("click", function () {
            window.location.reload();
        });
        appDiv.appendChild(backButton);
    }
});
