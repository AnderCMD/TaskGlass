// Package i18n holds every user-visible string of TaskGlass in English and
// Spanish. The tray menu reads it directly and the settings window receives the
// whole table as JSON, so there is a single source of truth for both.
package i18n

import (
	"sort"
	"strings"
)

// Lang is a supported UI language.
type Lang string

const (
	Auto    Lang = "auto" // follow the Windows display language
	English Lang = "en"
	Spanish Lang = "es"
)

// Langs lists the concrete languages.
var Langs = []Lang{English, Spanish}

// Valid reports whether l may be stored in the configuration.
func (l Lang) Valid() bool { return l == Auto || l == English || l == Spanish }

// Resolve turns Auto into a concrete language using the Windows locale name
// (for example "es-MX" or "en-US").
func Resolve(l Lang, systemLocale string) Lang {
	switch l {
	case English, Spanish:
		return l
	}
	if strings.HasPrefix(strings.ToLower(systemLocale), "es") {
		return Spanish
	}
	return English
}

// T returns the string for key in lang, falling back to English and finally to
// the key itself so a missing translation is visible but never fatal.
func T(lang Lang, key string) string {
	if m, ok := table[lang]; ok {
		if s, ok := m[key]; ok {
			return s
		}
	}
	if s, ok := table[English][key]; ok {
		return s
	}
	return key
}

// Table returns a copy of all strings for lang (English fills any gap).
func Table(lang Lang) map[string]string {
	out := make(map[string]string, len(table[English]))
	for k, v := range table[English] {
		out[k] = v
	}
	for k, v := range table[lang] {
		out[k] = v
	}
	return out
}

// Keys returns the sorted keys of a language, used by tests.
func Keys(lang Lang) []string {
	var ks []string
	for k := range table[lang] {
		ks = append(ks, k)
	}
	sort.Strings(ks)
	return ks
}

var table = map[Lang]map[string]string{English: en, Spanish: es}

var en = map[string]string{
	// tray menu
	"tray.enabled":      "Enabled",
	"tray.enabled.tip":  "Turn the effect on or off",
	"tray.settings":     "Settings…",
	"tray.settings.tip": "Open the settings window with a live preview",
	"tray.style":        "Style",
	"tray.theme":        "Theme",
	"tray.language":     "Language",
	"tray.maximized":    "Opaque when a window is maximized",
	"tray.hotkey":       "Global hotkey (%s)",
	"tray.hotkey.tip":   "Toggle TaskGlass from anywhere; change the keys in Settings",
	"tray.flyout":       "Translucent tray flyout",
	"tray.flyout.tip":   "Make the hidden-icons flyout (the ^ button) translucent",
	"tray.xaml":         "Fully transparent (Windows 11)",
	"tray.xaml.tip":     "Remove the taskbar's own XAML fill",
	"tray.autostart":    "Start with Windows",
	"tray.openconfig":   "Open config file",
	"tray.openlog":      "Open log",
	"tray.reload":       "Reload config",
	"tray.quit":         "Quit TaskGlass",
	"accent.clear":      "Clear",
	"accent.blur":       "Blur",
	"accent.acrylic":    "Acrylic",
	"accent.opaque":     "Opaque",
	"accent.normal":     "Normal (Windows default)",
	"theme.auto":        "Follow Windows",
	"theme.light":       "Light",
	"theme.dark":        "Dark",
	"language.auto":     "Automatic (Windows)",
	"language.en":       "English",
	"language.es":       "Español",

	// settings window
	"ui.title":            "TaskGlass Settings",
	"ui.save":             "Save",
	"ui.saved":            "Saved",
	"ui.saveFailed":       "Could not save: %s",
	"ui.revert":           "Discard changes",
	"ui.reset":            "Restore defaults",
	"ui.resetConfirm":     "Restore every setting to its default?",
	"ui.unsaved":          "Unsaved changes",
	"ui.preview":          "Live preview",
	"ui.previewNote":      "Approximation of your taskbar. Changes apply to the real taskbar when you press Save.",
	"ui.previewWallpaper": "Preview background",
	"ui.previewTray":      "Hidden icons flyout",
	"ui.previewTaskbar":   "Taskbar",
	"ui.previewMax":       "With a maximized window",
	"ui.previewApp":       "App rule: %s",
	"ui.tab.general":      "General",
	"ui.tab.appearance":   "Appearance",
	"ui.tab.tray":         "Hidden icons",
	"ui.tab.apps":         "Per app",
	"ui.tab.maximized":    "Maximized",
	"ui.tab.advanced":     "Advanced",
	"ui.tab.about":        "About",

	"general.enabled":          "Enable TaskGlass",
	"general.enabled.help":     "When off, the stock Windows taskbar is shown.",
	"general.autostart":        "Start with Windows",
	"general.autostart.help":   "Launch TaskGlass automatically when you sign in.",
	"general.language":         "Language",
	"general.theme":            "Windows theme",
	"general.theme.help":       "Used for automatic colours and the tray icon.",
	"general.taskbars":         "Taskbars to style",
	"general.taskbars.all":     "All monitors",
	"general.taskbars.primary": "Primary monitor only",
	"general.trayIcon":         "Tray icon style",
	"general.trayIcon.auto":    "Automatic (matches Windows)",
	"general.trayIcon.color":   "Colour logo",
	"general.trayIcon.light":   "White glyph",
	"general.trayIcon.dark":    "Dark glyph",
	"general.hotkey":           "Global shortcut",
	"general.hotkey.help":      "Toggles TaskGlass from anywhere. Needs Ctrl, Alt, Shift or Win plus A-Z, 0-9 or F1-F24.",
	"general.hotkey.enable":    "Enable the shortcut",

	"appear.style":          "Taskbar style",
	"appear.color":          "Tint colour",
	"appear.colorAuto":      "Follow the Windows theme",
	"appear.opacity":        "Tint strength",
	"appear.opacity.help":   "0 is the most transparent, 100 the most solid.",
	"appear.preset":         "Quick presets",
	"appear.preset.glass":   "Pure glass",
	"appear.preset.frosted": "Frosted",
	"appear.preset.acrylic": "Acrylic tint",
	"appear.preset.solid":   "Solid",
	"appear.maxRule":        "Different look when a window is maximized",
	"appear.maxRule.help":   "Lets the bar turn solid (or any other style) while a maximized window is on that monitor.",
	"appear.xaml":           "Fully transparent taskbar (Windows 11)",
	"appear.xaml.help":      "Removes Windows 11's own grey fill. Without it Clear, Blur and Acrylic cannot get fully transparent. Needs taskglass_tap.dll.",
	"appear.xaml.missing":   "taskglass_tap.dll was not found next to the program, so this option has no effect.",
	"appear.accentNote":     "Windows' \"Show accent color on Start and taskbar\" must be off for the effects to show.",
	"style.clear.desc":      "Completely see-through background.",
	"style.blur.desc":       "Gaussian blur of what is behind the bar, with an optional tint.",
	"style.acrylic.desc":    "Blur plus the frosted noise texture, with an optional tint.",
	"style.opaque.desc":     "A solid colour.",
	"style.normal.desc":     "Leave the taskbar exactly as Windows draws it.",

	"tray.surface":      "Translucent hidden-icons flyout",
	"tray.surface.help": "Makes the panel that opens with the ^ button (hidden tray icons) translucent. Needs taskglass_tap.dll.",
	"tray.opacity":      "Flyout background opacity",
	"tray.opacity.help": "0 removes the background, 100 is the stock Windows look.",

	"apps.title":  "Per-application styles",
	"apps.help":   "While one of these programs has focus, the taskbar uses the style below instead of the default. The first match wins.",
	"apps.add":    "Add application",
	"apps.exe":    "Program (e.g. chrome.exe)",
	"apps.remove": "Remove",
	"apps.empty":  "No rules yet.",
	"apps.limit":  "Up to %d rules.",

	"adv.poll":       "Fallback check interval (ms)",
	"adv.debounce":   "Event debounce (ms)",
	"adv.refresh":    "Re-apply interval (ms)",
	"adv.log":        "Log level",
	"adv.openConfig": "Open config file",
	"adv.openLog":    "Open log",
	"adv.help":       "Leave these alone unless you know you need them.",

	"about.version": "Version",
	"about.license": "Licensed under the Apache License 2.0.",
	"about.site":    "Project page",

	"err.invalid":        "Invalid settings: %s",
	"appear.pickColor":   "Choose…",
	"appear.blurNote":    "Windows 11 draws the classic blur as solid black, so Blur uses the acrylic backdrop; its strength lightens or darkens the glass.",
	"ui.closeUnsaved":    "You have unsaved changes. Save them before closing?",
	"ui.scenario":        "Show:",
	"apps.pick":          "Select a rule on the left, or add an application.",
	"ui.hotkeyInvalid":   "Invalid shortcut",
	"ui.previewDisabled": "TaskGlass is off: stock Windows taskbar",
}

var es = map[string]string{
	// tray menu
	"tray.enabled":      "Activado",
	"tray.enabled.tip":  "Activa o desactiva el efecto",
	"tray.settings":     "Configuración…",
	"tray.settings.tip": "Abre la ventana de configuración con vista previa en vivo",
	"tray.style":        "Estilo",
	"tray.theme":        "Tema",
	"tray.language":     "Idioma",
	"tray.maximized":    "Opaca con una ventana maximizada",
	"tray.hotkey":       "Atajo global (%s)",
	"tray.hotkey.tip":   "Activa o desactiva TaskGlass desde cualquier lugar; cambia las teclas en Configuración",
	"tray.flyout":       "Panel de iconos ocultos translúcido",
	"tray.flyout.tip":   "Vuelve translúcido el panel de iconos ocultos (botón ^)",
	"tray.xaml":         "Totalmente transparente (Windows 11)",
	"tray.xaml.tip":     "Quita el relleno XAML propio de la barra de tareas",
	"tray.autostart":    "Iniciar con Windows",
	"tray.openconfig":   "Abrir archivo de configuración",
	"tray.openlog":      "Abrir registro",
	"tray.reload":       "Recargar configuración",
	"tray.quit":         "Salir de TaskGlass",
	"accent.clear":      "Transparente",
	"accent.blur":       "Desenfoque",
	"accent.acrylic":    "Acrílico",
	"accent.opaque":     "Sólido",
	"accent.normal":     "Normal (predeterminado de Windows)",
	"theme.auto":        "Seguir a Windows",
	"theme.light":       "Claro",
	"theme.dark":        "Oscuro",
	"language.auto":     "Automático (Windows)",
	"language.en":       "English",
	"language.es":       "Español",

	// settings window
	"ui.title":            "Configuración de TaskGlass",
	"ui.save":             "Guardar",
	"ui.saved":            "Guardado",
	"ui.saveFailed":       "No se pudo guardar: %s",
	"ui.revert":           "Descartar cambios",
	"ui.reset":            "Restaurar valores predeterminados",
	"ui.resetConfirm":     "¿Restaurar todos los ajustes a sus valores predeterminados?",
	"ui.unsaved":          "Cambios sin guardar",
	"ui.preview":          "Vista previa en vivo",
	"ui.previewNote":      "Aproximación de tu barra de tareas. Los cambios se aplican a la barra real al pulsar Guardar.",
	"ui.previewWallpaper": "Fondo de la vista previa",
	"ui.previewTray":      "Panel de iconos ocultos",
	"ui.previewTaskbar":   "Barra de tareas",
	"ui.previewMax":       "Con una ventana maximizada",
	"ui.previewApp":       "Regla de aplicación: %s",
	"ui.tab.general":      "General",
	"ui.tab.appearance":   "Apariencia",
	"ui.tab.tray":         "Iconos ocultos",
	"ui.tab.apps":         "Por aplicación",
	"ui.tab.maximized":    "Maximizada",
	"ui.tab.advanced":     "Avanzado",
	"ui.tab.about":        "Acerca de",

	"general.enabled":          "Activar TaskGlass",
	"general.enabled.help":     "Si está desactivado, se muestra la barra de tareas normal de Windows.",
	"general.autostart":        "Iniciar con Windows",
	"general.autostart.help":   "Inicia TaskGlass automáticamente al iniciar sesión.",
	"general.language":         "Idioma",
	"general.theme":            "Tema de Windows",
	"general.theme.help":       "Se usa para los colores automáticos y el icono de la bandeja.",
	"general.taskbars":         "Barras de tareas a estilizar",
	"general.taskbars.all":     "Todos los monitores",
	"general.taskbars.primary": "Solo el monitor principal",
	"general.trayIcon":         "Estilo del icono de la bandeja",
	"general.trayIcon.auto":    "Automático (según Windows)",
	"general.trayIcon.color":   "Logotipo a color",
	"general.trayIcon.light":   "Símbolo blanco",
	"general.trayIcon.dark":    "Símbolo oscuro",
	"general.hotkey":           "Atajo global",
	"general.hotkey.help":      "Activa o desactiva TaskGlass desde cualquier lugar. Requiere Ctrl, Alt, Shift o Win más A-Z, 0-9 o F1-F24.",
	"general.hotkey.enable":    "Activar el atajo",

	"appear.style":          "Estilo de la barra de tareas",
	"appear.color":          "Color del tinte",
	"appear.colorAuto":      "Seguir el tema de Windows",
	"appear.opacity":        "Intensidad del tinte",
	"appear.opacity.help":   "0 es lo más transparente, 100 lo más sólido.",
	"appear.preset":         "Ajustes rápidos",
	"appear.preset.glass":   "Cristal puro",
	"appear.preset.frosted": "Esmerilado",
	"appear.preset.acrylic": "Acrílico con tinte",
	"appear.preset.solid":   "Sólido",
	"appear.maxRule":        "Aspecto distinto con una ventana maximizada",
	"appear.maxRule.help":   "Permite que la barra se vuelva sólida (u otro estilo) mientras haya una ventana maximizada en ese monitor.",
	"appear.xaml":           "Barra totalmente transparente (Windows 11)",
	"appear.xaml.help":      "Quita el relleno gris propio de Windows 11. Sin esto, Transparente, Desenfoque y Acrílico no pueden ser totalmente transparentes. Requiere taskglass_tap.dll.",
	"appear.xaml.missing":   "No se encontró taskglass_tap.dll junto al programa, así que esta opción no tiene efecto.",
	"appear.accentNote":     "La opción de Windows «Mostrar el color de énfasis en Inicio y la barra de tareas» debe estar desactivada para que se vean los efectos.",
	"style.clear.desc":      "Fondo completamente transparente.",
	"style.blur.desc":       "Desenfoque gaussiano de lo que hay detrás de la barra, con tinte opcional.",
	"style.acrylic.desc":    "Desenfoque más la textura de ruido esmerilado, con tinte opcional.",
	"style.opaque.desc":     "Un color sólido.",
	"style.normal.desc":     "Deja la barra de tareas exactamente como la dibuja Windows.",

	"tray.surface":      "Panel de iconos ocultos translúcido",
	"tray.surface.help": "Vuelve translúcido el panel que se abre con el botón ^ (iconos ocultos de la bandeja). Requiere taskglass_tap.dll.",
	"tray.opacity":      "Opacidad del fondo del panel",
	"tray.opacity.help": "0 quita el fondo, 100 es el aspecto normal de Windows.",

	"apps.title":  "Estilos por aplicación",
	"apps.help":   "Mientras uno de estos programas tenga el foco, la barra usa el estilo indicado en lugar del predeterminado. Gana la primera coincidencia.",
	"apps.add":    "Añadir aplicación",
	"apps.exe":    "Programa (p. ej. chrome.exe)",
	"apps.remove": "Quitar",
	"apps.empty":  "Aún no hay reglas.",
	"apps.limit":  "Hasta %d reglas.",

	"adv.poll":       "Intervalo de comprobación de respaldo (ms)",
	"adv.debounce":   "Agrupación de eventos (ms)",
	"adv.refresh":    "Intervalo de reaplicación (ms)",
	"adv.log":        "Nivel de registro",
	"adv.openConfig": "Abrir archivo de configuración",
	"adv.openLog":    "Abrir registro",
	"adv.help":       "No toques esto salvo que sepas que lo necesitas.",

	"about.version": "Versión",
	"about.license": "Licenciado bajo la Licencia Apache 2.0.",
	"about.site":    "Página del proyecto",

	"err.invalid":        "Ajustes no válidos: %s",
	"appear.pickColor":   "Elegir…",
	"appear.blurNote":    "Windows 11 dibuja el desenfoque clásico en negro sólido, así que Desenfoque usa el fondo acrílico; su intensidad aclara u oscurece el cristal.",
	"ui.closeUnsaved":    "Tienes cambios sin guardar. ¿Guardarlos antes de cerrar?",
	"ui.scenario":        "Mostrar:",
	"apps.pick":          "Selecciona una regla a la izquierda o añade una aplicación.",
	"ui.hotkeyInvalid":   "Atajo no válido",
	"ui.previewDisabled": "TaskGlass está apagado: barra de tareas normal de Windows",
}
