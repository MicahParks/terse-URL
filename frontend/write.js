const Operation = {
    INSERT: 0,
    UPDATE: 1,
    UPSERT: 2,
}

function WriteInit() {
    this.headers = {"Content-Type": "application/json"};
    this.method = "POST";
}

async function write(operation, terse) {
    swaggerClient.then(
        client => client.apis.api.terseWrite({terse: [terse], operation: operation}),
        reason => console.error('failed to load the spec: ' + reason)
    )
        .then(
            reason => console.error('failed on api call: ' + reason)
        );
}
