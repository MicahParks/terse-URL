async function exportShortened(shortenedURLs) {
    let resultPromise;
    let promise = swaggerClient
        .then(
            client => client.apis.api.export({shortenedURLs: shortenedURLs}),
            reason => console.error('failed to load the spec: ' + reason)
        )
        .then(
            exportResult => resultPromise = JSON.parse(exportResult.data),
            reason => console.error('failed on api call: ' + reason)
        );
    await promise;
    return resultPromise;
}
