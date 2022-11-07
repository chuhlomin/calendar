"use strict";

let today = new Date();

let config = {
    pageSize: "A4",
    firstDay: "0",
    fontSizeDays: "12",
    fontSizeMonth: "12",
    fontSizeWeekdays: "5",
    fontSizeWeekNumbers: "5 ",
    textColor: "#222222",
    weekendColor: "#aa5555",
    year: today.getFullYear(),
    month: today.getMonth(),
    monthColor: "#222222",
    daysXStep: "25",
    daysXShift: "40",
    daysYStep: "35",
    daysYShift: "50",
    weeknumbersColor: "#999999",
    weeknumbersXShift: "20",
    weeknumbersYStep: "35",
    weeknumbersYShift: "42",
    weekdaysColor: "#999999",
    inactiveColor: "#c8c8c8",
};

let configInputTypes = {
    pageSize: "select",
    firstDay: "radio",
    fontSizeDays: "number",
    fontSizeMonth: "number",
    fontSizeWeekdays: "number",
    fontSizeWeekNumbers: "number",
    textColor: "color",
    weekendColor: "color",
    daysXStep: "number",
    daysXShift: "number",
    daysYStep: "number",
    daysYShift: "number",
    weeknumbersColor: "color",
    weeknumbersXShift: "number",
    weeknumbersYStep: "number",
    weeknumbersYShift: "number",
    monthColor: "color",
    weekdaysColor: "color",
    inactiveColor: "color",
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
        }
    }
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
    changeConfigKV(element.name, element.value);
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
            year: config.year,
            month: getMonthName(cfg.month),
            halfWidth: width/2,
            days: d,
            weekdays: weekdays(config.firstDay),
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
    if (end.getDay() === lastDay) {
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
        days.push({
            day: date.getDate(),
            x: column * cfg.daysXStep + cfg.daysXShift,
            y: row * cfg.daysYStep + cfg.daysYShift,
            inactive: date.getMonth() != cfg.month,
            weekend: date.getDay() == 0 || date.getDay() == 6,
        });

        if (date.getDay() == cfg.firstDay) {
            weeknumbers.push(weeknumber);
            weeknumber++;

            if (weeknumber > 52) {
                weeknumber = 1;
            }
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
     while (date.getDay() != cfg.firstDay) {
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

    for (let key in cfg) {
        if (configInputTypes[key] == "number" || configIntegerFields.includes(key)) {
            let value = parseInt(cfg[key]);
            if (isNaN(value)) {
                errors.push("Invalid value for " + key + ": " + cfg[key]);
                continue;
            }

            if (key == "month" && (value < 0 || value > 11)) {
                errors.push("Invalid value for " + key + ": " + cfg[key]);
                continue;
            }

            cfg[key] = value;
        }
    }

    return [cfg, errors];
}
