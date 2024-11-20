function collapseHeader() {
    var body = document.body;
    // give it the class header-is-collapsed
    body.classList.add('header-is-collapsed');
    
    setTimeout(redrawArrows, 1500);
}

function expandHeader() {
    var body = document.body;
    // remove the class header-is-collapsed
    body.classList.remove('header-is-collapsed');
    setTimeout(redrawArrows, 1500);
}
