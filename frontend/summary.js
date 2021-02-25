function SummaryInit(shortenedURLs) {
    this.body = JSON.stringify(shortenedURLs);
    this.headers = {
        "Content-Type": "application/json"
    };
    this.method = "POST";
}

async function summarize(shortenedURLs) {
    let returnThis;
    let promise = swaggerClient.then(
        client => client.apis.api.shortenedSummary({shortenedURLs: shortenedURLs}),
        reason => console.error('failed to load the spec: ' + reason)
    )
        .then(
            shortenedSummaryResult => returnThis = JSON.parse(shortenedSummaryResult.data),
            reason => console.error('failed on api call: ' + reason)
        );
    await promise;
    return returnThis;
}