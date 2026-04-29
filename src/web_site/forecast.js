const forecastData = [
    {
        period: "Июнь 2026",
        name: "Расходомер РМ-08",
        type: "Расходомер",
        reason: "Требует замены по сроку службы",
        plan: "До 02.06.2026",
        planDate: "2026-06-02",
        quantity: "2 шт."
    },
    {
        period: "Август 2026",
        name: "Тахометр ТХ-11",
        type: "Тахометр",
        reason: "Критическое выбытие оборудования",
        plan: "До 12.08.2026",
        planDate: "2026-08-12",
        quantity: "1 шт."
    },
    {
        period: "Сентябрь 2026",
        name: "Манометр МВП-63",
        type: "Манометр",
        reason: "Плановая замена в компрессорном контуре",
        plan: "До 19.09.2026",
        planDate: "2026-09-19",
        quantity: "3 шт."
    },
    {
        period: "Ноябрь 2026",
        name: "Газоанализатор ГА-300",
        type: "Газоанализатор",
        reason: "Обновление парка по прогнозу риска",
        plan: "До 01.11.2026",
        planDate: "2026-11-01",
        quantity: "2 шт."
    },
    {
        period: "Декабрь 2026",
        name: "Манометр МП-100",
        type: "Манометр",
        reason: "Приближается к снятию",
        plan: "До 15.12.2026",
        planDate: "2026-12-15",
        quantity: "4 шт."
    },
    {
        period: "Январь 2027",
        name: "Весы платформенные ВП-500",
        type: "Весы",
        reason: "Резерв под замену складского узла",
        plan: "До 30.01.2027",
        planDate: "2027-01-30",
        quantity: "2 шт."
    }
];

const forecastState = {
    sourceData: forecastData,
    filteredData: [...forecastData],
    currentPage: 1,
    rowsPerPage: 10,
    startDate: "",
    endDate: ""
};

const forecastElements = {
    startDate: document.getElementById("forecast-start-date"),
    endDate: document.getElementById("forecast-end-date"),
    tableBody: document.getElementById("forecast-table-body"),
    summary: document.getElementById("forecast-summary"),
    pageInfo: document.getElementById("forecast-page-info"),
    prevPage: document.getElementById("forecast-prev-page"),
    nextPage: document.getElementById("forecast-next-page"),
    paginationNumbers: document.getElementById("forecast-pagination-numbers")
};

function applyForecastFilters() {
    forecastState.filteredData = forecastState.sourceData.filter((item) => {
        const matchesStartDate = !forecastState.startDate || item.planDate >= forecastState.startDate;
        const matchesEndDate = !forecastState.endDate || item.planDate <= forecastState.endDate;

        return matchesStartDate && matchesEndDate;
    });

    const totalPages = Math.max(1, Math.ceil(forecastState.filteredData.length / forecastState.rowsPerPage));
    if (forecastState.currentPage > totalPages) {
        forecastState.currentPage = totalPages;
    }

    forecastElements.startDate.max = forecastState.endDate || "";
    forecastElements.endDate.min = forecastState.startDate || "";
}

function renderForecastPagination(totalPages) {
    forecastElements.paginationNumbers.innerHTML = "";

    for (let page = 1; page <= totalPages; page += 1) {
        const button = document.createElement("button");
        button.type = "button";
        button.className = `pagination__number${page === forecastState.currentPage ? " is-active" : ""}`;
        button.textContent = String(page);
        button.addEventListener("click", () => {
            forecastState.currentPage = page;
            renderForecastTable();
        });
        forecastElements.paginationNumbers.appendChild(button);
    }
}

function renderForecastTable() {
    applyForecastFilters();

    const startIndex = (forecastState.currentPage - 1) * forecastState.rowsPerPage;
    const pageItems = forecastState.filteredData.slice(startIndex, startIndex + forecastState.rowsPerPage);
    const totalPages = Math.max(1, Math.ceil(forecastState.filteredData.length / forecastState.rowsPerPage));

    forecastElements.tableBody.innerHTML = pageItems.length
        ? pageItems.map((item) => `
            <tr>
                <td>${item.period}</td>
                <td><strong>${item.name}</strong></td>
                <td>${item.type}</td>
                <td>${item.reason}</td>
                <td>${item.plan}</td>
                <td>${item.quantity}</td>
            </tr>
        `).join("")
        : `
            <tr>
                <td colspan="6">
                    <div class="empty-state">За выбранный период плановых закупок не найдено.</div>
                </td>
            </tr>
        `;

    forecastElements.summary.textContent = `Показано ${pageItems.length} из ${forecastState.filteredData.length} записей`;
    forecastElements.pageInfo.textContent = `Страница ${forecastState.currentPage} из ${totalPages}`;
    forecastElements.prevPage.disabled = forecastState.currentPage === 1;
    forecastElements.nextPage.disabled = forecastState.currentPage === totalPages;
    renderForecastPagination(totalPages);
}

function bindForecastEvents() {
    forecastElements.startDate.addEventListener("change", () => {
        forecastState.startDate = forecastElements.startDate.value;
        forecastState.currentPage = 1;
        renderForecastTable();
    });

    forecastElements.endDate.addEventListener("change", () => {
        forecastState.endDate = forecastElements.endDate.value;
        forecastState.currentPage = 1;
        renderForecastTable();
    });

    forecastElements.prevPage.addEventListener("click", () => {
        if (forecastState.currentPage > 1) {
            forecastState.currentPage -= 1;
            renderForecastTable();
        }
    });

    forecastElements.nextPage.addEventListener("click", () => {
        const totalPages = Math.max(1, Math.ceil(forecastState.sourceData.length / forecastState.rowsPerPage));
        if (forecastState.currentPage < totalPages) {
            forecastState.currentPage += 1;
            renderForecastTable();
        }
    });
}

document.addEventListener("DOMContentLoaded", () => {
    bindForecastEvents();
    renderForecastTable();
});