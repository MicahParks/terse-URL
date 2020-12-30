function DeleteInit() {
    this.body = JSON.stringify({
        terse: true, // TODO Configurable?
        visits: true
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

async function deleteOne(shortened) {
    return fetch(`/api/delete/${shortened}`, new DeleteInit())
        .then(function (response) {
            return response.json();
        });
}
