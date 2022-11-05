"use strict";

let currentPageSize = "A4";
let currentOrientation = "P";
let currentFirstDay = "0";
let currentTextColor = "#000000";

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

    let firstDay = localStorage.getItem("firstDay");
    if (firstDay) {
        currentFirstDay = firstDay;
        let firstDaySelect = document.querySelector("input[name='firstDay'][value='" + firstDay + "']");
        firstDaySelect.checked = true;
    }

    let textColor = localStorage.getItem("textColor");
    if (textColor) {
        currentTextColor = textColor;
    }

    // let lineWidth = localStorage.getItem("patternLineWidth");
    // if (lineWidth) {
    //     currentLineWidth = lineWidth;
    //     let lineWidthInput = document.getElementsByName("patternLineWidth")[0];
    //     lineWidthInput.value = lineWidth;
    // }

    updateCalendar();

    // initialize color picker
    let button = document.getElementById('picker');
    let picker = new ColorPicker(button, currentTextColor);
    button.addEventListener('colorChange', function (event) {
        changeTextColor(event.detail.color.hexa);
    });
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
    // pattern.setAttribute("width", width);
    // pattern.setAttribute("height", height);
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

function changeFirstDay(element) {
    currentFirstDay = element.value;
    updateCalendar();
    localStorage.setItem("firstDay", currentFirstDay);
}

function changeTextColor(color) {
    currentTextColor = color;
    updateCalendar();
    localStorage.setItem("textColor", currentTextColor);
}

function updateCalendar() {
    let pageData = document.getElementById("pageSize").querySelector("option[value='" + currentPageSize + "']").dataset;
    let width = pageData.width;

    const [d, w] = days(2022, 11, currentFirstDay);

    let templateMonth = document.getElementById('template_month').innerHTML;
    let renderedMonth = Mustache.render(
        templateMonth,
        {
            year: "2022",
            month: "November",
            halfWidth: width/2,
            days: d,
            weekdays: weekdays(currentFirstDay),
            weeknumbers: weeknumbers(w)
        }
    );

    let templateStyles = document.getElementById('template_styles').innerHTML;
    let renderedStyles = Mustache.render(
        templateStyles,
        {
            textColor: currentTextColor,
        }
    );

    document.getElementById('defs').innerHTML = renderedMonth;
    document.getElementById('style').innerHTML = renderedStyles;
}

function days(year, month, firstDay) {
    firstDay = parseInt(firstDay);
    let days = [];
    let date = new Date(year, month - 1, 1);
    let end = new Date(year, month, 0);

    if (date.getDay() != firstDay) {
        date.setDate(date.getDate() - (date.getDay() - firstDay));
    }

    let row = 0;
    let column = 0;
    let weeknumbers = [];

    let weeknumber = getWeekNumber(date);
    weeknumbers.push(weeknumber);

    while (date <= end) {
        days.push({
            day: date.getDate(),
            x: column * 25 + 40,
            y: row * 35 + 50,
            inactive: date.getMonth() != month - 1,
            weekend: date.getDay() == 0 || date.getDay() == 6,
        });

        date.setDate(date.getDate() + 1);

        if (column == 6) {
            row++;
            column = 0;
            weeknumber++;
            weeknumbers.push(weeknumber);
        } else {
            column++;
        }
    }

    // finish last row
    while (date.getDay() != firstDay) {
        days.push({
            day: date.getDate(),
            x: column * 25 + 40,
            y: row * 35 + 50,
            inactive: true,
        });

        date.setDate(date.getDate() + 1);
        column++;
    }

    return [days, weeknumbers];
}

function weekdays(firstDay) {
    let weekdays = [
        // "Sunday",
        // "Monday",
        // "Tuesday",
        // "Wednesday",
        // "Thursday",
        // "Friday",
        // "Saturday",
        "Sun",
        "Mon",
        "Tue",
        "Wed",
        "Thu",
        "Fri",
        "Sat",
    ];

    // parse int
    firstDay = parseInt(firstDay);
    
    let days = [];
    for (let i = 0; i < 7; i++) {
        days.push({
            day: weekdays[(firstDay + i) % 7],
            x: i * 25 + 40,
            y: 30,
        });
    }

    return days;
}

function weeknumbers(w) {
    let result = [];

    for (let i = 0; i < w.length; i++) {
        result.push({
            weeknumber: w[i],
            x: 20,
            y: i * 35 + 42,
        });
    }

    return result;
}

function getWeekNumber(d) {
    d = new Date(Date.UTC(d.getFullYear(), d.getMonth(), d.getDate()));
    d.setUTCDate(d.getUTCDate() + 4 - (d.getUTCDay()||7));
    let yearStart = new Date(Date.UTC(d.getUTCFullYear(),0,1));
    return Math.ceil((((d - yearStart) / 86400000) + 1)/7);
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
        })
    }).then(response => {
        if (response.ok) {
            response.blob().then(blob => {
                let url = URL.createObjectURL(blob);
                let a = document.createElement("a");
                a.href = url;
                a.download = "calendar.pdf";
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
