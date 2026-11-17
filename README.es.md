# TaskGlass

[![CI](https://github.com/AnderCMD/TaskGlass/actions/workflows/ci.yml/badge.svg)](https://github.com/AnderCMD/TaskGlass/actions/workflows/ci.yml)
[![Licencia: MIT](https://img.shields.io/badge/licencia-MIT-blue.svg)](LICENSE)

App de bandeja, ligera y de código abierto, que vuelve la barra de tareas de Windows **transparente, con desenfoque o acrílica**, de modo que solo se vean los iconos. Escrita en Go, sin cgo ni runtime: un único `.exe` de ~3 MB.

[Read in English](README.md)

## Características

- Estilos **Clear / Blur / Acrylic / Opaque / Normal**
- **Reglas dinámicas**: opaca cuando hay una ventana maximizada, estilo propio con el menú Inicio abierto
- **Multimonitor**: cada barra se gestiona por separado, incluso las que se conectan después
- **Resiste reinicios de Explorer** y reaplica el efecto
- **Config JSON con recarga en caliente**; un error tipográfico no rompe la sesión (se conserva la configuración previa y se registra el error)
- Menú de bandeja, inicio con Windows opcional, instancia única
- Restaura la barra original al salir, y `--restore` para recuperarla tras un cierre inesperado

## Instalación

Descarga `taskglass.exe` desde [Releases](https://github.com/AnderCMD/TaskGlass/releases) y ejecútalo. Vive en el área de notificaciones.

Compilar desde el código (Go 1.26+, Windows):

```bash
go build -ldflags "-H=windowsgui -s -w" -o dist/taskglass.exe ./cmd/taskglass
```

## Uso

Clic derecho en el icono de bandeja para elegir estilo, activar las reglas dinámicas, el inicio automático o abrir la configuración.

| Flag | Función |
|------|---------|
| `--config <ruta>` | Usa otro archivo de configuración |
| `--reset` | Reescribe la configuración por defecto y sale |
| `--restore` | Devuelve la barra a la normalidad y sale |
| `--verbose` | Registra también en stderr |
| `--version` | Muestra la versión |

Archivos: config `%APPDATA%\TaskGlass\config.json`, log `%LOCALAPPDATA%\TaskGlass\taskglass.log`.

Más detalle en [docs/configuration.md](docs/configuration.md) y [docs/architecture.md](docs/architecture.md).

## Limitaciones

Usa la API no documentada `SetWindowCompositionAttribute`, la misma técnica que TranslucentTB. Microsoft puede cambiar su comportamiento en cualquier actualización. Objetivo: Windows 10 y 11; el aspecto exacto de "clear" en la barra XAML de Windows 11 puede variar entre compilaciones. Si algo se ve mal, abre un issue indicando tu versión de Windows.

## Contribuir

Ver [CONTRIBUTING.md](CONTRIBUTING.md). Seguridad: [SECURITY.md](SECURITY.md).

## Licencia

[MIT](LICENSE)
