function ExportInit(shortenedURLs) {
    this.body = JSON.stringify(shortenedURLs);
    this.headers = {
        "Content-Type": "application/json"
    };
    this.method = "POST";
}

async function exportAll() {
    return fetch("/api/export")
        .then(function (response) {
            return response.json();
        })
}

async function exportSome(shortenedURLs) {
    return fetch(`/api/export/some`, new ExportInit(shortenedURLs))
        .then(function (response) {
            return response.json();
        });
}
