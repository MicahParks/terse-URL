async function exportAll() {
    return fetch("/api/export")
        .then(function (response) {
            return response.json();
        })
}

async function exportOne(shortened) {
    return fetch(`/api/export/${shortened}`)
        .then(function (response) {
            return response.json();
        })
}
