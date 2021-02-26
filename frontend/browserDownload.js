function download(filename, text) {
    let element = document.createElement('a');
    element.setAttribute('href', 'data:text/plain;charset=utf-8,' + encodeURIComponent(text));
    element.setAttribute('download', filename);

    element.style.display = 'none';
    document.body.appendChild(element);

    element.click();

    document.body.removeChild(element);
}

function downloadExport(shortenedURLs) {
    exportShortened(shortenedURLs).then(function (exportData) {
        let filename = 'export.json';
        if (shortenedURLs.length === 1) {
            filename = shortenedURLs + '.json';
        }
        download(filename, JSON.stringify(exportData));
    });
}
