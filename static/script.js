"use strict";

let currentPageSize = "A4";
let currentOrientation = "P";
let currentPattern = "triangles";
let currentPatternWidth = "10";
let currentPatternHeight = "6";
let currentPatternColor = "#cccccc";
let currentLineWidth = "250";

let body = document.getElementsByTagName("body")[0];
let panel = document.getElementsByClassName('panel')[0];

window.onload = function() {
    // read state from localStorage
    let sizeID = localStorage.getItem("sizeID");
    if (sizeID) {
        currentPageSize = sizeID;
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
        updateForm(currentPattern);
    }

    let patternWidth = localStorage.getItem("patternWidth");
    if (patternWidth) {
        currentPatternWidth = patternWidth;
        let widthInput = document.getElementsByName("patternWidth")[0];
        widthInput.value = patternWidth;
    }

    let patternHeight = localStorage.getItem("patternHeight");
    if (patternHeight) {
        currentPatternHeight = patternHeight;
        let heightInput = document.getElementsByName("patternHeight")[0];
        heightInput.value = patternHeight;
    }

    let patternColor = localStorage.getItem("patternColor");
    if (patternColor) {
        currentPatternColor = patternColor;
    }

    let lineWidth = localStorage.getItem("patternLineWidth");
    if (lineWidth) {
        currentLineWidth = lineWidth;
        let lineWidthInput = document.getElementsByName("patternLineWidth")[0];
        lineWidthInput.value = lineWidth;
    }

    updatePatterns();

    // initialize color picker
    let button = document.getElementById('picker');
    let picker = new ColorPicker(button, currentPatternColor);
    button.addEventListener('colorChange', function (event) {
        changeColor(event.detail.color.hexa);
    });

    let patterns = document.getElementsByClassName("pattern");
    for (let patternEl of patterns) {
        patternEl.addEventListener("keypress", function(event) {
            if (event.code == "Space" || event.code == "Enter") {
                patternEl.childNodes[1].checked = true;
                changePattern(patternEl.childNodes[1]);
            }
        });
    }
};

let yOffset = 0;

panel.onscroll = function() {
    if (!body.classList.contains("preview")) {
        yOffset = panel.scrollTop;
    }

    if (window.innerHeight + yOffset >= panel.scrollHeight) {
        body.classList.add("scrolled");
    } else {
        body.classList.remove("scrolled");
    }
};

function updatePage() {
    let svg = document.getElementById("svg");
    let rect = document.getElementById("rect");
    let pattern = document.getElementById("pattern");

    let pageData = document.getElementById("pageSize").querySelector("option[value='" + currentPageSize + "']").dataset;

    let width = pageData.width;
    let height = pageData.height;
    if (currentOrientation === "L") {
        width = pageData.height;
        height = pageData.width;
    }

    svg.setAttribute("viewBox", "0 0 " + width + " " + height);
    rect.setAttribute("width", width);
    rect.setAttribute("height", height);
    pattern.setAttribute("width", width);
    pattern.setAttribute("height", height);
}

function changePageSize(element) {
    currentPageSize = element.value;
    updatePage();
    localStorage.setItem("sizeID", currentPageSize);
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

    if (localStorage.getItem("patternSizeChanged") != "true") {
        // switching from square-like to diaming-like patterns makes patterns
        // look wrong: too wide or too narrow
        // so if the user has not changed the pattern size, we change it
        // automatically to make it look better.
        // This is not perfect, but it's better than nothing.

        if (element.dataset.width) {
            currentPatternWidth = element.dataset.width;
            let widthInput = document.getElementsByName("patternWidth")[0];
            widthInput.value = currentPatternWidth;
        }

        if (element.dataset.height) {
            currentPatternHeight = element.dataset.height;
            let heightInput = document.getElementsByName("patternHeight")[0];
            heightInput.value = currentPatternHeight;
        }
        updatePatterns();
    }

    updatePattern(currentPattern);
    localStorage.setItem("pattern", currentPattern);
    updateForm(currentPattern);
}

function updateForm(pattern) {
    updateWidthElements(currentPattern != "lines");
}

function updateWidthElements(enabled) {
    const field = document.getElementById("widthField");
    const input = document.getElementById("widthInput");

    if (enabled) {
        field.classList.remove("disabled");
        input.removeAttribute("disabled");
    } else {
        field.classList.add("disabled");
        input.setAttribute("disabled", "disabled");
    }
}

function changePatternWidth(element) {
    currentPatternWidth = element.value;
    updatePatterns();
    localStorage.setItem("patternWidth", currentPatternWidth);
    localStorage.setItem("patternSizeChanged", true);
}

function changePatternHeight(element) {
    currentPatternHeight = element.value;
    updatePatterns();
    localStorage.setItem("patternHeight", currentPatternHeight);
    localStorage.setItem("patternSizeChanged", true);
}

function changeColor(color) {
    currentPatternColor = color;
    updatePatterns();
    localStorage.setItem("patternColor", currentPatternColor);
}

function changePatternLineWidth(element) {
    currentLineWidth = element.value;
    updatePatterns();
    localStorage.setItem("patternLineWidth", currentLineWidth);
}

function updatePattern(patternName) {
    let pattern = document.getElementById("pattern");
    pattern.setAttribute("fill", "url(#" + patternName + ")");
}

function updatePatterns() {
    let template = document.getElementById('template_patterns').innerHTML;
    let rendered = Mustache.render(
        template,
        {
            width: currentPatternWidth,
            widthHalf: currentPatternWidth / 2,
            height: currentPatternHeight,
            heightHalf: currentPatternHeight / 2,
            color: currentPatternColor,
            lineWidth: currentLineWidth / 1000,
        }
    );
    document.getElementById('patterns').innerHTML = rendered;
}

function togglePreview(element) {
    if (body.classList.contains("preview")) {
        body.classList.remove("preview");
        panel.scrollTop = yOffset;
    } else {
        body.classList.add("preview");
    }
}

function submitForm(element) {
    let defaultLabel = element.value;
    element.setAttribute("disabled", "disabled");
    element.value = "Generating PDF...";

    fetch("/pdf", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify({
            size: currentPageSize,
            orientation: currentOrientation,
            pattern: {
                name: currentPattern,
                width: currentPatternWidth,
                height: currentPatternHeight,
                color: currentPatternColor,
                lineWidth: currentLineWidth,
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
            element.removeAttribute("disabled");
            element.value = defaultLabel;
        } else {
            response.text().then(text => {
                alert(text);
                console.log(text);
            });

            element.removeAttribute("disabled");
            element.value = defaultLabel;
        }
    }).catch(error => {
        alert("Error: " + error);
        element.removeAttribute("disabled");
        element.value = defaultLabel;
    });
}
