async function loadPrefix() {
    swaggerClient
        .then(
            client => client.apis.api.shortenedPrefix(),
            reason => console.error('failed to load the spec: ' + reason)
        )
        .then(
            shortenedPrefixResult => document.getElementById("httpPrefix").textContent = JSON.parse(shortenedPrefixResult.data),
            reason => console.error('failed on api call: ' + reason)
        )
}