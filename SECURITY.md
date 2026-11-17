# Security Policy

## Reporting a vulnerability

Please do **not** open a public issue. Use GitHub's private vulnerability reporting ("Security" tab → "Report a vulnerability") on this repository. You will get an acknowledgement within a few days.

## Scope

TaskGlass runs with user privileges, makes no network connections and only touches: its config and log files under the user profile, the `HKCU\...\Run\TaskGlass` registry value (when autostart is enabled), and the appearance of the Explorer taskbar windows. Reports about anything beyond that are welcome.

## Supported versions

Only the latest release receives fixes.
