[Setup]
AppId={{9CFAEC95-9F6A-4281-B815-4DC9AB3DB328}
AppName=Shelter
AppVersion=0.1
AppVerName=Shelter 0.1
AppPublisher=Rafael Dantas Justo
AppPublisherURL=http://github.com/rafaeljusto/shelter
AppSupportURL=http://github.com/rafaeljusto/shelter
AppUpdatesURL=http://github.com/rafaeljusto/shelter
DefaultDirName={pf}\shelter
DefaultGroupName=shelter
AllowNoIcons=yes
LicenseFile=..\..\LICENSE
OutputBaseFilename=setup
Compression=lzma
SolidCompression=yes

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"

[Tasks]
Name: "desktopicon"; Description: "{cm:CreateDesktopIcon}"; GroupDescription: "{cm:AdditionalIcons}"; Flags: unchecked

[Dirs]
Name: "{app}\templates\client"; Permissions: users-modify
Name: "{app}\templates\notification"; Permissions: users-modify
Name: "{app}\etc"; Permissions: users-modify

[Files]
Source: "..\..\shelter.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: "easyconf.exe"; DestDir: "{app}"; Flags: ignoreversion
Source: "..\..\etc\shelter.conf.windows.sample"; DestDir: "{app}\etc\shelter.conf.windows.sample"; Flags: ignoreversion

[Icons]
Name: "{group}\shelter"; Filename: "{app}\shelter.exe"
Name: "{commondesktop}\shelter"; Filename: "{app}\shelter.exe"; Tasks: desktopicon

[Run]
Filename: "{app}\easyconf.exe -sample={app}\etc\shelter.conf.windows.sample"