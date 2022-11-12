"use strict";

let today = new Date();

let config = {
    pageSize: "A4",
    firstDay: "0",
    year: today.getFullYear(),
    month: today.getMonth(), // changed by loadConfig if not found in localStorage

    // days
    daysFontSize: "45",
    daysFontFamily: "iosevka-regular",
    textColor: "#222222",
    weekendColor: "#aa5555",
    daysX: "40",
    daysY: "50",
    daysXStep: "25",
    daysYStep: "32",
    showInactiveDays: "true",
    inactiveColor: "#c8c8c8",

    // month
    showMonth: "true",
    monthFontFamily: "iosevka-aile-regular",
    monthFontSize: "50",
    monthColor: "#222222",
    monthY: "260",

    // weekdays
    showWeekdays: "true",
    weekdaysFontFamily: "iosevka-aile-regular",
    weekdaysFontSize: "18",
    weekdaysColor: "#999999",
    weekdaysX: "40",
    weekdaysY: "32",

    // week numbers
    showWeekNumbers: "true",
    weeknumbersFontFamily: "iosevka-regular",
    weeknumbersFontSize: "16",
    weeknumbersColor: "#999999",
    weeknumbersX: "20",
    weeknumbersY: "43",
};

let configInputTypes = {
    pageSize: "select",
    firstDay: "radio",

    // days
    daysFontSize: "number",
    daysFontFamily: "select",
    textColor: "color",
    weekendColor: "color",
    daysX: "number",
    daysY: "number",
    daysXStep: "number",
    daysYStep: "number",
    showInactiveDays: "checkbox",
    inactiveColor: "color",

    // month
    showMonth: "checkbox",
    monthFontFamily: "select",
    monthFontSize: "number",
    monthColor: "color",
    monthY: "number",
    
    // weekdays
    showWeekdays: "checkbox",
    weekdaysFontFamily: "select",
    weekdaysFontSize: "number",
    weekdaysColor: "color",
    weekdaysX: "number",
    weekdaysY: "number",

    // week numbers
    showWeekNumbers: "checkbox",
    weeknumbersFontFamily: "select",
    weeknumbersFontSize: "number",
    weeknumbersColor: "color",
    weeknumbersX: "number",
    weeknumbersY: "number",
};

let configIntegerFields = ["firstDay", "year", "month"];

let body = document.getElementsByTagName("body")[0];
let panel = document.getElementsByClassName('panel')[0];

let availableFonts = {
    "iosevka-regular": "Iosevka",
    "iosevka-aile-regular": "Iosevka Aile",
    "iosevka-etoile-regular": "Iosevka Etoile",
    "iosevka-curly-regular": "Iosevka Curly",
    "iosevka-curly-slab-regular": "Iosevka Curly Slab",
    "iosevka-ss01-regular": "Iosevka SS01 (Andale Mono Style)",
    "iosevka-ss02-regular": "Iosevka SS02 (Anonymous Pro Style)",
    "iosevka-ss03-regular": "Iosevka SS03 (Consolas Style)",
    "iosevka-ss04-regular": "Iosevka SS04 (Menlo Style)",
    "iosevka-ss05-regular": "Iosevka SS05 (Fira Mono Style)",
    "iosevka-ss06-regular": "Iosevka SS06 (Liberation Mono Style)",
    "iosevka-ss07-regular": "Iosevka SS07 (Monaco Style)",
    "iosevka-ss08-regular": "Iosevka SS08 (Pragmata Pro Style)",
    "iosevka-ss09-regular": "Iosevka SS09 (Source Code Pro Style)",
    "iosevka-ss10-regular": "Iosevka SS10 (Envy Code R Style)",
    "iosevka-ss11-regular": "Iosevka SS11 (X Windows Style)",
    "iosevka-ss12-regular": "Iosevka SS12 (Ubuntu Mono Style)",
    "iosevka-ss13-regular": "Iosevka SS13 (Lucida Style)",
    "iosevka-ss14-regular": "Iosevka SS14 (JetBrains Mono Style)",
    "iosevka-ss15-regular": "Iosevka SS15 (IMB Plex Mono Style)",
    "iosevka-ss16-regular": "Iosevka SS16 (PT Mono Style)",
    "iosevka-ss17-regular": "Iosevka SS17 (Recursive Mono Style)",
    "iosevka-ss18-regular": "Iosevka SS18 (Input Mono Style)",
};
let loadedFonts = {};
let loadingFonts = {};

function loadFont(fontName) {
    if (loadedFonts[fontName] || loadingFonts[fontName]) {
        return;
    }

    loadingFonts[fontName] = true;

    const font = new FontFace(fontName, "url('/" + fontName + ".ttf') format('truetype'), url('/" + fontName + ".woff') format('woff2')");
    font.load()
        .then(function(loadedFace) {
            document.fonts.add(loadedFace);
            loadedFonts[fontName] = true;
            delete loadingFonts[fontName];
            updateCalendar();
        })
        .catch(function(error) {
            console.log(error);
        });
}

function loadConfig(key) {
    var value = localStorage.getItem(key);
    if (value) {
        config[key] = value;

        switch (configInputTypes[key]) {
        case "select":
            let select = document.getElementsByName(key)[0];
            select.value = value;
            if (select.classList.contains("font")) {
                loadFont(value);
            }

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
            if (checkbox.dataset.fieldset) {
                toggleFieldset(checkbox.dataset.fieldset, checkbox.checked);
            }
            break;
        }
    } else { // value not found in localStorage
        if (configInputTypes[key] == "select") {
            let select = document.getElementsByName(key)[0];
            if (select.classList.contains("font")) {
                loadFont(config[key]);
            }
        }
        if (key == "year") {
            // set year to current year
            let input = document.getElementsByName(key)[0];
            input.value = config.year;
        }
        if (key == "month") {
            // make sure the month by default is the next month
            let date = new Date();
            config[key] = date.getMonth() + 1;
            if (config[key] == 12) {
                config.year++;
            }
            let select = document.getElementsByName(key)[0];
            select.value = config[key];
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
    // restore panel scroll position
    let panelOffset = localStorage.getItem("panelOffset");
    if (panelOffset) {
        panel.scrollTop = panelOffset;
    }

    // update select.font options
    let fontSelects = document.querySelectorAll("select.font");
    for (let i = 0; i < fontSelects.length; i++) {
        let select = fontSelects[i];
        // remove all options
        while (select.firstChild) {
            select.removeChild(select.firstChild);
        }
        // add options
        for (let font in availableFonts) {
            let option = document.createElement("option");
            option.value = font;
            option.text = availableFonts[font];
            select.add(option);
        }
    }

    // load config from local storage
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
let timerForYOffset = null;

panel.onscroll = function() {
    if (!body.classList.contains("preview")) {
        yOffset = panel.scrollTop;
    }

    // start timer if not already running
    if (timerForYOffset) {
        clearTimeout(timerForYOffset);
    }
    timerForYOffset = setTimeout(function() {
        localStorage.setItem("panelOffset", yOffset);
    }, 250);

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

        if (element.dataset.fieldset) {
            toggleFieldset(element.dataset.fieldset, element.checked);
        }
    }

    if (configInputTypes[key] == "select" && element.classList.contains("font")) {
        loadFont(value);
    }

    changeConfigKV(key, value);
}

function changeConfigKV(key, value) {
    config[key] = value;
    updateCalendar();
    localStorage.setItem(key, value);
}

function toggleFieldset(fieldsetName, enabled) {
    let fieldset = document.getElementById("fs-" + fieldsetName);
    // query every input, button and select inside .field block
    let inputs = fieldset.querySelectorAll(".field input, .field button, .field select");

    if (enabled) {
        fieldset.classList.remove("disabled");
        for (let i = 0; i < inputs.length; i++) {
            inputs[i].disabled = false;
        }
    } else {
        fieldset.classList.add("disabled");
        for (let i = 0; i < inputs.length; i++) {
            inputs[i].disabled = true;
        }
    }
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

    let templatePage = document.getElementById('template_page').innerHTML;
    let renderedPage = Mustache.render(
        templatePage,
        {
            month: getMonth(cfg),
            monthY: cfg.monthY,
            halfWidth: width/2,
            days: d,
            weekdays: weekdays(cfg),
            weeknumbers: weeknumbers(cfg, w),
            daysFontFamilyLoading: (cfg.daysFontFamily in loadingFonts),
            monthFontFamilyLoading: (cfg.monthFontFamily in loadingFonts),
            weekdaysFontFamilyLoading: (cfg.weekdaysFontFamily in loadingFonts),
            weeknumbersFontFamilyLoading: (cfg.weeknumbersFontFamily in loadingFonts),
        }
    );

    let k = 3.7795;

    cfg.daysFontSize /= k;
    cfg.monthFontSize /= k;
    cfg.weekdaysFontSize /= k;
    cfg.weeknumbersFontSize /= k;

    let templateStyles = document.getElementById('template_styles').innerHTML;
    let renderedStyles = Mustache.render(templateStyles, cfg);

    document.getElementById('defs').innerHTML = renderedPage;
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
            x: column * cfg.daysXStep + cfg.daysX,
            y: row * cfg.daysYStep + cfg.daysY,
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
            x: column * cfg.daysXStep + cfg.daysX,
            y: row * cfg.daysYStep + cfg.daysY,
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
            x: i * cfg.daysXStep + cfg.weekdaysX,
            y: cfg.weekdaysY,
        });
    }

    return days;
}

function weeknumbers(cfg, w) {
    let result = [];

    for (let i = 0; i < w.length; i++) {
        result.push({
            weeknumber: w[i],
            x: cfg.weeknumbersX,
            y: i * cfg.daysYStep + cfg.weeknumbersY,
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
