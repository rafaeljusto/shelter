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
LicenseFile={src}\LICENSE
OutputBaseFilename=setup
Compression=lzma
SolidCompression=yes

[Languages]
Name: "english"; MessagesFile: "compiler:Default.isl"

[Tasks]
Name: "desktopicon"; Description: "{cm:CreateDesktopIcon}"; GroupDescription: "{cm:AdditionalIcons}"; Flags: unchecked

[Files]
Source: "{src}\shelter.exe"; DestDir: "{app}"; Flags: ignoreversion

[Icons]
Name: "{group}\shelter"; Filename: "{app}\shelter.exe"
Name: "{commondesktop}\shelter"; Filename: "{app}\shelter.exe"; Tasks: desktopicon