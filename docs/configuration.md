# Configuration

File: `%APPDATA%\TaskGlass\config.json` (created with defaults on first run). Changes are picked up within about a second; no restart needed. If the file is invalid, the running settings are kept and the error is written to the log.

```json
{
  "version": 1,
  "default":    { "accent": "clear",  "color": "#000000", "opacity": 0 },
  "maximized":  { "enabled": true,  "accent": "opaque",  "color": "#000000", "opacity": 100 },
  "start_menu": { "enabled": false, "accent": "acrylic", "color": "#202020", "opacity": 60 },
  "poll_ms": 100,
  "refresh_ms": 1000,
  "log_level": "info"
}
```

| Key | Values | Meaning |
|-----|--------|---------|
| `version` | `1` | Schema version |
| `*.accent` | `clear`, `blur`, `acrylic`, `opaque`, `normal` | Effect. `normal` leaves the stock taskbar |
| `*.color` | `#RRGGBB` | Tint colour |
| `*.opacity` | `0`-`100` | Tint opacity in percent |
| `*.enabled` | bool | Rules `maximized` and `start_menu` only apply when enabled |
| `poll_ms` | 16-5000 | How often window state is sampled |
| `refresh_ms` | `poll_ms`-60000 | Re-apply interval even without changes |
| `log_level` | `debug`, `info`, `warn`, `error` | Log verbosity |

Rule priority: `start_menu` > `maximized` > `default`. Unknown keys are rejected so typos are caught.
