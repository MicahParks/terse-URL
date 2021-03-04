async function buildTable() {

    document.getElementById("checkAll").checked = false;
    removeAllChildNodes(table);

    let index = 0;
    return summarize(null).then(function (summaries) {
        for (const [_, summary] of Object.entries(summaries)) { // TODO Using _ as blank identifier...

            let row = template.cloneNode(true);
            row.id = summary.terse.shortenedURL;
            index++;

            let checkbox = document.createElement('input');
            checkbox.type = 'checkbox';
            checkbox.value = summary.terse.shortenedURL;
            checkbox.id = summary.terse.shortenedURL + 'Checkbox';
            let label = document.createElement('label');
            label.htmlFor = checkbox.id;

            row.cells[0].appendChild(checkbox);
            row.cells[0].appendChild(label);
            row.cells[1].innerHTML = summary.terse.shortenedURL;
            row.cells[2].innerHTML = summary.terse.originalURL;
            row.cells[3].innerHTML = summary.terse.redirectType;
            if (summary.visits === undefined || summary.visits.visitCount === undefined) {
                row.cells[4].innerHTML = "0";
            } else {
                row.cells[4].innerHTML = summary.visits.visitCount;
            }

            table.appendChild(row);

            document.getElementById(checkbox.id).onchange = rowChecked; // TODO Duplicate?

            for (let button of $('#' + row.id + ' :button')) {
                button.setAttribute('data-bs-shortened', summary.terse.shortenedURL);
            }
        }
    });
}