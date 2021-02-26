async function write(operation, terse) {
    let resultPromise;
    let promise = swaggerClient
        .then(
            client => client.apis.api.terseWrite({terse: [terse], operation: operation}),
            reason => console.error('failed to load the spec: ' + reason)
        )
        .then(
            terseWriteResult => resultPromise = JSON.parse(terseWriteResult.data),
            reason => console.error('failed on api call: ' + reason)
        );
    await promise;
    return resultPromise;
}
