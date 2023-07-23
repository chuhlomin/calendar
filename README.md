# Calendar PDF Generator

[![main](https://github.com/chuhlomin/calendar/actions/workflows/main.yml/badge.svg)](https://github.com/chuhlomin/calendar/actions/workflows/main.yml)

Powers [calendar.chuhlomin.com](https://calendar.chuhlomin.com).

<img width="100%" alt="calendar chuhlomin com" src="https://github.com/chuhlomin/calendar/assets/3620471/920492b7-388e-4a50-b690-9b5a37a2bc09">

## Directory structure

```bash
.
├── fonts          # ttf and woff2 fonts used in browser and PDF
├── langs          # language files
├── static         # static files
│   ├── script.js  # main script
│   ├── styles.css # main styles
│   └── index.html # main page
├── vendor         # vendored dependencies
└── main.go        # main application
```

## Local development

Install [Go](https://golang.org/doc/install).

In terminal:

```bash
make run
```

Example output:

```
2022/12/31 10:00:00 Starting...
2022/12/31 10:00:00 Starting server on http://127.0.0.1:8082
```

Open http://localhost:8082

## Adding new language

Add new option to select in `static/index.html`:

```html
<select name="language" id="language" ...>
...
<option value="en">English</option>
</select>
```

Create new file under `langs`.
