const DEFAULT_API_BASE_URL = "http://localhost:8080";
const SEARCH_DEBOUNCE_MS = 400;
const MILLISECONDS_IN_DAY = 24 * 60 * 60 * 1000;

const state = {
    pageData: [],
    totalItems: 0,
    currentPage: 1,
    rowsPerPage: 10,
    expiringOnly: false,
    currentQuery: "",
    expiringRange: "all",
    searchDebounceId: null,
    requestController: null,
    errorMessage: ""
};

const elements = {
    searchInput: document.getElementById("search-input"),
    searchInputWrap: document.getElementById("search-input-wrap"),
    loadingIndicator: document.getElementById("loading-indicator"),
    allButton: document.getElementById("all-button"),
    expiringButton: document.getElementById("expiring-button"),
    forecastButton: document.getElementById("forecast-button"),
    expiringRangePanel: document.getElementById("expiring-range-panel"),
    expiringMetaHost: document.getElementById("expiring-meta-host"),
    rangeFilters: Array.from(document.querySelectorAll(".range-filter")),
    rowsPerPage: document.getElementById("rows-per-page"),
    tableBody: document.getElementById("table-body"),
    tableSummary: document.getElementById("table-summary"),
    tableHeader: document.getElementById("table-header"),
    tableHeaderInfo: document.getElementById("table-header-info"),
    tableHeaderMeta: document.getElementById("table-header-meta"),
    workspaceCard: document.querySelector(".workspace-card"),
    pageInfo: document.getElementById("page-info"),
    prevPage: document.getElementById("prev-page"),
    nextPage: document.getElementById("next-page"),
    paginationNumbers: document.getElementById("pagination-numbers")
};

function getInitialViewMode() {
    const searchParams = new URLSearchParams(window.location.search);
    return searchParams.get("view") === "expiring";
}

function getApiBaseUrl() {
    if (typeof window.METROLOG_API_BASE_URL === "string" && window.METROLOG_API_BASE_URL.trim()) {
        return window.METROLOG_API_BASE_URL.trim().replace(/\/$/, "");
    }

    if (window.location.origin === DEFAULT_API_BASE_URL) {
        return "";
    }

    return DEFAULT_API_BASE_URL;
}

function getStatusClass(status) {
    if (status === "В норме") {
        return "status-badge status-badge--normal";
    }
    if (status === "Приближается к снятию") {
        return "status-badge status-badge--warning";
    }
    if (status === "Требует замены") {
        return "status-badge status-badge--critical";
    }
    return "status-badge";
}

function parseApiDate(value) {
    if (!value || typeof value !== "string") {
        return null;
    }

    const [year, month, day] = value.slice(0, 10).split("-").map(Number);

    if (!year || !month || !day) {
        return null;
    }

    return new Date(year, month - 1, day);
}

function formatDate(value) {
    if (!(value instanceof Date) || Number.isNaN(value.getTime())) {
        return "Нет данных";
    }

    return value.toLocaleDateString("ru-RU");
}

function addMonths(date, months) {
    const result = new Date(date);
    result.setMonth(result.getMonth() + months);
    return result;
}

function calculateDayDiff(targetDate) {
    const today = new Date();
    today.setHours(0, 0, 0, 0);

    const dueDate = new Date(targetDate);
    dueDate.setHours(0, 0, 0, 0);

    return Math.ceil((dueDate.getTime() - today.getTime()) / MILLISECONDS_IN_DAY);
}

function mapApiItem(item) {
    const commissioningDate = parseApiDate(item.commissioning_date);
    const mpiMonths = typeof item.mpi === "number" ? item.mpi : null;
    const nextInspectionDate = commissioningDate && Number.isInteger(mpiMonths)
        ? addMonths(commissioningDate, mpiMonths)
        : null;

    let status = "Нет данных";
    let remaining = "Нет данных";

    if (nextInspectionDate) {
        const daysRemaining = calculateDayDiff(nextInspectionDate);

        if (daysRemaining < 0) {
            status = "Требует замены";
            remaining = `Снятие ${formatDate(nextInspectionDate)}`;
        } else if (daysRemaining <= 365) {
            status = "Приближается к снятию";
            remaining = `Снятие ${formatDate(nextInspectionDate)}`;
        } else {
            status = "В норме";
            remaining = `${daysRemaining} дней`;
        }
    }

    const passport = item.passport || "Нет данных";

    return {
        id: passport,
        name: item.name || "Без названия",
        type: item.type || "Нет данных",
        commissioningDate: formatDate(commissioningDate),
        mpi: Number.isInteger(mpiMonths) ? `${mpiMonths} мес.` : "Нет данных",
        status,
        remaining,
        description: `Паспорт: ${passport}`
    };
}

function getPageCount() {
    return Math.max(1, Math.ceil(state.totalItems / state.rowsPerPage));
}

function updateLoadingIndicator(message = "") {
    elements.loadingIndicator.textContent = message;
}

function buildApiUrl() {
    const params = new URLSearchParams({
        limit: String(state.rowsPerPage),
        offset: String((state.currentPage - 1) * state.rowsPerPage)
    });

    if (state.currentQuery.trim()) {
        params.set("query", state.currentQuery.trim());
    }

    if (state.expiringOnly) {
        params.set("expiring_range", state.expiringRange);
    }

    return `${getApiBaseUrl()}/api/miinstance?${params.toString()}`;
}

async function loadTableData() {
    if (state.requestController) {
        state.requestController.abort();
    }

    const controller = new AbortController();
    state.requestController = controller;
    state.errorMessage = "";
    setLoading(true);
    updateLoadingIndicator("Загрузка данных из Service Park...");

    try {
        const response = await fetch(buildApiUrl(), {
            method: "GET",
            signal: controller.signal,
            headers: {
                Accept: "application/json"
            }
        });

        const payload = await response.json().catch(() => null);

        if (!response.ok || !payload || !payload.success) {
            throw new Error((payload && payload.error) || `Request failed with status ${response.status}`);
        }

        state.pageData = Array.isArray(payload.data) ? payload.data.map(mapApiItem) : [];
        state.totalItems = typeof payload.total === "number" ? payload.total : 0;

        const totalPages = getPageCount();
        if (state.currentPage > totalPages) {
            state.currentPage = totalPages;
            await loadTableData();
            return;
        }

        renderTable();
        updateLoadingIndicator(state.totalItems === 0 ? "Совпадений не найдено." : "");
    } catch (error) {
        if (error instanceof DOMException && error.name === "AbortError") {
            return;
        }

        state.pageData = [];
        state.totalItems = 0;
        state.errorMessage = "Не удалось загрузить данные из Service Park. Проверь, что сервер запущен на http://localhost:8080.";
        renderTable();
        updateLoadingIndicator(state.errorMessage);
    } finally {
        if (state.requestController === controller) {
            state.requestController = null;
            setLoading(false);
        }
    }
}

function updateRangeFilters() {
    elements.expiringRangePanel.classList.toggle("hidden", !state.expiringOnly);

    elements.rangeFilters.forEach((button) => {
        button.classList.toggle("is-active", button.dataset.range === state.expiringRange);
    });
}

function updateExpiringLayout() {
    elements.workspaceCard.classList.toggle("is-expiring", state.expiringOnly);

    if (state.expiringOnly) {
        elements.expiringMetaHost.classList.remove("hidden");
        elements.expiringMetaHost.appendChild(elements.tableHeaderMeta);
        elements.tableHeader.classList.add("hidden");
        return;
    }

    elements.tableHeader.classList.remove("hidden");
    elements.tableHeader.appendChild(elements.tableHeaderMeta);
    elements.expiringMetaHost.classList.add("hidden");
}

function updateModeButtons() {
    elements.allButton.classList.toggle("button--primary", !state.expiringOnly);
    elements.allButton.classList.toggle("button--secondary", state.expiringOnly);
    elements.expiringButton.classList.toggle("button--primary", state.expiringOnly);
    elements.expiringButton.classList.toggle("button--secondary", !state.expiringOnly);
}

function renderTable() {
    updateRangeFilters();
    updateExpiringLayout();
    updateModeButtons();

    const pageItems = state.pageData;
    const totalPages = getPageCount();

    elements.tableBody.innerHTML = pageItems.length
        ? pageItems.map((item) => `
            <tr>
                <td>${item.id}</td>
                <td>
                    <div class="cell-title">
                        <strong>${item.name}</strong>
                        <span>${item.description}</span>
                    </div>
                </td>
                <td>${item.type}</td>
                <td>${item.commissioningDate}</td>
                <td>${item.mpi}</td>
                <td class="status-cell"><span class="${getStatusClass(item.status)}">${item.status}</span></td>
                <td>${item.remaining}</td>
            </tr>
        `).join("")
        : `
            <tr>
                <td colspan="7">
                    <div class="empty-state">${state.errorMessage || "Ничего не найдено. Измените поисковый запрос или отключите фильтр."}</div>
                </td>
            </tr>
        `;

    if (state.totalItems === 0) {
        elements.tableSummary.textContent = "Показано 0 записей";
    } else {
        const startIndex = (state.currentPage - 1) * state.rowsPerPage + 1;
        const endIndex = startIndex + pageItems.length - 1;
        elements.tableSummary.textContent = `Показано ${startIndex}-${endIndex} из ${state.totalItems} записей`;
    }

    elements.pageInfo.textContent = `Страница ${state.currentPage} из ${totalPages}`;
    elements.prevPage.disabled = state.currentPage === 1;
    elements.nextPage.disabled = state.currentPage === totalPages;

    renderPagination(totalPages);
}

function renderPagination(totalPages) {
    elements.paginationNumbers.innerHTML = "";

    if (totalPages <= 1) {
        return;
    }

    const visibleWindow = 5;
    let startPage = Math.max(1, state.currentPage - 2);
    let endPage = Math.min(totalPages, startPage + visibleWindow - 1);

    if (endPage - startPage + 1 < visibleWindow) {
        startPage = Math.max(1, endPage - visibleWindow + 1);
    }

    const pages = [];

    if (startPage > 1) {
        pages.push(1);
        if (startPage > 2) {
            pages.push("...");
        }
    }

    for (let page = startPage; page <= endPage; page += 1) {
        pages.push(page);
    }

    if (endPage < totalPages) {
        pages.push("...");
    }

    pages.forEach((page) => {
        if (page === "...") {
            const ellipsis = document.createElement("span");
            ellipsis.className = "pagination__ellipsis";
            ellipsis.textContent = "...";
            elements.paginationNumbers.appendChild(ellipsis);
            return;
        }

        const button = document.createElement("button");
        button.type = "button";
        button.className = `pagination__number${page === state.currentPage ? " is-active" : ""}`;
        button.textContent = String(page);
        button.addEventListener("click", () => {
            state.currentPage = page;
            loadTableData();
        });
        elements.paginationNumbers.appendChild(button);
    });
}

function setLoading(isLoading) {
    elements.searchInputWrap.classList.toggle("is-loading", isLoading);
}

function abortPendingSearch() {
    if (state.searchDebounceId) {
        window.clearTimeout(state.searchDebounceId);
        state.searchDebounceId = null;
    }
}

function scheduleSearch() {
    const nextQuery = elements.searchInput.value;

    abortPendingSearch();

    state.searchDebounceId = window.setTimeout(() => {
        state.currentPage = 1;
        state.currentQuery = nextQuery;
        loadTableData();
    }, SEARCH_DEBOUNCE_MS);
}

function setFilterMode(expiringOnly) {
    state.expiringOnly = expiringOnly;
    state.currentPage = 1;

    if (!state.expiringOnly) {
        state.expiringRange = "all";
    }

    loadTableData();
}

function setExpiringRange(range) {
    state.expiringRange = range;
    state.currentPage = 1;
    loadTableData();
}

function bindEvents() {
    elements.searchInput.addEventListener("input", () => {
        scheduleSearch();
    });

    elements.allButton.addEventListener("click", () => {
        setFilterMode(false);
    });

    elements.expiringButton.addEventListener("click", () => {
        setFilterMode(true);
    });

    elements.rangeFilters.forEach((button) => {
        button.addEventListener("click", () => {
            if (state.expiringOnly) {
                setExpiringRange(button.dataset.range || "all");
            }
        });
    });

    elements.rowsPerPage.addEventListener("change", () => {
        state.rowsPerPage = Number(elements.rowsPerPage.value);
        state.currentPage = 1;
        loadTableData();
    });

    elements.prevPage.addEventListener("click", () => {
        if (state.currentPage > 1) {
            state.currentPage -= 1;
            loadTableData();
        }
    });

    elements.nextPage.addEventListener("click", () => {
        const totalPages = getPageCount();
        if (state.currentPage < totalPages) {
            state.currentPage += 1;
            loadTableData();
        }
    });
}

document.addEventListener("DOMContentLoaded", () => {
    state.expiringOnly = getInitialViewMode();
    bindEvents();
    loadTableData();
});