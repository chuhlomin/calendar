"use strict";

let today = new Date();

let month = today.getMonth();
let year = today.getFullYear();

// show next month by default
month++;
if (month == 12) {
  month = 0;
  year++;
}

let config = {
  language: "en",
  pageSize: "A4",
  width: 210,
  height: 297,
  firstDay: "0",
  year: year,
  month: month, // changed by loadConfig if not found in localStorage

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
  monthFormat: "January 2006",
  monthFontFamily: "iosevka-aile-regular",
  monthFontSize: "50",
  monthY: "260",
  monthColor: "#222222",

  // weekdays
  showWeekdays: "true",
  weekdaysFontFamily: "iosevka-aile-regular",
  weekdaysFontSize: "18",
  weekdaysX: "40",
  weekdaysY: "32",
  weekdaysColor: "#999999",

  // week numbers
  showWeekNumbers: "true",
  weeknumbersFontFamily: "iosevka-regular",
  weeknumbersFontSize: "16",
  weeknumbersX: "20",
  weeknumbersY: "43",
  weeknumbersColor: "#999999",
};

let configInputTypes = {
  language: "select",
  pageSize: "select",
  firstDay: "radio",
  month: "select",
  year: "number",

  // days
  daysFontSize: "number",
  daysFontFamily: "select",
  daysX: "number",
  daysY: "number",
  daysXStep: "number",
  daysYStep: "number",
  textColor: "color",
  weekendColor: "color",
  showInactiveDays: "checkbox",
  inactiveColor: "color",

  // month
  showMonth: "checkbox",
  monthFormat: "select",
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

const styleOnlyConfig = [
  "daysFontSize",
  "textColor",
  "weekendColor",
  "inactiveColor",
  "monthColor",
  "weekdaysColor",
  "weeknumbersColor",
];

const gap = 10; // mm

let configIntegerFields = ["firstDay", "year", "month"];

let width = 0;
let height = 0;

let body = document.getElementsByTagName("body")[0];
let panel = document.getElementsByClassName("panel")[0];
let language = document.getElementById("language");
let pageSize = document.getElementById("pageSize");
let yearInput = document.getElementById("year");
let canvas = document.getElementById("canvas");
let calendar = document.getElementById("calendar");
let style = document.getElementById("style");
let svg = document.getElementById("svg");
let rect = document.getElementById("rect");
let templateStyles = document.getElementById("template_styles").innerHTML;
let templateDays = document.getElementById("template_days").innerHTML;
let templateMonth = document.getElementById("template_month").innerHTML;
let templateWeekdays = document.getElementById("template_weekdays").innerHTML;
let templateWeeknumbers = document.getElementById(
  "template_weeknumbers",
).innerHTML;

// calendarThresholds calculated based on page size (see `updatePage`)
// and affects optimal number of rows and columns
let calendarThresholds = [];

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

let weekdays = ["Sun", "Mon", "Tue", "Wed", "Thu", "Fri", "Sat"];

let generatingLabel = "Generating";

let availableFonts = {
  "firacode-regular": "Fira Code",
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
  "jetbrains-mono-regular": "JetBrains Mono",
  "naga-regular": "Naga",
  "martianmono-regular": "MartianMono",
  "martianmono-condenced-regular": "MartianMono Condenced",
  "onest-regular": "Onest",
};
let loadedFonts = {};
let loadingFonts = {};
let erroredFonts = {};

const errorNode = document.createElement("span");
errorNode.classList.add("error");
errorNode.textContent = " ⚠";

const fontBaseURL = "https://fonts.calendar.chuhlomin.com/";

function loadFont(fontName) {
  if (loadedFonts[fontName] || loadingFonts[fontName]) {
    return;
  }

  loadingFonts[fontName] = true;

  const font = new FontFace(
    fontName,
    "url('" +
      fontBaseURL +
      fontName +
      ".ttf') format('truetype'), url('https://fonts.calendar.chuhlomin.com/" +
      fontName +
      ".woff') format('woff2')",
    {
      // range for numbers, latin characters and cyrillic
      unicodeRange: "U+0020-007F, U+00A0-00FF, U+0400-04FF",
    },
  );
  font
    .load()
    .then(function (loadedFace) {
      document.fonts.add(loadedFace);
      loadedFonts[fontName] = true;
      delete loadingFonts[fontName];
      updateCalendar();

      // delete all .error child nodes
      document
        .querySelectorAll("option[value='" + fontName + "'].error")
        .forEach(function (option) {
          option.classList.remove("error");
          option.removeChild(errorNode);
        });
    })
    .catch(function (error) {
      console.log(error);
      delete loadingFonts[fontName];
      erroredFonts[fontName] = true;

      // mark font as errored in select option
      document
        .querySelectorAll("option[value='" + fontName + "']:not(.error)")
        .forEach(function (option) {
          option.classList.add("error");
          option.appendChild(errorNode.cloneNode(true));
        });
    });
}

function loadConfig(key) {
  let value = localStorage.getItem(key);
  if (value) {
    config[key] = value;

    switch (configInputTypes[key]) {
      case "select":
        let select = document.getElementsByName(key)[0];
        select.value = value;
        if (select.classList.contains("font")) {
          loadFont(value);
        }

        if (key == "language") {
          updateLanguage();
        }

        if (key == "pageSize") {
          updatePage();
        }
        break;

      case "radio":
        let option = document.querySelector(
          "input[name='" + key + "'][value='" + value + "']",
        );
        option.checked = true;
        break;

      case "number":
        let input = document.getElementsByName(key)[0];
        input.value = value;

        if (key == "year" || key == "month") {
          updateFormats();
        }
        break;

      case "checkbox":
        let checkbox = document.getElementsByName(key)[0];
        checkbox.checked = value == "true";
        if (checkbox.dataset.fieldset) {
          toggleFieldset(checkbox.dataset.fieldset, checkbox.checked);
        }
        if (checkbox.dataset.field) {
          toggleField(checkbox.dataset.field, checkbox.checked);
        }
        break;
    }
  } else {
    // value not found in localStorage
    if (configInputTypes[key] == "select") {
      let select = document.getElementsByName(key)[0];
      if (select.classList.contains("font")) {
        loadFont(config[key]);
      }
    }
    if (key == "year") {
      yearInput.value = config.year;
    }
    if (key == "month") {
      // make sure the month by default is the next month
      let select = document.getElementsByName(key)[0];
      select.value = config.month;
    }
    if (key == "pageSize") {
      updatePage();
    }
  }
}

function validateConfig(config) {
  let errors = [];

  let cfg = {};

  for (let key in config) {
    if (
      configInputTypes[key] == "number" ||
      configIntegerFields.includes(key)
    ) {
      let value = parseInt(config[key]);
      if (isNaN(value)) {
        errors.push("Invalid value for " + key + ": " + config[key]);
        continue;
      }

      if (key == "month" && (value < -1 || value > 11)) {
        // -1 is for all months
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

window.onload = function () {
  // restore panel scroll position
  let panelOffset = localStorage.getItem("panelOffset");
  if (panelOffset) {
    panel.scrollTop = panelOffset;
  }

  // update select.font options
  let fontSelects = document.querySelectorAll("select.font");
  for (const select of fontSelects) {
    // remove all options
    while (select.firstChild) {
      select.removeChild(select.firstChild);
    }
    // add options
    for (let font in availableFonts) {
      let option = document.createElement("option");
      option.value = font;
      option.text = availableFonts[font];

      if (select.name == "daysFontFamily" && font == config.daysFontFamily) {
        option.selected = true;
      }
      if (select.name == "monthFontFamily" && font == config.monthFontFamily) {
        option.selected = true;
      }
      if (
        select.name == "weekdaysFontFamily" &&
        font == config.weekdaysFontFamily
      ) {
        option.selected = true;
      }
      if (
        select.name == "weeknumbersFontFamily" &&
        font == config.weeknumbersFontFamily
      ) {
        option.selected = true;
      }

      select.add(option);
    }
  }

  // load config from local storage
  for (let key in config) {
    loadConfig(key);

    if (configInputTypes[key] == "color") {
      let button = document.getElementsByName(key)[0];
      new ColorPicker(button, config[key]);
      button.addEventListener("colorChange", (event) => {
        changeConfigKV(key, event.detail.color.hexa);
      });
    }
  }

  updateCalendar();
};

let currentRows = 0;

function maybeResizeCalendar() {
  // only execute it for year view
  if (config.month != "-1") {
    return;
  }

  let rows = getOptimalNumberOfRows(width, height);
  if (rows != currentRows) {
    currentRows = rows;
    resizeCalendar(rows, width, height);
  }
}

window.onresize = maybeResizeCalendar;

// handle panel scroll shadow effect
let yOffset = 0;
let timerForYOffset = null;

panel.onscroll = function () {
  if (!body.classList.contains("preview")) {
    yOffset = panel.scrollTop;
  }

  // start timer if not already running
  if (timerForYOffset) {
    clearTimeout(timerForYOffset);
  }
  timerForYOffset = setTimeout(function () {
    localStorage.setItem("panelOffset", yOffset);
  }, 250);

  if (window.innerHeight + yOffset >= panel.scrollHeight) {
    body.classList.add("scrolled");
  } else {
    body.classList.remove("scrolled");
  }
};

const possibleRows = [2, 3, 4, 6];

function updatePage() {
  let pageData = pageSize.querySelector(
    "option[value='" + config.pageSize + "']",
  ).dataset;

  width = parseInt(pageData.width);
  height = parseInt(pageData.height);

  config.width = width;
  config.height = height;

  calculateThresholds();
  if (config.month != "-1") {
    // year
    svg.setAttribute("viewBox", "0 0 " + width + " " + height);
  }
  rect.setAttribute("width", width);
  rect.setAttribute("height", height);
}

function calculateThresholds() {
  let prevCalendarAR;
  for (const rows of possibleRows) {
    let columns = 12 / rows;
    let calendarWidth = width * columns + gap * (columns - 1);
    let calendarHeight = height * rows + gap * (rows - 1);
    let calendarAR = calendarWidth / calendarHeight;

    if (prevCalendarAR !== undefined) {
      calendarThresholds[rows] = (prevCalendarAR + calendarAR) / 2;
    }
    prevCalendarAR = calendarAR;
  }
}

function updateLanguage() {
  let language = config.language;

  fetch("/langs/" + language + ".json")
    .then((response) => response.json())
    .then((data) => {
      // find all elements with data-i18n attribute
      let elements = document.querySelectorAll("[data-i18n]");

      for (const element of elements) {
        let key = element.dataset.i18n;
        if (data[key]) {
          element.innerHTML = data[key];
        }
      }

      months = [
        data.month_january,
        data.month_february,
        data.month_march,
        data.month_april,
        data.month_may,
        data.month_june,
        data.month_july,
        data.month_august,
        data.month_september,
        data.month_october,
        data.month_november,
        data.month_december,
      ];

      weekdays = [
        data.weekday_short_sunday,
        data.weekday_short_monday,
        data.weekday_short_tuesday,
        data.weekday_short_wednesday,
        data.weekday_short_thursday,
        data.weekday_short_friday,
        data.weekday_short_saturday,
      ];

      generatingLabel = data.generating;

      updateFormats();
      updateCalendar();
    })
    .catch((error) => {
      alert("Error: " + error);
    });
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
  if (key == "pageSize") {
    updatePage();
  }
  if (key == "year" || key == "month") {
    updateFormats();
  }
  if (key == "language") {
    updateLanguage();
  }
  if (key == "showInactiveDays") {
    toggleField("inactive-color", value == "true");
  }

  updateCalendar(styleOnlyConfig.includes(key));
  localStorage.setItem(key, value);
}

function toggleFieldset(fieldsetName, enabled) {
  let fieldset = document.getElementById("fs-" + fieldsetName);
  // query every input, button and select inside .field block
  let inputs = fieldset.querySelectorAll(
    ".field input, .field button, .field select",
  );

  if (enabled) {
    fieldset.classList.remove("disabled");
    for (const input of inputs) {
      input.disabled = false;
    }
  } else {
    fieldset.classList.add("disabled");
    for (const input of inputs) {
      input.disabled = true;
    }
  }
}

function toggleField(fieldName, enabled) {
  let field = document.getElementById("field-" + fieldName);
  let input = field.querySelector("input, button, select");

  if (enabled) {
    field.classList.remove("disabled");
    input.disabled = false;
  } else {
    field.classList.add("disabled");
    input.disabled = true;
  }
}

function updateFormats() {
  let options = document.getElementsByClassName("month-format");
  for (const option of options) {
    let format = option.value;
    let result = format;

    option.innerHTML = result.replace("January", months[0]);
  }
}

function updateCalendar(styleOnly = false) {
  let [cfg, errors] = validateConfig(config);
  if (errors.length > 0) {
    console.error(errors);
    return;
  }

  updateSVGStyles(cfg);

  if (styleOnly) {
    return;
  }

  if (cfg.month == -1) {
    drawYear(cfg, width, height);
    return;
  }

  drawMonth(cfg, width, height);
}

function days(cfg, month) {
  let lastDay = cfg.firstDay - 1;
  if (lastDay < 0) {
    lastDay = 6;
  }

  let days = [];
  let weekNumbers = [];
  let date = new Date(cfg.year, month, 1);
  let end = new Date(cfg.year, month + 1, 0);
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

  // day shift to get week number
  let shift = new Date(date);
  shift.setDate(shift.getDate() + 3);

  let weeknumber = getWeekNumber(shift);

  while (date <= end) {
    let inactive = date.getMonth() != month;

    if (date.getDay() == cfg.firstDay && cfg.showWeekNumbers) {
      weekNumbers.push({
        weeknumber: weeknumber,
        x: cfg.weeknumbersX,
        y: row * cfg.daysYStep + cfg.weeknumbersY,
      });
      weeknumber++;

      if (weeknumber > 52) {
        weeknumber = 1;
      }
    }

    if (row == 0 && inactive && !cfg.showInactiveDays) {
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

  return [days, weekNumbers];
}

function weekdaysLayout(cfg) {
  if (!cfg.showWeekdays) {
    return [];
  }

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

function drawMonth(cfg, width, height) {
  svg.setAttribute("viewBox", "0 0 " + width + " " + height);
  calendar.innerHTML = "";
  addMonthToPage(calendar, cfg, width, cfg.month);
}

function getOptimalNumberOfRows(width, height) {
  let canvasAR = (canvas.offsetWidth / canvas.offsetHeight).toFixed(2);

  // check calendarThresholds in possibleRows
  // iterate over possibleRows from last to first
  for (let i = possibleRows.length - 1; i >= 0; i--) {
    if (canvasAR < calendarThresholds[possibleRows[i]]) {
      return possibleRows[i];
    }
  }

  return 2;
}

function drawYear(cfg, width, height) {
  let rows = getOptimalNumberOfRows(width, height);
  currentRows = rows;
  let columns = 12 / rows;

  svg.setAttribute(
    "viewBox",
    "0 0 " +
      (width * columns + gap * (columns - 1)) +
      " " +
      (height * rows + gap * (rows - 1)),
  );
  calendar.innerHTML = "";

  let month = 0;
  for (let y = 0; y < rows; y++) {
    for (let x = 0; x < columns; x++) {
      let page = document.createElementNS("http://www.w3.org/2000/svg", "g");
      page.setAttribute("id", "month_" + month);
      page.setAttribute(
        "transform",
        "translate(" + x * (width + gap) + "," + y * (height + gap) + ")",
      );
      page.setAttribute("width", width);
      page.setAttribute("height", height);

      addMonthToPage(page, cfg, width, month);

      calendar.appendChild(page);

      month++;
    }
  }
}

function resizeCalendar(rows, width, height) {
  let columns = 12 / rows;

  svg.setAttribute(
    "viewBox",
    "0 0 " +
      (width * columns + gap * (columns - 1)) +
      " " +
      (height * rows + gap * (rows - 1)),
  );

  let month = 0;
  for (let y = 0; y < rows; y++) {
    for (let x = 0; x < columns; x++) {
      let page = document.getElementById("month_" + month);
      page.setAttribute(
        "transform",
        "translate(" + x * (width + gap) + "," + y * (height + gap) + ")",
      );
      month++;
    }
  }
}

function addMonthToPage(page, cfg, width, month) {
  let rect = document.createElementNS("http://www.w3.org/2000/svg", "use");
  rect.setAttribute("href", "#rect");
  rect.setAttribute("x", "0");
  rect.setAttribute("y", "0");
  rect.setAttribute("fill", "url('#rect')");
  page.appendChild(rect);

  // days
  const [d, w] = days(cfg, month);
  let renderedDays = Mustache.render(templateDays, {
    days: d,
    daysFontFamilyLoading:
      cfg.daysFontFamily in loadingFonts || cfg.daysFontFamily in erroredFonts,
  });
  let daysGroup = document.createElementNS("http://www.w3.org/2000/svg", "g");
  daysGroup.classList.add("days");
  daysGroup.innerHTML = renderedDays;
  page.appendChild(daysGroup);

  // month
  let monthGroup = document.createElementNS("http://www.w3.org/2000/svg", "g");
  monthGroup.classList.add("month");
  monthGroup.innerHTML = Mustache.render(templateMonth, {
    month: getMonth(cfg, month),
    x: width / 2,
    y: cfg.monthY,
    year: cfg.year,
    halfWidth: width / 2,
    monthFontFamilyLoading:
      cfg.monthFontFamily in loadingFonts ||
      cfg.monthFontFamily in erroredFonts,
  });
  page.appendChild(monthGroup);

  // weekdays
  let weekdaysGroup = document.createElementNS(
    "http://www.w3.org/2000/svg",
    "g",
  );
  weekdaysGroup.classList.add("weekdays");
  weekdaysGroup.innerHTML = Mustache.render(templateWeekdays, {
    weekdays: weekdaysLayout(cfg),
    weekdaysFontFamilyLoading:
      cfg.weekdaysFontFamily in loadingFonts ||
      cfg.weekdaysFontFamily in erroredFonts,
  });
  page.appendChild(weekdaysGroup);

  // week numbers
  let weeknumbersGroup = document.createElementNS(
    "http://www.w3.org/2000/svg",
    "g",
  );
  weeknumbersGroup.classList.add("weeknumbers");
  weeknumbersGroup.innerHTML = Mustache.render(templateWeeknumbers, {
    weeknumbers: w,
    weeknumbersFontFamilyLoading:
      cfg.weeknumbersFontFamily in loadingFonts ||
      cfg.weeknumbersFontFamily in erroredFonts,
  });
  page.appendChild(weeknumbersGroup);
}

function updateSVGStyles(cfg) {
  let k = 3.7795;

  cfg.daysFontSize /= k;
  cfg.monthFontSize /= k;
  cfg.weekdaysFontSize /= k;
  cfg.weeknumbersFontSize /= k;

  style.innerHTML = Mustache.render(templateStyles, cfg);
}

function getMonth(cfg, month) {
  if (!cfg.showMonth) {
    return "";
  }

  return cfg.monthFormat
    .replace("2006", "_year_")
    .replace("01", "_monthPad_")
    .replace("1", "_month_")
    .replace("_year_", cfg.year)
    .replace("_monthPad_", ("0" + (month + 1).toString()).slice(-2))
    .replace("_month_", month + 1)
    .replace("January", months[month]);
}

function getWeekNumber(d) {
  d = new Date(Date.UTC(d.getFullYear(), d.getMonth(), d.getDate()));
  d.setUTCDate(d.getUTCDate() + 4 - (d.getUTCDay() || 7));
  let yearStart = new Date(Date.UTC(d.getUTCFullYear(), 0, 1));
  return Math.ceil(((d - yearStart) / 86400000 + 1) / 7);
}

function togglePreview(element) {
  if (body.classList.contains("preview")) {
    body.classList.remove("preview");
    panel.scrollTop = yOffset;
  } else {
    body.classList.add("preview");
    maybeResizeCalendar();
  }
}

function buildFilename(cfg) {
  let filename = "calendar_" + cfg.year;
  if (cfg.month != -1) {
    filename += "_" + (cfg.month + 1);
  }
  filename += ".pdf";
  return filename;
}

Date.prototype.addDays = function (days) {
  var date = new Date(this.valueOf());
  date.setDate(date.getDate() + days);
  return date;
};

async function submitForm(element) {
  let [cfg, errors] = validateConfig(config);
  if (errors.length > 0) {
    alert(errors.join("\r"));
    return;
  }

  let defaultLabel = element.innerHTML;
  element.setAttribute("disabled", "disabled");
  element.innerHTML = generatingLabel + '<span class="may-hide"> PDF</span>...';

  setTimeout(() => {
    generatePDF(cfg).finally(() => {
      element.removeAttribute("disabled");
      element.innerHTML = defaultLabel;
    });
  }, 100);
}

async function generatePDF(cfg) {
  const doc = new jspdf.jsPDF({
    orientation: "portrait",
    unit: "mm",
    format: [cfg.width, cfg.height],
  });

  addFonts(doc, cfg);

  if (cfg.month == -1) {
    for (let month = 0; month < 12; month++) {
      if (month > 0) {
        doc.addPage();
      }
      addMonthPage(doc, cfg, month);
    }
  } else {
    addMonthPage(doc, cfg, cfg.month);
  }

  doc.save(buildFilename(cfg));
}

function addFonts(doc, cfg) {
  const fonts = {};
  fonts[cfg.daysFontFamily] = null;
  fonts[cfg.monthFontFamily] = null;
  fonts[cfg.weekdaysFontFamily] = null;
  fonts[cfg.weeknumbersFontFamily] = null;

  Object.keys(fonts).map((font) => {
    doc.addFont(fontBaseURL + font + ".ttf", font, "normal");
  });
}

function addMonthPage(doc, cfg, month) {
  const [d, w] = days(cfg, month);

  drawDaysPDF(doc, cfg, d);

  if (cfg.showMonth) {
    drawMonthPDF(doc, cfg, cfg.year, month);
  }

  if (cfg.showWeekdays) {
    drawWeekdaysPDF(doc, cfg);
  }

  if (cfg.showWeekNumbers) {
    drawWeekNumbersPDF(doc, cfg, w);
  }
}

function drawDaysPDF(doc, cfg, d) {
  doc.setFont(cfg.daysFontFamily);
  doc.setFontSize(cfg.daysFontSize);

  const rgba = hexToRgb(cfg.textColor);
  const rgbaInactive = hexToRgb(cfg.inactiveColor);
  const rgbaWeekend = hexToRgb(cfg.weekendColor);

  let color = rgba;

  for (let i = 0; i < d.length; i++) {
    let day = d[i];

    color = rgba;
    if (day.inactive) {
      color = rgbaInactive;
    } else if (day.weekend) {
      color = rgbaWeekend;
    }

    if (day.inactive && !cfg.showInactiveDays) {
      continue;
    }

    doc.setTextColor(color.r, color.g, color.b, { a: color.a });
    doc.text(day.day.toString(), day.x, day.y, { align: "right" });
  }
}

function drawMonthPDF(doc, cfg, year, month) {
  doc.setFont(cfg.monthFontFamily);
  doc.setFontSize(cfg.monthFontSize);

  const rgba = hexToRgb(cfg.monthColor);
  doc.setTextColor(rgba.r, rgba.g, rgba.b, { a: rgba.a });

  doc.text(getMonth(cfg, month), cfg.width / 2, cfg.monthY, {
    align: "center",
  });
}

function drawWeekdaysPDF(doc, cfg) {
  doc.setFont(cfg.weekdaysFontFamily);
  doc.setFontSize(cfg.weekdaysFontSize);

  const rgba = hexToRgb(cfg.weekdaysColor);
  doc.setTextColor(rgba.r, rgba.g, rgba.b, { a: rgba.a });

  weekdays = weekdaysLayout(cfg);

  for (let day = 0; day < weekdays.length; day++) {
    let w = weekdays[day];
    doc.text(w.day, w.x, w.y, { align: "right" });
  }
}

function drawWeekNumbersPDF(doc, cfg, weekNumbers) {
  doc.setFont(cfg.weeknumbersFontFamily);
  doc.setFontSize(cfg.weeknumbersFontSize);

  const rgba = hexToRgb(cfg.weeknumbersColor);
  doc.setTextColor(rgba.r, rgba.g, rgba.b, { a: rgba.a });

  for (let i = 0; i < weekNumbers.length; i++) {
    let wn = weekNumbers[i];
    doc.text(wn.weeknumber.toString(), wn.x, wn.y, { align: "right" });
  }
}

function hexToRgb(hex) {
  return {
    r: parseInt(hex.substring(1, 3), 16),
    g: parseInt(hex.substring(3, 5), 16),
    b: parseInt(hex.substring(5, 7), 16),
    a: 1,
  };
}
