const siData = [
    {
        id: "SI-001",
        name: "Весы лабораторные ВЛ-210",
        type: "Весы",
        commissioningDate: "14.03.2018",
        mpi: "12 мес.",
        status: "В норме",
        remaining: "237 дней",
        description: "Лабораторные весы высокой точности для аналитических измерений."
    },
    {
        id: "SI-002",
        name: "Манометр МП-100",
        type: "Манометр",
        commissioningDate: "05.11.2016",
        mpi: "12 мес.",
        status: "Приближается к снятию",
        remaining: "Снятие 15.12.2026",
        description: "Манометр для контроля давления в трубопроводах и технологических линиях."
    },
    {
        id: "SI-003",
        name: "Датчик температуры ТСМ-920",
        type: "Датчик",
        commissioningDate: "22.07.2019",
        mpi: "6 мес.",
        status: "В норме",
        remaining: "412 дней",
        description: "Промышленный датчик температуры для производственных агрегатов."
    },
    {
        id: "SI-004",
        name: "Расходомер РМ-08",
        type: "Расходомер",
        commissioningDate: "11.02.2015",
        mpi: "12 мес.",
        status: "Требует замены",
        remaining: "Снятие 02.06.2026",
        description: "Средство измерения расхода жидкости на узле учета."
    },
    {
        id: "SI-005",
        name: "Весы платформенные ВП-500",
        type: "Весы",
        commissioningDate: "29.08.2017",
        mpi: "12 мес.",
        status: "Приближается к снятию",
        remaining: "Снятие 30.01.2027",
        description: "Платформенные весы для складских и приемо-сдаточных операций."
    },
    {
        id: "SI-006",
        name: "Термометр ТЛ-4",
        type: "Термометр",
        commissioningDate: "10.12.2020",
        mpi: "6 мес.",
        status: "В норме",
        remaining: "518 дней",
        description: "Технический термометр для контроля температуры в лаборатории."
    },
    {
        id: "SI-007",
        name: "Манометр МВП-63",
        type: "Манометр",
        commissioningDate: "18.04.2014",
        mpi: "12 мес.",
        status: "Требует замены",
        remaining: "Снятие 19.09.2026",
        description: "Вакуумметрический манометр для компрессорного оборудования."
    },
    {
        id: "SI-008",
        name: "Датчик давления ДД-51",
        type: "Датчик",
        commissioningDate: "02.05.2021",
        mpi: "12 мес.",
        status: "В норме",
        remaining: "604 дня",
        description: "Датчик давления для мониторинга технологических емкостей."
    },
    {
        id: "SI-009",
        name: "Газоанализатор ГА-300",
        type: "Газоанализатор",
        commissioningDate: "16.01.2018",
        mpi: "12 мес.",
        status: "Приближается к снятию",
        remaining: "Снятие 01.11.2026",
        description: "Анализатор состава газовой среды для промышленных площадок."
    },
    {
        id: "SI-010",
        name: "Осциллограф ОС-220",
        type: "Осциллограф",
        commissioningDate: "25.09.2019",
        mpi: "24 мес.",
        status: "В норме",
        remaining: "690 дней",
        description: "Измерительный прибор для анализа электрических сигналов."
    },
    {
        id: "SI-011",
        name: "Тахометр ТХ-11",
        type: "Тахометр",
        commissioningDate: "07.06.2013",
        mpi: "12 мес.",
        status: "Требует замены",
        remaining: "Снятие 12.08.2026",
        description: "Тахометр для измерения скорости вращения оборудования."
    },
    {
        id: "SI-012",
        name: "Влагомер ВМ-75",
        type: "Влагомер",
        commissioningDate: "30.10.2020",
        mpi: "12 мес.",
        status: "В норме",
        remaining: "540 дней",
        description: "Прибор для контроля влажности сырья и производственных материалов."
    }
];

const state = {
    sourceData: siData,
    filteredData: [...siData],
    currentPage: 1,
    rowsPerPage: 10,
    expiringOnly: false,
    currentQuery: "",
    expiringRange: "all",
    searchDebounceId: null,
    searchController: null
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

function getStatusClass(status) {
    if (status === "В норме") {
        return "status-badge status-badge--normal";
    }
    if (status === "Приближается к снятию") {
        return "status-badge status-badge--warning";
    }
    return "status-badge status-badge--critical";
}

function parseRemovalDate(value) {
    const match = value.match(/(\d{2})\.(\d{2})\.(\d{4})/);

    if (!match) {
        return null;
    }

    const [, day, month, year] = match;
    return new Date(Number(year), Number(month) - 1, Number(day));
}

function getExpiringRangeLimit(range) {
    const today = new Date();
    today.setHours(0, 0, 0, 0);

    const limit = new Date(today);

    if (range === "week") {
        limit.setDate(limit.getDate() + 7);
        return { start: today, end: limit };
    }

    if (range === "month") {
        limit.setMonth(limit.getMonth() + 1);
        return { start: today, end: limit };
    }

    if (range === "year") {
        limit.setFullYear(limit.getFullYear() + 1);
        return { start: today, end: limit };
    }

    return null;
}

function matchesExpiringRange(item) {
    if (!state.expiringOnly || state.expiringRange === "all") {
        return true;
    }

    const removalDate = parseRemovalDate(item.remaining);
    const range = getExpiringRangeLimit(state.expiringRange);

    if (!removalDate || !range) {
        return false;
    }

    return removalDate >= range.start && removalDate <= range.end;
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

function applyFilters() {
    const query = state.currentQuery.trim().toLowerCase();

    state.filteredData = state.sourceData.filter((item) => {
        const matchesSearch = !query
            || item.name.toLowerCase().includes(query)
            || item.type.toLowerCase().includes(query)
            || item.description.toLowerCase().includes(query);

        const matchesExpiring = !state.expiringOnly
            || item.status === "Приближается к снятию"
            || item.status === "Требует замены";

        return matchesSearch && matchesExpiring && matchesExpiringRange(item);
    });

    const totalPages = Math.max(1, Math.ceil(state.filteredData.length / state.rowsPerPage));
    if (state.currentPage > totalPages) {
        state.currentPage = totalPages;
    }
}

function renderTable() {
    applyFilters();
    updateRangeFilters();
    updateExpiringLayout();
    updateModeButtons();

    const startIndex = (state.currentPage - 1) * state.rowsPerPage;
    const pageItems = state.filteredData.slice(startIndex, startIndex + state.rowsPerPage);
    const totalPages = Math.max(1, Math.ceil(state.filteredData.length / state.rowsPerPage));

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
                    <div class="empty-state">Ничего не найдено. Измените поисковый запрос или отключите фильтр.</div>
                </td>
            </tr>
        `;

    elements.tableSummary.textContent = `Показано ${pageItems.length} из ${state.filteredData.length} записей`;
    elements.pageInfo.textContent = `Страница ${state.currentPage} из ${totalPages}`;
    elements.prevPage.disabled = state.currentPage === 1;
    elements.nextPage.disabled = state.currentPage === totalPages;

    renderPagination(totalPages);
}

function renderPagination(totalPages) {
    elements.paginationNumbers.innerHTML = "";

    for (let page = 1; page <= totalPages; page += 1) {
        const button = document.createElement("button");
        button.type = "button";
        button.className = `pagination__number${page === state.currentPage ? " is-active" : ""}`;
        button.textContent = String(page);
        button.addEventListener("click", () => {
            state.currentPage = page;
            renderTable();
        });
        elements.paginationNumbers.appendChild(button);
    }
}

function setLoading(isLoading) {
    elements.searchInputWrap.classList.toggle("is-loading", isLoading);
}

function abortPendingSearch() {
    if (state.searchDebounceId) {
        window.clearTimeout(state.searchDebounceId);
        state.searchDebounceId = null;
    }

    if (state.searchController) {
        state.searchController.abort();
        state.searchController = null;
    }
}

function performSearchRequest(query, signal) {
    return new Promise((resolve, reject) => {
        const timeoutId = window.setTimeout(() => {
            resolve(query);
        }, 300);

        signal.addEventListener("abort", () => {
            window.clearTimeout(timeoutId);
            reject(new DOMException("Search aborted", "AbortError"));
        }, { once: true });
    });
}

function scheduleSearch() {
    const nextQuery = elements.searchInput.value;

    abortPendingSearch();

    state.searchDebounceId = window.setTimeout(async () => {
        const controller = new AbortController();
        state.searchController = controller;
        state.currentPage = 1;
        setLoading(true);

        try {
            const resolvedQuery = await performSearchRequest(nextQuery, controller.signal);
            state.currentQuery = resolvedQuery;
            renderTable();
        } catch (error) {
            if (!(error instanceof DOMException) || error.name !== "AbortError") {
                throw error;
            }
        } finally {
            if (state.searchController === controller) {
                state.searchController = null;
                setLoading(false);
            }
        }
    }, 400);
}

function setFilterMode(expiringOnly) {
    state.expiringOnly = expiringOnly;
    state.currentPage = 1;
    if (!state.expiringOnly) {
        state.expiringRange = "all";
    }
    renderTable();
    setLoading(false);
}

function setExpiringRange(range) {
    state.expiringRange = range;
    state.currentPage = 1;
    renderTable();
    setLoading(false);
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
        renderTable();
    });

    elements.prevPage.addEventListener("click", () => {
        if (state.currentPage > 1) {
            state.currentPage -= 1;
            renderTable();
        }
    });

    elements.nextPage.addEventListener("click", () => {
        const totalPages = Math.max(1, Math.ceil(state.filteredData.length / state.rowsPerPage));
        if (state.currentPage < totalPages) {
            state.currentPage += 1;
            renderTable();
        }
    });
}

document.addEventListener("DOMContentLoaded", () => {
    state.expiringOnly = getInitialViewMode();
    bindEvents();
    renderTable();
    setLoading(false);
});