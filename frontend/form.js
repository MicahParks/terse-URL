async function submitForm(e) {

    e.preventDefault();

    const submitButton = document.getElementById("submitButton");
    submitButton.disabled = true;
    setTimeout(() => submitButton.disabled = false, 1000); // TODO

    let terse = new Terse();
    terse.originalURL = document.getElementById("originalURL").value;
    terse.shortenedURL = document.getElementById("shortenedURL").value;
    terse.javascriptTracking = document.getElementById("jsTracking").checked;

    let operation = document.getElementById("writeOperation").value;

    terse.redirectType = $("input[name=redirectType]:checked", "#redirectType").val();

    if (terse.redirectType === "meta" || terse.redirectType === "js") {

        let htmlTitle = $("#htmlTitle").val();

        let og = makeMetaMap("#ogMeta :input");
        let twitter = makeMetaMap("#twitterMeta :input");

        terse.mediaPreview = new MediaPreview(og, htmlTitle, twitter);
    }

    await write(operation, terse);
    $("#formModal").modal("hide");
    buildTable(); // TODO Maybe don't do this every time?
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

function populateForm(terse) {

    document.getElementById("originalURL").value = terse.originalURL;
    document.getElementById("shortenedURL").value = terse.shortenedURL;
    terse.javascriptTracking = document.getElementById("jsTracking").checked = terse.javascriptTracking;

    document.getElementById("writeOperation").selectedIndex = 2;

    $("#" + terse.redirectType + "Redirect").prop("checked", true);

    hideShowHTMLForm();

    if (terse.redirectType === "meta" || terse.redirectType === "js") {

        $('#advanced').collapse('show'); // TODO Make sure this doesn't close the advanced collapse in certain conditions.

        if (terse.mediaPreview !== undefined) {
            $("#htmlTitle").value = terse.mediaPreview.title;
            clearPreview();
            populateMetaMap(terse.mediaPreview.og, terse.mediaPreview.twitter);
        }
    }
}

function populateMetaMap(og, twitter) {

    for (let key in og) {
        let clone = ogInput.cloneNode(true);
        clone.id = ogInput.id + ogInputCounter;
        ogInputCounter++;

        clone.childNodes[1].value = key;
        clone.childNodes[3].value = og[key];

        document.getElementById("ogMeta").appendChild(clone);
    }

    for (let key in twitter) {
        let clone = twitterInput.cloneNode(true);
        clone.id = twitterInput.id + twitterInputCounter;
        twitterInputCounter++;

        clone.childNodes[1].value = key;
        clone.childNodes[3].value = twitter[key];

        document.getElementById("twitterMeta").appendChild(clone);
    }
}

function clearPreview() {
    removeAllChildNodes(document.getElementById("ogMeta"));
    removeAllChildNodes(document.getElementById("twitterMeta"));
}
