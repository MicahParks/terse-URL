function SummaryInit(shortenedURLs) {
    this.body = JSON.stringify(shortenedURLs);
    this.headers = {
        "Content-Type": "application/json"
    };
    this.method = "POST";
}

async function summarize(shortenedURLs) {
    return fetch("/api/summary", new SummaryInit(shortenedURLs))
        .then(function (response) {
            return response.json();
        })
}