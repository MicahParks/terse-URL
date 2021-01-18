async function getPrefix() { // TODO Use local storage for this?
    let response = await fetch(`/api/prefix`);
    response.json().then(prefix => {
        document.getElementById("httpPrefix").textContent = prefix;
    });
}
