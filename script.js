"use strict";

const pageSizes = [
    {
        id: "letter",
        name: "US Letter",
        description: "8.5\" × 11\"",
        width: 216,
        height: 279
    },
    {
        id: "a4",
        name: "A4",
        description: "210mm × 297mm",
        width: 210,
        height: 297
    }
];

let currentPage = pageSizes[0];
let currentOrientation = "portrait";

function updatePage() {
    let svg = document.getElementById("svg");
    let rect = document.getElementById("rect");

    let width = currentPage.width;
    let height = currentPage.height;
    if (currentOrientation === "landscape") {
        width = currentPage.height;
        height = currentPage.width;
    }

    svg.setAttribute("viewBox", "0 0 " + width + " " + height);
    rect.setAttribute("width", width);
    rect.setAttribute("height", height);
}

function changePageSize(element) {
    currentPage = pageSizes.find(pageSize => pageSize.id === element.value);
    updatePage();
}

function changeOrientation(element) {
    switch (element.value) {
    case "portrait":
        currentOrientation = "portrait";
        updatePage();
        break;
    case "landscape":
        currentOrientation = "landscape";
        updatePage();
        break;
    default:
        console.log("Unknown orientation: " + element.value);
    }
}
