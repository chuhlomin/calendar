# Calendar PDF Generator

[![main](https://github.com/chuhlomin/calendar/actions/workflows/main.yml/badge.svg)](https://github.com/chuhlomin/calendar/actions/workflows/main.yml)

Powers [calendar.chuhlomin.com](https://calendar.chuhlomin.com).

<img width="100%" alt="calendar chuhlomin com" src="https://github.com/chuhlomin/calendar/assets/3620471/920492b7-388e-4a50-b690-9b5a37a2bc09">

## Directory structure

```bash
.

├── pages          # static files
│   ├── fonts      # ttf and woff2 fonts used in browser and PDF
│   ├── langs      # language files
│   ├── styles.css # main styles
│   ├── index.js   # main script
│   └── index.html # main page
└── Makefile       # build and run commands
```

## Local development

Install [Wrangler](https://developers.cloudflare.com/workers/wrangler/install-and-update/).

In terminal:

```bash
make dev
```

Open http://localhost:8788

## Adding new language

Add new option to select in `pages/index.html`:

```html
<select name="language" id="language" ...>
  ...
  <option value="en">English</option>
</select>
```

Create new file under `pages/langs`.
