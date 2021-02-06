function generateTableBody(table, tableData) {
    for (let entry of tableData) {
        let row = table.insertRow();
        for (let key of entry) {
            let cell = row.insertCell();
            let text = document.createTextNode(entry[key]);
            cell.appendChild(text);
        }
    }
}

function generateTableHead(table, firstRow) {
    let thead = table.createTHead();
    let row = thead.insertRow();
    for (let key of firstRow) {
        let th = document.createElement("th");
        let text = document.createTextNode(key);
        th.appendChild(text);
        row.appendChild(th);
    }
}

function createTable(selector, tableData) {
    let table = document.querySelector(selector);
    generateTableHead(table, tableData[0]); // TODO Confirm data is at least 1 long.
    generateTableBody(table, tableData);
}

function createTerseList() {
    let list = document.getElementById("terseList");
    // list.
}
