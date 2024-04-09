// This script is used to create a menu for a documet
// It has to be automatically generated from the headings of the document

// Create a variable that stores the skeleton of the menu

function initializeMenu() {
    var menu = document.createElement("div");
    menu.id = "menu";
    document.body.appendChild(menu);

    var headings = document.querySelectorAll("h1, h2, h3, h4, h5, h6");

    // Give each heading an id containing the heding level
    for (var i = 0; i < headings.length; i++) {
        var heading = headings[i];
        heading.id = heading.tagName + '-'+ i;
    }

    for (var i = 0; i < headings.length; i++) {
        var heading = headings[i];
        var menuItem = document.createElement("a");
        menuItem.href = "#" + heading.id;
        menuItem.textContent = heading.textContent;
        menuItem.classList.add('collapsed');
        menu.appendChild(menuItem);
    }
}

// There is a function in other file that assigns over-viewport and under-viewport classes to elements
function checkHeadings() {

    // Select the headings that have none of the classes
    var headingsOutOfViewPort = document.querySelectorAll('h1.over-viewport, h1.under-viewport, h2.over-viewport, h2.under-viewport, h3.over-viewport, h3.under-viewport, h4.over-viewport, h4.under-viewport, h5.over-viewport, h5.under-viewport, h6.over-viewport, h6.under-viewport');
    var headingsInViewPort = document.querySelectorAll('h1:not(.over-viewport):not(.under-viewport), h2:not(.over-viewport):not(.under-viewport), h3:not(.over-viewport):not(.under-viewport), h4:not(.over-viewport):not(.under-viewport), h5:not(.over-viewport):not(.under-viewport), h6:not(.over-viewport):not(.under-viewport)');

    // In case there are no headings in the viewport
    if (headingsInViewPort.length == 0) {
        // Search for the closest heading over the viewport
        var headingsOverViewport = document.querySelector('h1.over-viewport, h2.over-viewport, h3.over-viewport, h4.over-viewport, h5.over-viewport, h6.over-viewport');
        // Take the last
        var heading = headingsOverViewport[headingsOverViewport.length-1];
        // Append to the list
        headingsInViewPort.push(heading);
    }


    // Add the class collapsed to the menu items that correspond to
    // the headings that are out of the viewport
    for (var i = 0; i < headingsOutOfViewPort.length; i++) {
        var heading = headingsOutOfViewPort[i];
        var headingId = heading.id;
        var menuItem = document.querySelector("#menu a[href='#" + headingId + "']");
        menuItem.classList.add('collapsed');
    }

    // Remove the class collapsed from the headings that are in the viewport
    // All also the ones they belong to
    for (var i = 0; i < headingsInViewPort.length; i++) {
        var heading = headingsInViewPort[i];
        var headingId = heading.id;
        var menuItem = document.querySelector("#menu a[href='#" + headingId + "']");
        menuItem.classList.remove('collapsed');

        // Now search through the previous menu items and remove the class collapsed
        // that have lower level than highest previous one found

        var previousHighestLevel = heading.tagName[1];
        var previousMenuItem = menuItem.previousElementSibling;
        
        while (previousMenuItem) {
            if (previousMenuItem.tagName == "A") {
                var currentLevel = previousMenuItem.href.split('-')[0].slice(-1);
                if (currentLevel <= previousHighestLevel) {
                    previousMenuItem.classList.remove('collapsed');
                    previousHighestLevel = currentLevel;
                }
            }
            previousMenuItem = previousMenuItem.previousElementSibling;
        }
    }

    // Now the same forwards
    for (var i = 0; i < headingsInViewPort.length; i++) {
        var heading = headingsInViewPort[i];
        var headingId = heading.id;
        var menuItem = document.querySelector("#menu a[href='#" + headingId + "']");
        menuItem.classList.remove('collapsed');

        // Now search through the next menu items and remove the class collapsed
        // that same level than highest next one found

        var nextHighestLevel = heading.tagName[1];
        var nextMenuItem = menuItem.nextElementSibling;
        
        while (nextMenuItem) {
            if (nextMenuItem.tagName == "A") {
                var currentLevel = nextMenuItem.href.split('-')[0].slice(-1);
                if (currentLevel <= nextHighestLevel) {
                    nextMenuItem.classList.remove('collapsed');
                    nextHighestLevel = currentLevel;
                }
            }
            nextMenuItem = nextMenuItem.nextElementSibling;
        }
    }
}

initializeMenu();
// Call the function when the user scrolls but not too often
window.addEventListener('scroll', function() {
    clearTimeout(window.scrollTimer);
    window.scrollTimer = setTimeout(checkHeadings, 200);
});