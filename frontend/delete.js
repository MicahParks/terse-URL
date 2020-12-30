async function deleteAll() {
    let init = {
        body: JSON.stringify({
            terse: true,
            visits: true
        }),
        headers: {
            "Content-Type": "application/json"
        },
        method: "DELETE"
    }
    return fetch("/api/delete", init)
        .then(function (response) {
            return response.json();
        })
}