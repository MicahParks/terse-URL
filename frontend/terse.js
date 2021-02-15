async function getTerse(shortenedURL) {
    return fetch("/api/terse/" + shortenedURL)
        .then(function (terse) {
            return terse.json();
        });
}
