async function getPrefix() {
    let thing = fetch(`/api/prefix`);
    return thing;
}
