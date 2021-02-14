function DeleteInit(shortenedURLs) {
    this.body = JSON.stringify({
        "delete": {
            terse: true, // TODO Configurable?
            visits: true
        },
        "shortenedURLs": shortenedURLs
    });
    this.headers = {
        "Content-Type": "application/json"
    };
    this.method = "DELETE";
}

async function deleteAll() {
    return fetch("/api/delete", new DeleteInit())
        .then(function (response) {
            return response.json();
        });
}

async function deleteSome(shortenedURLs) {
    return fetch(`/api/delete/some`, new DeleteInit(shortenedURLs));
}
