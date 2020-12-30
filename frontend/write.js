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
    let init = new WriteInit();
    init.body = JSON.stringify(terse);
    return fetch(`/api/write/${operation}`, init);
}
