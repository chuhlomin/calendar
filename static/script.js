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
}

function changeOrientation(element) {
    switch (element.value) {
    case "portrait":
        currentOrientation = "P";
        updatePage();
        break;
    case "landscape":
        currentOrientation = "L";
        updatePage();
        break;
    default:
        console.log("Unknown orientation: " + element.value);
    }
}

function changePattern(element) {
    let pattern = document.getElementById("pattern");
    pattern.setAttribute("fill", "url(#" + element.value + ")");
    currentPattern = element.value;
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
            pattern: currentPattern
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
