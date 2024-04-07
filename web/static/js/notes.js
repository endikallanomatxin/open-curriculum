function toggleNote(noteName) {
    var content = document.querySelector("div#note-" + noteName);
    content.classList.toggle("expanded");
}