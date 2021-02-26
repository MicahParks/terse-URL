async function getTerse(shortenedURL) {
    let resultPromise;
    let promise = swaggerClient
        .then(
            client => client.apis.api.terseRead({shortenedURLs: [shortenedURL]}),
            reason => console.error('failed to load the spec: ' + reason)
        )
        .then(
            terseReadResult => resultPromise = JSON.parse(terseReadResult.data),
            reason => console.error('failed on api call: ' + reason)
        );
    await promise;
    return resultPromise;
}
