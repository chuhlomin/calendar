"use strict";

let today = new Date();

let config = {
    pageSize: "A4",
    firstDay: "0",
    year: today.getFullYear(),
    month: today.getMonth(),

    // days
    fontSizeDays: "12",
    textColor: "#222222",
    weekendColor: "#aa5555",
    daysXShift: "40",
    daysXStep: "25",
    daysYShift: "50",
    daysYStep: "35",
    showInactiveDays: "true",
    inactiveColor: "#c8c8c8",

    // month
    showMonth: "true",
    fontSizeMonth: "12",
    monthColor: "#222222",

    // weekdays
    showWeekdays: "true",
    fontSizeWeekdays: "5",
    weekdaysColor: "#999999",
    weekdaysXShift: "40",
    weekdaysXStep: "25",
    weekdaysYShift: "30",

    // week numbers
    showWeekNumbers: "true",
    fontSizeWeekNumbers: "5 ",
    weeknumbersColor: "#999999",
    weeknumbersXShift: "20",
    weeknumbersYShift: "42",
    weeknumbersYStep: "35",
};

let configInputTypes = {
    pageSize: "select",
    firstDay: "radio",

    // days
    fontSizeDays: "number",
    textColor: "color",
    weekendColor: "color",
    daysXShift: "number",
    daysXStep: "number",
    daysYShift: "number",
    daysYStep: "number",
    showInactiveDays: "checkbox",
    inactiveColor: "color",

    // month
    showMonth: "checkbox",
    fontSizeMonth: "number",
    monthColor: "color",
    
    // weekdays
    showWeekdays: "checkbox",
    fontSizeWeekdays: "number",
    weekdaysColor: "color",
    weekdaysXShift: "number",
    weekdaysXStep: "number",
    weekdaysYShift: "number",

    // week numbers
    showWeekNumbers: "checkbox",
    fontSizeWeekNumbers: "number",
    weeknumbersColor: "color",
    weeknumbersXShift: "number",
    weeknumbersYShift: "number",
    weeknumbersYStep: "number",
};

let configIntegerFields = ["firstDay", "year", "month"];

let body = document.getElementsByTagName("body")[0];
let panel = document.getElementsByClassName('panel')[0];

function loadConfig(key) {
    var value = localStorage.getItem(key);
    if (value) {
        config[key] = value;

        switch (configInputTypes[key]) {
        case "select":
            let select = document.getElementsByName(key)[0];
            select.value = value;

            if (key == "pageSize") {
                updatePage();
            }
            break;

        case "radio":
            let option = document.querySelector("input[name='" + key + "'][value='" + value + "']");
            option.checked = true;
            break;

        case "number":
            let input = document.getElementsByName(key)[0];
            input.value = value;
            break;

        case "checkbox":
            let checkbox = document.getElementsByName(key)[0];
            checkbox.checked = value == "true";
            break;
        }
    }
}

function validateConfig(config) {
    let errors = [];

    let cfg = {};

    for (let key in config) {
        if (configInputTypes[key] == "number" || configIntegerFields.includes(key)) {
            let value = parseInt(config[key]);
            if (isNaN(value)) {
                errors.push("Invalid value for " + key + ": " + config[key]);
                continue;
            }

            if (key == "month" && (value < 0 || value > 11)) {
                errors.push("Invalid value for " + key + ": " + config[key]);
                continue;
            }

            cfg[key] = value;
            continue;
        }

        if (configInputTypes[key] == "checkbox") {
            cfg[key] = config[key] == "true";
            continue;
        }

        cfg[key] = config[key];
    }

    return [cfg, errors];
}

window.onload = function() {
    for (let key in config) {
        loadConfig(key);

        if (configInputTypes[key] == "color") {
            let button = document.getElementsByName(key)[0];
            new ColorPicker(button, config[key]);
            button.addEventListener('colorChange', function (event) {
                changeConfigKV(key, event.detail.color.hexa);
            });
        }
    }

    updateCalendar();
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

function changeConfig(element) {
    let key = element.name;
    let value = element.value;

    if (configInputTypes[key] == "checkbox") {
        if (element.checked) {
            value = "true";
        } else {
            value = "false";
        }
    }

    changeConfigKV(key, value);
}

function changeConfigKV(key, value) {
    config[key] = value;
    updateCalendar();
    localStorage.setItem(key, value);
}

function updateCalendar() {
    let pageData = document.getElementById("pageSize").querySelector("option[value='" + config.pageSize + "']").dataset;
    let width = pageData.width;

    let [cfg, errors] = validateConfig(config);
    if (errors.length > 0) {
        console.error(errors);
        return;
    }

    const [d, w] = days(cfg);

    let templateMonth = document.getElementById('template_month').innerHTML;
    let renderedMonth = Mustache.render(
        templateMonth,
        {
            month: getMonth(cfg),
            halfWidth: width/2,
            days: d,
            weekdays: weekdays(cfg),
            weeknumbers: weeknumbers(cfg, w)
        }
    );

    let templateStyles = document.getElementById('template_styles').innerHTML;
    let renderedStyles = Mustache.render(templateStyles, config);

    document.getElementById('defs').innerHTML = renderedMonth;
    document.getElementById('style').innerHTML = renderedStyles;
}

function days(cfg) {
    let lastDay = cfg.firstDay - 1;
    if (lastDay < 0) {
        lastDay = 6;
    }

    let days = [];
    let date = new Date(cfg.year, cfg.month, 1);
    let end = new Date(cfg.year, cfg.month + 1, 0);
    if (end.getDay() === lastDay && cfg.showInactiveDays) {
        end.setDate(end.getDate() + 7);
    }

    // add days from previous month
    if (date.getDay() != cfg.firstDay) {
        if (date.getDay() == lastDay) {
            date.setDate(date.getDate() - 6);
        } else {
            date.setDate(date.getDate() - (date.getDay() - cfg.firstDay));
        }
    }

    let row = 0;
    let column = 0;
    let weeknumbers = [];

    let weeknumber = getWeekNumber(date);

    while (date <= end) {
        let inactive = date.getMonth() != cfg.month;

        if (date.getDay() == cfg.firstDay && cfg.showWeekNumbers) {
            weeknumbers.push(weeknumber);
            weeknumber++;

            if (weeknumber > 52) {
                weeknumber = 1;
            }
        }

        if (inactive && !cfg.showInactiveDays) {
            date.setDate(date.getDate() + 1);
            column++;
            continue;
        }

        days.push({
            day: date.getDate(),
            x: column * cfg.daysXStep + cfg.daysXShift,
            y: row * cfg.daysYStep + cfg.daysYShift,
            inactive: inactive,
            weekend: date.getDay() == 0 || date.getDay() == 6,
        });

        date.setDate(date.getDate() + 1);

        if (column == 6) {
            row++;
            column = 0;
        } else {
            column++;
        }
    }

    // finish last row
     while (date.getDay() != cfg.firstDay && cfg.showInactiveDays) {
        days.push({
            day: date.getDate(),
            x: column * cfg.daysXStep + cfg.daysXShift,
            y: row * cfg.daysYStep + cfg.daysYShift,
            inactive: true,
        });

        date.setDate(date.getDate() + 1);
        column++;
    }

    return [days, weeknumbers];
}

function weekdays(cfg) {
    if (!cfg.showWeekdays) {
        return [];
    }

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

    let days = [];
    for (let i = 0; i < 7; i++) {
        days.push({
            day: weekdays[(cfg.firstDay + i) % 7],
            x: i * cfg.weekdaysXStep + cfg.weekdaysXShift,
            y: cfg.weekdaysYShift,
        });
    }

    return days;
}

function weeknumbers(cfg, w) {
    let result = [];

    for (let i = 0; i < w.length; i++) {
        result.push({
            weeknumber: w[i],
            x: cfg.weeknumbersXShift,
            y: i * cfg.weeknumbersYStep + cfg.weeknumbersYShift,
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

function getMonth(cfg) {
    if (!cfg.showMonth) {
        return "";
    }

    return months[cfg.month] + " " + cfg.year;
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
