function toggleExpand(noteName) {
    var content = document.querySelector("div#expand-" + noteName);
    content.classList.toggle("expanded");
}