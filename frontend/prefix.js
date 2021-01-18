async function getPrefix() {
    let response = await fetch(`/api/prefix`);
    response.json().then(prefix => {
        document.getElementById("httpPrefix").textContent = prefix;
    });
}
