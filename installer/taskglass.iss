; TaskGlass installer (Inno Setup 6). Build with scripts/build-installer.ps1.
; Per-user install: no administrator rights, nothing outside the user profile.
; Bilingual (English / Español): the wizard, the license and the app's first
; configuration follow the language chosen at the start.

#ifndef AppVersion
  #define AppVersion "0.0.0"
#endif
#ifndef DistDir
  #define DistDir "..\dist"
#endif

; Windows version resources only accept N.N.N.N, so strip any suffix like -ci.
#define NumericVersion Copy(AppVersion, 1, Pos("-", AppVersion + "-") - 1)

#define AppName "TaskGlass"
#define AppExe "taskglass.exe"

[Setup]
AppId={{6F1B7C3E-2D4A-4E8B-9C57-A1D3E5F70B29}
AppName={#AppName}
AppVersion={#AppVersion}
AppVerName={#AppName} {#AppVersion}
AppPublisher=AnderCMD
AppPublisherURL=https://github.com/AnderCMD/TaskGlass
AppSupportURL=https://github.com/AnderCMD/TaskGlass/issues
AppUpdatesURL=https://github.com/AnderCMD/TaskGlass/releases
DefaultDirName={autopf}\{#AppName}
DefaultGroupName={#AppName}
DisableProgramGroupPage=yes
PrivilegesRequired=lowest
PrivilegesRequiredOverridesAllowed=dialog
ArchitecturesAllowed=x64compatible
ArchitecturesInstallIn64BitMode=x64compatible
MinVersion=10.0
SetupIconFile=..\assets\taskglass.ico
UninstallDisplayIcon={app}\{#AppExe}
OutputDir={#DistDir}
OutputBaseFilename=TaskGlass-{#AppVersion}-Setup
Compression=lzma2/ultra
SolidCompression=yes
WizardStyle=modern
CloseApplications=yes
RestartApplications=no
ShowLanguageDialog=auto
VersionInfoVersion={#NumericVersion}
VersionInfoProductName={#AppName}

[Languages]
; The license page shows the text in the installer language. The English text
; is the legally binding original; the Spanish one is an unofficial translation.
Name: "english"; MessagesFile: "compiler:Default.isl"; LicenseFile: "..\LICENSE"
Name: "spanish"; MessagesFile: "compiler:Languages\Spanish.isl"; LicenseFile: "..\LICENSE.es.txt"

[Tasks]
Name: "autostart"; Description: "{cm:AutoStart}"; GroupDescription: "{cm:Startup}"
Name: "desktopicon"; Description: "{cm:CreateDesktopIcon}"; GroupDescription: "{cm:AdditionalIcons}"; Flags: unchecked

[CustomMessages]
english.AutoStart=Start TaskGlass when I sign in to Windows
english.Startup=Startup:
english.LaunchNow=Launch TaskGlass now
english.SettingsIcon=TaskGlass Settings
english.DeleteSettings=Also delete your TaskGlass settings and logs?
english.StylePageTitle=Taskbar style
english.StylePageDesc=How should your taskbar look?
english.StylePageSub=You can change this at any time from the TaskGlass settings window.
english.StyleClear=Clear: completely transparent
english.StyleBlur=Blur: blurred background with a light tint
english.StyleAcrylic=Acrylic: blur plus frosted texture
english.StyleOpaque=Opaque: a solid colour
english.StyleNormal=Normal: leave the Windows taskbar untouched
english.LangPageTitle=Language
english.LangPageDesc=Which language should TaskGlass use?
english.LangPageSub=Applies to the tray menu and the settings window.
english.LangAuto=Automatic (follow Windows)
english.LangEn=English
english.LangEs=Español
english.FeatPageTitle=Features
english.FeatPageDesc=Choose the extras you want.
english.FeatPageSub=Each of these can be changed later in the settings window.
english.FeatFlyout=Translucent hidden-icons flyout (the ^ button)
english.FeatHotkey=Global shortcut Ctrl+Alt+G to turn TaskGlass on or off
english.FeatXaml=Fully transparent taskbar on Windows 11 (recommended)
spanish.AutoStart=Iniciar TaskGlass al iniciar sesión en Windows
spanish.Startup=Inicio:
spanish.LaunchNow=Iniciar TaskGlass ahora
spanish.SettingsIcon=Configuración de TaskGlass
spanish.DeleteSettings=¿Eliminar también la configuración y los registros de TaskGlass?
spanish.StylePageTitle=Estilo de la barra de tareas
spanish.StylePageDesc=¿Cómo quieres que se vea tu barra de tareas?
spanish.StylePageSub=Puedes cambiarlo en cualquier momento desde la ventana de configuración de TaskGlass.
spanish.StyleClear=Transparente: totalmente transparente
spanish.StyleBlur=Desenfoque: fondo desenfocado con un tinte suave
spanish.StyleAcrylic=Acrílico: desenfoque con textura esmerilada
spanish.StyleOpaque=Sólido: un color liso
spanish.StyleNormal=Normal: dejar la barra de tareas de Windows sin cambios
spanish.LangPageTitle=Idioma
spanish.LangPageDesc=¿Qué idioma debe usar TaskGlass?
spanish.LangPageSub=Se aplica al menú de la bandeja y a la ventana de configuración.
spanish.LangAuto=Automático (seguir a Windows)
spanish.LangEn=English
spanish.LangEs=Español
spanish.FeatPageTitle=Funciones
spanish.FeatPageDesc=Elige los extras que quieras.
spanish.FeatPageSub=Todas se pueden cambiar después en la ventana de configuración.
spanish.FeatFlyout=Panel de iconos ocultos translúcido (el botón ^)
spanish.FeatHotkey=Atajo global Ctrl+Alt+G para activar o desactivar TaskGlass
spanish.FeatXaml=Barra totalmente transparente en Windows 11 (recomendado)

[Files]
Source: "{#DistDir}\{#AppExe}"; DestDir: "{app}"; Flags: ignoreversion
Source: "{#DistDir}\taskglass_tap.dll"; DestDir: "{app}"; Flags: ignoreversion
Source: "..\LICENSE"; DestDir: "{app}"; Flags: ignoreversion
Source: "..\LICENSE.es.txt"; DestDir: "{app}"; Flags: ignoreversion
Source: "..\NOTICE"; DestDir: "{app}"; Flags: ignoreversion
Source: "..\README.md"; DestDir: "{app}"; Flags: ignoreversion
Source: "..\README.es.md"; DestDir: "{app}"; Flags: ignoreversion
Source: "..\docs\config.schema.json"; DestDir: "{app}"; Flags: ignoreversion

[Icons]
Name: "{group}\{#AppName}"; Filename: "{app}\{#AppExe}"
Name: "{group}\{cm:SettingsIcon}"; Filename: "{app}\{#AppExe}"; Parameters: "--settings"
Name: "{group}\{cm:UninstallProgram,{#AppName}}"; Filename: "{uninstallexe}"
Name: "{autodesktop}\{#AppName}"; Filename: "{app}\{#AppExe}"; Tasks: desktopicon

[Registry]
; Same value the app's own "Start with Windows" toggle writes (the path in plain
; double quotes, single backslashes).
Root: HKCU; Subkey: "Software\Microsoft\Windows\CurrentVersion\Run"; ValueType: string; ValueName: "TaskGlass"; ValueData: """{app}\{#AppExe}"""; Flags: uninsdeletevalue; Tasks: autostart
; Clear a "disabled in Task Manager" mark left by an earlier install, which would
; stop the entry from launching even though it is registered.
Root: HKCU; Subkey: "Software\Microsoft\Windows\CurrentVersion\Explorer\StartupApproved\Run"; ValueName: "TaskGlass"; Flags: deletevalue; Tasks: autostart

[Run]
Filename: "{app}\{#AppExe}"; Description: "{cm:LaunchNow}"; Flags: nowait postinstall skipifsilent

[UninstallRun]
; Ask the running instance to quit so it restores the stock taskbar, then
; --restore as a safety net in case it was not running or was stuck.
Filename: "{app}\{#AppExe}"; Parameters: "--quit"; Flags: runhidden waituntilterminated; RunOnceId: "TaskGlassQuit"
Filename: "{app}\{#AppExe}"; Parameters: "--restore"; Flags: runhidden waituntilterminated; RunOnceId: "TaskGlassRestore"

[UninstallDelete]
; Cached copies of the helper DLL (loaded into Explorer until it restarts, so
; deletion may fail silently; that is harmless).
Type: filesandordirs; Name: "{localappdata}\TaskGlass\tap"

[Code]
var
  StylePage, LangPage, FeatPage: TInputOptionWizardPage;

function ConfigPath(): String;
begin
  Result := ExpandConstant('{userappdata}\TaskGlass\config.json');
end;

{ The questions only make sense on a first install: an existing config.json
  (upgrade or reinstall) is never touched. }
function FirstInstall(): Boolean;
begin
  Result := not FileExists(ConfigPath());
end;

procedure InitializeWizard();
begin
  StylePage := CreateInputOptionPage(wpSelectTasks, CustomMessage('StylePageTitle'), CustomMessage('StylePageDesc'),
    CustomMessage('StylePageSub'), True, False);
  StylePage.Add(CustomMessage('StyleClear'));
  StylePage.Add(CustomMessage('StyleBlur'));
  StylePage.Add(CustomMessage('StyleAcrylic'));
  StylePage.Add(CustomMessage('StyleOpaque'));
  StylePage.Add(CustomMessage('StyleNormal'));
  StylePage.SelectedValueIndex := 0;

  LangPage := CreateInputOptionPage(StylePage.ID, CustomMessage('LangPageTitle'), CustomMessage('LangPageDesc'),
    CustomMessage('LangPageSub'), True, False);
  LangPage.Add(CustomMessage('LangAuto'));
  LangPage.Add(CustomMessage('LangEn'));
  LangPage.Add(CustomMessage('LangEs'));
  { Default to the language the installer itself is running in. }
  if ActiveLanguage = 'spanish' then
    LangPage.SelectedValueIndex := 2
  else
    LangPage.SelectedValueIndex := 1;

  FeatPage := CreateInputOptionPage(LangPage.ID, CustomMessage('FeatPageTitle'), CustomMessage('FeatPageDesc'),
    CustomMessage('FeatPageSub'), False, False);
  FeatPage.Add(CustomMessage('FeatFlyout'));
  FeatPage.Add(CustomMessage('FeatHotkey'));
  FeatPage.Add(CustomMessage('FeatXaml'));
  FeatPage.Values[0] := True;
  FeatPage.Values[1] := False;
  FeatPage.Values[2] := True;
end;

function ShouldSkipPage(PageID: Integer): Boolean;
begin
  Result := False;
  if (PageID = StylePage.ID) or (PageID = LangPage.ID) or (PageID = FeatPage.ID) then
    Result := not FirstInstall();
end;

function BoolJson(B: Boolean): String;
begin
  if B then Result := 'true' else Result := 'false';
end;

{ Writes the first config.json from the wizard answers. Keys left out keep their
  built-in defaults, so the file stays valid across versions. }
procedure WriteInitialConfig();
var
  Accent, Lang, Json: String;
  Opacity: Integer;
begin
  case StylePage.SelectedValueIndex of
    0: begin Accent := 'clear';   Opacity := 0;   end;
    1: begin Accent := 'blur';    Opacity := 35;  end;
    2: begin Accent := 'acrylic'; Opacity := 55;  end;
    3: begin Accent := 'opaque';  Opacity := 100; end;
  else
    begin Accent := 'normal'; Opacity := 100; end;
  end;
  case LangPage.SelectedValueIndex of
    1: Lang := 'en';
    2: Lang := 'es';
  else
    Lang := 'auto';
  end;

  Json := '{' + #13#10 +
    '  "version": 1,' + #13#10 +
    '  "language": "' + Lang + '",' + #13#10 +
    '  "default": { "accent": "' + Accent + '", "color": "auto", "opacity": ' + IntToStr(Opacity) + ' },' + #13#10 +
    '  "surfaces": { "tray_overflow": { "enabled": ' + BoolJson(FeatPage.Values[0]) + ', "opacity": 10 } },' + #13#10 +
    '  "hotkey_enabled": ' + BoolJson(FeatPage.Values[1]) + ',' + #13#10 +
    '  "xaml_fix": ' + BoolJson(FeatPage.Values[2]) + #13#10 +
    '}' + #13#10;

  ForceDirectories(ExtractFilePath(ConfigPath()));
  SaveStringToFile(ConfigPath(), Json, False);
end;

procedure CurStepChanged(CurStep: TSetupStep);
begin
  if (CurStep = ssPostInstall) and (not WizardSilent) and FirstInstall() then
    WriteInitialConfig();
end;

procedure CurUninstallStepChanged(CurUninstallStep: TUninstallStep);
begin
  { User settings are kept unless the user asks otherwise. }
  if CurUninstallStep = usPostUninstall then
    if MsgBox(CustomMessage('DeleteSettings'), mbConfirmation, MB_YESNO or MB_DEFBUTTON2) = IDYES then
    begin
      DelTree(ExpandConstant('{userappdata}\TaskGlass'), True, True, True);
      DelTree(ExpandConstant('{localappdata}\TaskGlass'), True, True, True);
    end;
end;
