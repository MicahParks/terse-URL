function SummaryInit(summary) {
    this.body = JSON.stringify(summary);
    this.headers = {
        "Content-Type": "application/json"
    };
    this.method = "POST";
}

async function summaryAll(summary) {
    return fetch("/api/summary", new SummaryInit(summary))
        .then(function (response) {
            return response.json();
        })
}