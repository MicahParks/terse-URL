function MetaInit(originalURL) {
    this.body = JSON.stringify(originalURL);
    this.headers = {
        "Content-Type": "application/json"
    }
    this.method = "POST";
}

async function getMeta(originalURL) {
    return fetch("/api/frontend/meta", new MetaInit(originalURL))
        .then(function (response) {
            return response.json();
        })
}