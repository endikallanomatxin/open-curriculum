function addClass(el, className) {
    if (el.classList)
        el.classList.add(className);
}

function removeClass(el, className) {
    if (el.classList)
        el.classList.remove(className);
}


function isElementOverViewport(el) {
    var rect = el.getBoundingClientRect();
    return rect.bottom < 0;
}

function isElementUnderViewport(el) {
    var rect = el.getBoundingClientRect();
    return rect.top > (window.innerHeight || document.documentElement.clientHeight);
}


function checkElements() {
    var elements = document.querySelectorAll('*');
    // Seguramente lo mejor sea usar una clase .animated
    // y que solo se aplique a los elementos que queramos
    // que tengan animación
    for (var i = 0; i < elements.length; i++) {

        if (isElementOverViewport(elements[i])) {
            addClass(elements[i], 'over-viewport');
        }
        else {
            removeClass(elements[i], 'over-viewport');
        }

        if (isElementUnderViewport(elements[i])) {
            addClass(elements[i], 'under-viewport');
        }
        else {
            removeClass(elements[i], 'under-viewport');
        }
    }
}

document.addEventListener('DOMContentLoaded', function () {
    window.addEventListener('scroll', checkElements);
});