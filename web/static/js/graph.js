// Now create a function that creates an svg arrow from one units node to another
function createArrow(from, to) {
    const nodeFrom = document.getElementById(`node-${from}`);
    const nodeTo = document.getElementById(`node-${to}`);

    const rectFrom = nodeFrom.getBoundingClientRect();
    const rectTo = nodeTo.getBoundingClientRect();

    const pointFrom = {
        x: rectFrom.left + window.scrollX + rectFrom.width / 2,
        y: rectFrom.top + window.scrollY + rectFrom.height / 2
    };

    const pointTo = {
        x: rectTo.left + window.scrollX + rectTo.width / 2,
        y: rectTo.top + window.scrollY + rectTo.height / 2
    };

    const dx = pointTo.x - pointFrom.x;
    const dy = pointTo.y - pointFrom.y;
    const distance = Math.sqrt(dx * dx + dy * dy);
    const unitX = dx / distance;
    const unitY = dy / distance;

    const fromOffset = 10;  // Starting point offset
    const toOffset = 16;    // Ending point offset to avoid overlap with arrowhead
    const adjustedFrom = {
        x: pointFrom.x + unitX * fromOffset,
        y: pointFrom.y + unitY * fromOffset
    };
    const adjustedTo = {
        x: pointTo.x - unitX * toOffset,
        y: pointTo.y - unitY * toOffset
    };

    // Calculate SVG position and size, adding a padding for the marker
    const padding = 10;  // Padding to ensure markers are within the SVG
    const minX = Math.min(adjustedFrom.x, adjustedTo.x) - padding;
    const maxX = Math.max(adjustedFrom.x, adjustedTo.x) + padding;
    const minY = Math.min(adjustedFrom.y, adjustedTo.y) - padding;
    const maxY = Math.max(adjustedFrom.y, adjustedTo.y) + padding;

    const svgWidth = maxX - minX;
    const svgHeight = maxY - minY;

    // Create the svg element with adjusted size and position
    const svg = document.createElementNS("http://www.w3.org/2000/svg", "svg");
    svg.classList.add("dependency-arrow");

    svg.setAttribute("width", svgWidth + "px");
    svg.setAttribute("height", svgHeight + "px");
    svg.style.position = "absolute";
    svg.style.top = minY + "px";
    svg.style.left = minX + "px";
    svg.style.zIndex = 2;

    const defs = document.createElementNS("http://www.w3.org/2000/svg", "defs");
    const marker = document.createElementNS("http://www.w3.org/2000/svg", "marker");
    marker.setAttribute("id", "arrowhead");
    marker.setAttribute("markerWidth", "6");
    marker.setAttribute("markerHeight", "4");
    marker.setAttribute("refX", "3");  // Small offset inside the marker for proper alignment
    marker.setAttribute("refY", "2");
    marker.setAttribute("orient", "auto");
    const polygon = document.createElementNS("http://www.w3.org/2000/svg", "polygon");
    polygon.setAttribute("points", "0 0, 6 2, 0 4");
    polygon.setAttribute("fill", "black");
    marker.appendChild(polygon);
    defs.appendChild(marker);
    svg.appendChild(defs);

    // Adjust path coordinates relative to the new SVG position
    const path = document.createElementNS("http://www.w3.org/2000/svg", "path");
    path.setAttribute("d", `M ${adjustedFrom.x - minX} ${adjustedFrom.y - minY} L ${adjustedTo.x - minX} ${adjustedTo.y - minY}`);
    path.setAttribute("stroke", "black");
    path.setAttribute("stroke-width", "2");
    path.setAttribute("fill", "none");
    path.setAttribute("marker-end", "url(#arrowhead)");

    svg.appendChild(path);
    document.body.appendChild(svg);
}

function drawArrows() {
    if (typeof dependencies !== "undefined") {
        dependencies.forEach(dep => createArrow(dep.from, dep.to));
    }
}

drawArrows();


// Redraw the arrows

function redrawArrows() {
    document.querySelectorAll("svg.dependency-arrow").forEach(svg => svg.remove())
    dependencies.forEach(dep => createArrow(dep.from, dep.to))
}

window.addEventListener("resize", () => {redrawArrows()})

window.addEventListener("scroll", () => {redrawArrows()})

const sections = document.querySelectorAll("section")

sections.forEach(section => {
    section.addEventListener("scroll", () => {redrawArrows()});
})

window.addEventListener("htmx:afterSwap", () => {
    setTimeout(() => {redrawArrows()}, 300)
})


