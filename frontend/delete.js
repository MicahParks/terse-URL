async function deleteShortened(shortenedURLs) {
    swaggerClient.then(
        client => client.apis.api.shortenedDelete({shortenedURLs: shortenedURLs}),
        reason => console.error('failed to load the spec: ' + reason)
    )
        .then(
            reason => console.error('failed on api call: ' + reason)
        );
}
