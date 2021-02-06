async function submitForm(e) {

    e.preventDefault();

    const submitButton = document.getElementById("submitButton");
    submitButton.disabled = true;
    setTimeout(() => submitButton.disabled = false, 1000); // TODO

    let terse = new Terse();
    terse.originalURL = document.getElementById("originalURL").value;
    terse.shortenedURL = document.getElementById("shortenedURL").value;

    let operation = document.getElementById("writeOperation").value;

    terse.redirectType = $("input[name=redirectType]:checked", "#redirectType").val();

    if (terse.redirectType === "meta" || terse.redirectType === "js") {

        let htmlTitle = $("#htmlTitle").val();

        let og = makeMetaMap("#ogMeta :input");
        let twitter = makeMetaMap("#twitterMeta :input");

        let inherit = true; // TODO

        terse.mediaPreview = new MediaPreview(inherit, og, htmlTitle, twitter);
    }

    await write(operation, terse);
}

function makeMetaMap(query) {
    let metaMap = new Map();
    let index = 0;
    let key = "";
    for (let child of $(query)) {
        if (child.type === "text") {
            if (index % 2 === 0) {
                key = child.value;
            } else {
                metaMap[key] = child.value;
            }
            index++;
        }
    }
    return metaMap;
}
