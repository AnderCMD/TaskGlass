# Security Policy

## Reporting a vulnerability

Please do **not** open a public issue. Use GitHub's private vulnerability reporting ("Security" tab → "Report a vulnerability") on this repository. You will get an acknowledgement within a few days.

## Scope

TaskGlass runs with user privileges, makes a single kind of network connection, to GitHub, to look for and download new releases (once a day; turn it off with `"updates": "off"`, see [Updates](docs/configuration.md#updates)), and only touches: its config and log files under the user profile, the `HKCU\...\Run\TaskGlass` registry value (when autostart is enabled), and the appearance of the Explorer taskbar windows. Reports about anything beyond that are welcome.

## Supported versions

Only the latest release receives fixes.
