const Controller = {
    search: (ev) => {
        ev.preventDefault();
        const form = document.getElementById("form");
        const data = Object.fromEntries(new FormData(form));
        const response = fetch(`/api/search?q=${data.query}`).then((response) => {
            response.json().then((results) => {
                Controller.updateTable(results);
            });
        });
    },

    updateTable: (results) => {
        const table = document.getElementById("table-body");
        const rows = [];
        rows.push(`<tr><td>Line</td><td>Appears on these line number(s)</td><tr/>`);
        for (let result of results) {
            let line = "";
            let bold = false;
            for (let i = 0; i < result.line.length; i++) {
                if (result.matchedIndexes.includes(i)) {
                    if (!bold) {
                        line += "<strong>";
                    }
                    line += result.line.charAt(i);
                    bold = true;
                } else {
                    if (bold) {
                        line += "</strong>";
                    }
                    line += result.line.charAt(i);
                    bold = false;
                }
            }
            if (bold) {
                line += "</strong>"
            }
            rows.push(`<tr><td>${line}</td><td>${result.lineNumbers}</td><tr/>`);
        }
        table.innerHTML = rows;
    },
};

const form = document.getElementById("form");
form.addEventListener("submit", Controller.search);
