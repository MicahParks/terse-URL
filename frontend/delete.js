async function deleteShortened(shortenedURLs) {
    let resultPromise;
    let promise = swaggerClient
        .then(
            client => client.apis.api.shortenedDelete({shortenedURLs: shortenedURLs}),
            reason => console.error('failed to load the spec: ' + reason)
        )
        .then(
            shortenedDeleteResult => resultPromise = {},
            reason => console.error('failed on api call: ' + reason)
        );
    await promise;
    return resultPromise;
}

async function deleteRow(shortenedURLs) {
    deleteShortened(shortenedURLs).then(function () {
        buildTable();
    });
    $('#deleteModal').modal('hide');
    $('#deleteCheckedModal').modal('hide');
}
