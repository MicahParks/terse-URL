async function summarize(shortenedURLs) {
    let resultPromise;
    let promise = swaggerClient.then(
        client => client.apis.api.shortenedSummary({shortenedURLs: shortenedURLs}),
        reason => console.error('failed to load the spec: ' + reason)
    )
        .then(
            shortenedSummaryResult => resultPromise = JSON.parse(shortenedSummaryResult.data),
            reason => console.error('failed on api call: ' + reason)
        );
    await promise;
    return resultPromise;
}