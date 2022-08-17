"use strict";

const pageSizes = [
    {
        id: "Letter",
        name: "US Letter",
        description: "8.5\" × 11\"",
        width: 216,
        height: 279
    },
    {
        id: "A4",
        name: "A4",
        description: "210mm × 297mm",
        width: 210,
        height: 297
    }
];

let currentPage = pageSizes[0];
let currentOrientation = "P";
let currentPattern = "rect";
let currentPatternSize = "15";
let currentPatternColor = "#000000";

window.onload = function() {
    // read state from localStorage
    let sizeID = localStorage.getItem("sizeID");
    if (sizeID) {
        currentPage = pageSizes.find(size => size.id == sizeID);
        let sizeSelect = document.getElementsByName("size")[0];
        sizeSelect.value = sizeID;
    }
    let orientation = localStorage.getItem("orientation");
    if (orientation) {
        currentOrientation = orientation;
        let orientationSelect = document.querySelector("input[name='orientation'][value='" + orientation + "']");
        orientationSelect.checked = true;
    }

    if (sizeID || orientation) {
        updatePage();
    }

    let pattern = localStorage.getItem("pattern");
    if (pattern) {
        currentPattern = pattern;
        let patternSelect = document.querySelector("input[name='pattern'][value='" + pattern + "']");
        patternSelect.checked = true;
        updatePattern(currentPattern);
    }

    let patternSize = localStorage.getItem("patternSize");
    if (patternSize) {
        currentPatternSize = patternSize;
        let sizeInput = document.getElementsByName("patternSize")[0];
        sizeInput.value = patternSize;
        updatePatternSize(currentPatternSize);
    }

    let patternColor = localStorage.getItem("patternColor");
    if (patternColor) {
        currentPatternColor = patternColor;
        updatePatternColor(currentPatternColor);
    }

    // initialize color picker
    let button = document.getElementById('picker');
    let picker = new ColorPicker(button, currentPatternColor);
    button.addEventListener('colorChange', function (event) {
        changeColor(event.detail.color.hexa);
    });
}

function updatePage() {
    let svg = document.getElementById("svg");
    let rect = document.getElementById("rect");
    let pattern = document.getElementById("pattern");

    let width = currentPage.width;
    let height = currentPage.height;
    if (currentOrientation === "L") {
        width = currentPage.height;
        height = currentPage.width;
    }

    svg.setAttribute("viewBox", "0 0 " + width + " " + height);
    rect.setAttribute("width", width);
    rect.setAttribute("height", height);
    pattern.setAttribute("width", width);
    pattern.setAttribute("height", height);
}

function changePageSize(element) {
    currentPage = pageSizes.find(pageSize => pageSize.id === element.value);
    updatePage();
    localStorage.setItem("sizeID", currentPage.id);
}

function changeOrientation(element) {
    if (element.value != "P" && element.value != "L") {
        console.log("Unknown orientation: " + element.value);
        return;
    }

    currentOrientation = element.value;
    updatePage();
    localStorage.setItem("orientation", element.value);
}

function changePattern(element) {
    currentPattern = element.value;
    updatePattern(currentPattern);
    localStorage.setItem("pattern", currentPattern);
}

function changePatternSize(element) {
    currentPatternSize = element.value;
    updatePatternSize(currentPatternSize);
    localStorage.setItem("patternSize", currentPatternSize);
}

function changeColor(color) {
    currentPatternColor = color;
    updatePatternColor(currentPatternColor);
    localStorage.setItem("patternColor", currentPatternColor);
}

function updatePattern(patternName) {
    let pattern = document.getElementById("pattern");
    pattern.setAttribute("fill", "url(#" + patternName + ")");
}

function updatePatternSize(size) {
    let patterns = document.getElementsByTagName("pattern");

    for (let pattern of patterns) {
        let width = size;
        let height = size;

        if (pattern.id == "rhombus") {
            height = width * 0.6;
        }

        pattern.setAttribute("width", width);
        pattern.setAttribute("height", height);
    }
}

function updatePatternColor(color) {
    let svg = document.getElementById("svg");
    const paths = svg.getElementsByTagName("path");
    for (let path of paths) {
        path.setAttribute("stroke", color);
    }
    const circles = svg.getElementsByTagName("circle")
    for (let circle of circles) {
        circle.setAttribute("fill", color);
    }
}

function submitForm() {
    fetch("/pdf", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            size: currentPage.id,
            orientation: currentOrientation,
            pattern: {
                name: currentPattern,
                size: currentPatternSize,
                color: currentPatternColor
            }
        })
    }).then(response => {
        if (response.ok) {
            response.blob().then(blob => {
                let url = URL.createObjectURL(blob);
                let a = document.createElement("a");
                a.href = url;
                a.download = "grid.pdf";
                a.click();
            });
        } else {
            console.log("Error: " + response.statusText);
        }
    }).catch(error => {
        console.log("Error: " + error);
    });
}
