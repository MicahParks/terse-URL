const deleteInit = {
    body: JSON.stringify({
        terse: true, // TODO Configurable?
        visits: true
    }),
    headers: {
        "Content-Type": "application/json"
    },
    method: "DELETE"
}

async function deleteAll() {
    return fetch("/api/delete", deleteInit)
        .then(function (response) {
            return response.json();
        })
}

async function deleteOne(shortened) {
    return fetch(`/api/delete/${shortened}`, deleteInit)
        .then(function (response) {
            return response.json();
        })
}
