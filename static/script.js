"use strict";

let today = new Date();

let config = {
    pageSize: "A4",
    firstDay: "0",
    textColor: "#222222",
    weekendColor: "#aa5555",
    year: today.getFullYear(),
    month: today.getMonth(),
    daysXStep: "25",
    daysXShift: "40",
    daysYStep: "35",
    daysYShift: "50",
};


let body = document.getElementsByTagName("body")[0];
let panel = document.getElementsByClassName('panel')[0];

window.onload = function() {
    // read state from localStorage
    let sizeID = localStorage.getItem("sizeID");
    if (sizeID) {
        config["pageSize"] = sizeID;
        let sizeSelect = document.getElementsByName("size")[0];
        sizeSelect.value = sizeID;
    }

    if (sizeID) {
        updatePage();
    }

    let firstDay = localStorage.getItem("firstDay");
    if (firstDay) {
        config["firstDay"] = firstDay;
        let firstDaySelect = document.querySelector("input[name='firstDay'][value='" + firstDay + "']");
        firstDaySelect.checked = true;
    }

    let textColor = localStorage.getItem("textColor");
    if (textColor) {
        config["textColor"] = textColor;
    }

    let weekendColor = localStorage.getItem("weekendColor");
    if (weekendColor) {
        config["weekendColor"] = weekendColor;
    }

    let year = localStorage.getItem("year");
    if (year) {
        config["year"] = year;
    }

    let month = localStorage.getItem("month");
    if (month) {
        config["month"] = month;
    }

    let daysXStep = localStorage.getItem("daysXStep");
    if (daysXStep) {
        config["daysXStep"] = daysXStep;
        let daysXStepInput = document.getElementsByName("daysXStep")[0];
        daysXStepInput.value = daysXStep;
    }

    let daysXShift = localStorage.getItem("daysXShift");
    if (daysXShift) {
        config["daysXShift"] = daysXShift;
        let daysXShiftInput = document.getElementsByName("daysXShift")[0];
        daysXShiftInput.value = daysXShift;
    }

    let daysYStep = localStorage.getItem("daysYStep");
    if (daysYStep) {
        config["daysYStep"] = daysYStep;
        let daysYStepInput = document.getElementsByName("daysYStep")[0];
        daysYStepInput.value = daysYStep;
    }

    let daysYShift = localStorage.getItem("daysYShift");
    if (daysYShift) {
        config["daysYShift"] = daysYShift;
        let daysYShiftInput = document.getElementsByName("daysYShift")[0];
        daysYShiftInput.value = daysYShift;
    }

    updateCalendar();

    // initialize color pickers
    let pickerTextColor = document.getElementById('pickerTextColor');
    new ColorPicker(pickerTextColor, config.textColor);
    pickerTextColor.addEventListener('colorChange', function (event) {
        changeTextColor(event.detail.color.hexa);
    });

    let pickerWeekendColor = document.getElementById('pickerWeekendColor');
    new ColorPicker(pickerWeekendColor, config.weekendColor);
    pickerWeekendColor.addEventListener('colorChange', function (event) {
        changeWeekendColor(event.detail.color.hexa);
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

    let pageData = document.getElementById("pageSize").querySelector("option[value='" + config.pageSize + "']").dataset;

    let width = pageData.width;
    let height = pageData.height;

    svg.setAttribute("viewBox", "0 0 " + width + " " + height);
    rect.setAttribute("width", width);
    rect.setAttribute("height", height);
}

function changePageSize(element) {
    config["pageSize"] = element.value;
    updatePage();
    localStorage.setItem("sizeID", config.pageSize);
}

function changeFirstDay(element) {
    config["firstDay"] = element.value;
    updateCalendar();
    localStorage.setItem("firstDay", config.firstDay);
}

function changeTextColor(color) {
    config["textColor"] = color;
    updateCalendar();
    localStorage.setItem("textColor", config.textColor);
}

function changeWeekendColor(color) {
    config["weekendColor"] = color;
    updateCalendar();
    localStorage.setItem("weekendColor", config.weekendColor);
}

function changeMonth(step) {
    let month = parseInt(config.month);
    let year = parseInt(config.year);

    month += step;
    if (month < 0) {
        month = 11;
        year--;
    } else if (month > 11) {
        month = 0;
        year++;
    }

    config.month = month;
    config.year = year;

    updateCalendar();

    localStorage.setItem("month", config.month);
    localStorage.setItem("year", config.year);
}

function changeDaysXStep(element) {
    config["daysXStep"] = element.value;
    updateCalendar();
    localStorage.setItem("daysXStep", config.daysXStep);
}

function changeDaysXShift(element) {
    config["daysXShift"] = element.value;
    updateCalendar();
    localStorage.setItem("daysXShift", config.daysXShift);
}

function changeDaysYStep(element) {
    config["daysYStep"] = element.value;
    updateCalendar();
    localStorage.setItem("daysYStep", config.daysYStep);
}

function changeDaysYShift(element) {
    config["daysYShift"] = element.value;
    updateCalendar();
    localStorage.setItem("daysYShift", config.daysYShift);
}

function updateCalendar() {
    let pageData = document.getElementById("pageSize").querySelector("option[value='" + config.pageSize + "']").dataset;
    let width = pageData.width;

    // validate year
    let year = parseInt(config.year);
    if (isNaN(year)) {
        console.log("Invalid year: " + year);
        return;
    }

    // validate month
    let month = parseInt(config.month);
    if (isNaN(month)) {
        console.log("Invalid month: " + month);
        return;
    }
    if (month < 0 || month > 11) {
        console.log("Invalid month: " + month);
        return;
    }

    let firstDay = parseInt(config.firstDay);
    let lastDay = firstDay - 1;
    if (lastDay < 0) {
        lastDay = 6;
    }

    let daysXShift = parseInt(config.daysXShift);
    if (isNaN(daysXShift)) {
        console.log("Invalid daysXShift: " + daysXShift);
        return;
    }

    let daysXStep = parseInt(config.daysXStep);
    if (isNaN(daysXStep)) {
        console.log("Invalid daysXStep: " + daysXStep);
        return;
    }

    let daysYShift = parseInt(config.daysYShift);
    if (isNaN(daysYShift)) {
        console.log("Invalid daysYShift: " + daysYShift);
        return;
    }

    let daysYStep = parseInt(config.daysYStep);
    if (isNaN(daysYStep)) {
        console.log("Invalid daysYStep: " + daysYStep);
        return;
    }

    const [d, w] = days(
        year,
        month,
        firstDay,
        lastDay,
        daysXStep,
        daysXShift,
        daysYStep,
        daysYShift
    );

    let templateMonth = document.getElementById('template_month').innerHTML;
    let renderedMonth = Mustache.render(
        templateMonth,
        {
            year: config.year,
            month: getMonthName(month),
            halfWidth: width/2,
            days: d,
            weekdays: weekdays(config.firstDay),
            weeknumbers: weeknumbers(w)
        }
    );

    let templateStyles = document.getElementById('template_styles').innerHTML;
    let renderedStyles = Mustache.render(
        templateStyles,
        {
            textColor: config.textColor,
            weekendColor: config.weekendColor
        }
    );

    document.getElementById('defs').innerHTML = renderedMonth;
    document.getElementById('style').innerHTML = renderedStyles;
}

function days(year, month, firstDay, lastDay, daysXStep, daysXShift, daysYStep, daysYShift) {
    let days = [];
    let date = new Date(year, month, 1);
    let end = new Date(year, month + 1, 0);
    if (end.getDay() === lastDay) {
        end.setDate(end.getDate() + 7);
    }

    // add days from previous month
    if (date.getDay() != firstDay) {
        if (date.getDay() == lastDay) {
            date.setDate(date.getDate() - 6);
        } else {
            date.setDate(date.getDate() - (date.getDay() - firstDay));
        }
    }

    let row = 0;
    let column = 0;
    let weeknumbers = [];

    let weeknumber = getWeekNumber(date);

    while (date <= end) {
        days.push({
            day: date.getDate(),
            x: column * daysXStep + daysXShift,
            y: row * daysYStep + daysYShift,
            inactive: date.getMonth() != month,
            weekend: date.getDay() == 0 || date.getDay() == 6,
        });

        if (date.getDay() == firstDay) {
            weeknumbers.push(weeknumber);
            weeknumber++;
        }

        date.setDate(date.getDate() + 1);

        if (column == 6) {
            row++;
            column = 0;
        } else {
            column++;
        }
    }

    // finish last row
     while (date.getDay() != firstDay) {
        days.push({
            day: date.getDate(),
            x: column * daysXStep + daysXShift,
            y: row * daysYStep + daysYShift,
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

let months = [
    "January",
    "February",
    "March",
    "April",
    "May",
    "June",
    "July",
    "August",
    "September",
    "October",
    "November",
    "December",
];

function getMonthName(month) {
    return months[month];
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

    let [cfg, errors] = validateConfig(config);
    if (errors.length > 0) {
        alert(errors.join("\r"));
        return;
    }

    fetch("/pdf", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(cfg)
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

function validateConfig(cfg) {
    let errors = [];

    let year = parseInt(cfg.year);
    if (isNaN(year)) {
        errors.push("Year is not a number");
    }

    let month = parseInt(cfg.month);
    if (isNaN(month)) {
        errors.push("Month is not a number");
    }

    let firstDay = parseInt(cfg.firstDay);
    if (isNaN(firstDay)) {
        errors.push("First day is not a number");
    }

    let daysXStep = parseInt(cfg.daysXStep);
    if (isNaN(daysXStep)) {
        errors.push("Days X step is not a number");
    }

    let daysXShift = parseInt(cfg.daysXShift);
    if (isNaN(daysXShift)) {
        errors.push("Days X shift is not a number");
    }

    let daysYStep = parseInt(cfg.daysYStep);
    if (isNaN(daysYStep)) {
        errors.push("Days Y step is not a number");
    }

    let daysYShift = parseInt(cfg.daysYShift);
    if (isNaN(daysYShift)) {
        errors.push("Days Y shift is not a number");
    }

    cfg.year = year;
    cfg.month = month;
    cfg.firstDay = firstDay;
    cfg.daysXStep = daysXStep;
    cfg.daysXShift = daysXShift;
    cfg.daysYStep = daysYStep;
    cfg.daysYShift = daysYShift;

    return [cfg, errors];
}
