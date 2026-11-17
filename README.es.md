# TaskGlass

[![CI](https://github.com/AnderCMD/TaskGlass/actions/workflows/ci.yml/badge.svg)](https://github.com/AnderCMD/TaskGlass/actions/workflows/ci.yml)
[![Licencia: MIT](https://img.shields.io/badge/licencia-MIT-blue.svg)](LICENSE)

App de bandeja, ligera y de código abierto, que vuelve la barra de tareas de Windows **totalmente transparente, con desenfoque o acrílica**, de modo que solo se vean los iconos. Escrita en Go (sin cgo) más una pequeña DLL nativa auxiliar, sin runtime que instalar.

[Read in English](README.md)

## Características

- Estilos **Clear / Blur / Acrylic / Opaque / Normal**, incluida la transparencia total en la barra XAML de Windows 11
- **Reglas dinámicas**: opaca cuando hay una ventana maximizada, estilo propio con el menú Inicio abierto
- **Multimonitor**: cada barra se gestiona por separado, incluso las que se conectan después
- **Resiste reinicios de Explorer** y reaplica el efecto
- **Config JSON con recarga en caliente**; un error tipográfico no rompe la sesión (se conserva la configuración previa y se registra el error)
- Menú de bandeja, inicio con Windows opcional, instancia única
- Restaura la barra original al salir, y `--restore` para recuperarla tras un cierre inesperado

## Instalación

Descarga el zip desde [Releases](https://github.com/AnderCMD/TaskGlass/releases), extráelo y ejecuta `taskglass.exe`. Deja `taskglass_tap.dll` en la misma carpeta (ver [Cómo funciona](#cómo-funciona)). Vive en el área de notificaciones.

Compilar desde el código (Windows, Go 1.26+ y un compilador MinGW-w64 para la DLL, p. ej. `winget install MartinStorsjo.LLVM-MinGW.UCRT`):

```bash
go build -ldflags "-H=windowsgui -s -w" -o dist/taskglass.exe ./cmd/taskglass
powershell -File scripts/build-tap.ps1   # genera dist/taskglass_tap.dll
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

## Cómo funciona

Windows tiene una API de acento no documentada, `SetWindowCompositionAttribute`, con la que TaskGlass difumina o tiñe la barra. En Windows 11 la barra pinta además su propio relleno XAML gris semitransparente y un borde encima, así que el acento por sí solo nunca llega a ser totalmente transparente.

Para quitar ese relleno, `taskglass_tap.dll` se carga dentro de `explorer.exe` mediante la API XAML Diagnostics, localiza los rectángulos `BackgroundFill` / `BackgroundStroke` de la barra y pone su opacidad a 0 mientras TaskGlass está en marcha. Si TaskGlass se cierra o falla, la DLL lo detecta y el aspecto original vuelve solo. Detalles en [docs/architecture.md](docs/architecture.md).

## Limitaciones y seguridad

- Depende de funciones internas no documentadas de Windows. Una actualización puede romperlo; en ese caso la app sigue con efectos de acento (tintes y desenfoque) y registra el motivo.
- Cargar código en Explorer es lo que las herramientas de seguridad llaman "inyección", y algunos antivirus pueden avisar. La DLL es pequeña y su código está en [native/tap](native/tap); puedes compilarla tú. Con `"xaml_fix": false` en la configuración nunca se toca Explorer.
- Verificado en Windows 11 compilación 26300 con un monitor y ocultado automático. Windows 10 (sin barra XAML) solo necesita la API de acento y no está probado; multimonitor tampoco. Abre un issue indicando tu versión de Windows si algo se ve mal.
- Si la barra se ve mal tras un cierre inesperado: `taskglass.exe --restore`, o reinicia Explorer.

## Contribuir

Ver [CONTRIBUTING.md](CONTRIBUTING.md). Seguridad: [SECURITY.md](SECURITY.md).

## Licencia

[MIT](LICENSE)
