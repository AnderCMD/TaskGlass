; TaskGlass installer (Inno Setup 6). Build with scripts/build-installer.ps1.
; Per-user install: no administrator rights, nothing outside the user profile.

#ifndef AppVersion
  #define AppVersion "0.0.0"
#endif
#ifndef DistDir
  #define DistDir "..\dist"
#endif

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
LicenseFile=..\LICENSE
SetupIconFile=..\assets\taskglass.ico
UninstallDisplayIcon={app}\{#AppExe}
OutputDir={#DistDir}
OutputBaseFilename=TaskGlass-{#AppVersion}-Setup
Compression=lzma2/ultra
SolidCompression=yes
WizardStyle=modern
CloseApplications=yes
RestartApplications=no
VersionInfoVersion={#AppVersion}
VersionInfoProductName={#AppName}

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"
Name: "spanish"; MessagesFile: "compiler:Languages\Spanish.isl"

[Tasks]
Name: "autostart"; Description: "{cm:AutoStart}"; GroupDescription: "{cm:Startup}"
Name: "desktopicon"; Description: "{cm:CreateDesktopIcon}"; GroupDescription: "{cm:AdditionalIcons}"; Flags: unchecked

[CustomMessages]
english.AutoStart=Start TaskGlass when I sign in to Windows
english.Startup=Startup:
english.LaunchNow=Launch TaskGlass now
spanish.AutoStart=Iniciar TaskGlass al iniciar sesión en Windows
spanish.Startup=Inicio:
spanish.LaunchNow=Iniciar TaskGlass ahora

[Files]
Source: "{#DistDir}\{#AppExe}"; DestDir: "{app}"; Flags: ignoreversion
Source: "{#DistDir}\taskglass_tap.dll"; DestDir: "{app}"; Flags: ignoreversion
Source: "..\LICENSE"; DestDir: "{app}"; Flags: ignoreversion
Source: "..\README.md"; DestDir: "{app}"; Flags: ignoreversion
Source: "..\docs\config.schema.json"; DestDir: "{app}"; Flags: ignoreversion

[Icons]
Name: "{group}\{#AppName}"; Filename: "{app}\{#AppExe}"
Name: "{group}\{cm:UninstallProgram,{#AppName}}"; Filename: "{uninstallexe}"
Name: "{autodesktop}\{#AppName}"; Filename: "{app}\{#AppExe}"; Tasks: desktopicon

[Registry]
; Same value the app's own "Start with Windows" toggle writes (quoted path).
Root: HKCU; Subkey: "Software\Microsoft\Windows\CurrentVersion\Run"; ValueType: string; ValueName: "TaskGlass"; ValueData: """{app}\{#AppExe}"""; Flags: uninsdeletevalue; Tasks: autostart

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
procedure CurUninstallStepChanged(CurUninstallStep: TUninstallStep);
begin
  { User settings are kept unless the user asks otherwise. }
  if CurUninstallStep = usPostUninstall then
    if MsgBox(ExpandConstant('Also delete your TaskGlass settings and logs?'), mbConfirmation, MB_YESNO or MB_DEFBUTTON2) = IDYES then
    begin
      DelTree(ExpandConstant('{userappdata}\TaskGlass'), True, True, True);
      DelTree(ExpandConstant('{localappdata}\TaskGlass'), True, True, True);
    end;
end;
