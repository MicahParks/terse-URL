const Operation = {
    INSERT: 0,
    UPDATE: 1,
    UPSERT: 2,
}

const writeInit = {
    headers: {
        "Content-Type": "application/json"
    },
    method: "POST"
}

async function write(operation, terse) {
    let init = writeInit;
    init.body = JSON.stringify(terse);
    return fetch(`/api/write/${operation}`, init)
}
