async function getMeta(originalURL) {
    let resultPromise;
    let promise = swaggerClient
        .then(
            client => client.apis.api.frontendMeta({originalURL: JSON.stringify(originalURL)}),
            reason => console.error('failed to load the spec: ' + reason)
        )
        .then(
            frontendMetaResult => resultPromise = JSON.parse(frontendMetaResult.data),
            reason => console.error('failed on api call: ' + reason)
        );
    await promise;
    return resultPromise;
}
