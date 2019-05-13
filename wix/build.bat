echo on

IF NOT DEFINED WIX (
  ECHO Environment variable "WIX" not set
  EXIT /B
)

CD %~dp0\..

CALL build.bat
CALL build-k.bat

FOR /F %%w in (.\wix\pluginlist.txt) DO (
  go build -o build\%%~nw.exe %%w
)

CD %~dp0

go build -o ..\build\wrapper.exe wrapper\wrapper_windows.go wrapper\install.go
go build -o ..\build\replace.exe replace\replace_windows.go
go build -o ..\build\generate_wxs.exe generate_wxs\generate_wxs.go

REM retrieve numeric version from git tag
FOR /F "usebackq" %%w IN (`git tag -l --sort=-version:refname "v*"`) DO (
  IF NOT DEFINED VERSION (
    SET VERSION=%%w
  )
)
set VERSION=%VERSION:v=%
FOR /F "tokens=1 delims=-+" %%w IN ('ECHO %%VERSION%%') DO SET VERSION=%%w
IF "%VERSION%"=="staging" (
  EXIT /B
)

if exist mackerel-agent.wxs del /F mackerel-agent.wxs
..\build\generate_wxs.exe -templateFile mackerel-agent.wxs.template -outputFile mackerel-agent.wxs -buildDir ..\build\ -productVersion "%VERSION%" -platform %PLATFORM_ID%

"%WIX%bin\candle.exe" -ext WixUIExtension -ext WixUtilExtension -arch %PLATFORM_ID% mackerel-agent.wxs
"%WIX%bin\light.exe" -ext WixUIExtension -ext WixUtilExtension -out "..\build\mackerel-agent%MSI_SUFFIX%.msi" mackerel-agent.wixobj
copy ..\build\mackerel-agent-kcps.exe ..\build\mackerel-agent.exe
"%WIX%bin\light.exe" -ext WixUIExtension -ext WixUtilExtension -out "..\build\mackerel-agent-k%MSI_SUFFIX%.msi" mackerel-agent.wixobj

REM code signing if build on tags
if defined APPVEYOR_REPO_TAG_NAME (
  certutil -f -decode c:\mackerel-secure\cert.p12.base64 c:\mackerel-secure\cert.p12

  FOR /F "usebackq" %%w IN (`type c:\mackerel-secure\certpass`) DO "%SIGNTOOL%" sign /fd sha256 /t "http://timestamp.verisign.com/scripts/timestamp.dll" /f "c:\mackerel-secure\cert.p12" /p "%%w" /v "..\build\mackerel-agent%MSI_SUFFIX%.msi"

  FOR /F "usebackq" %%w IN (`type c:\mackerel-secure\certpass`) DO "%SIGNTOOL%" sign /fd sha256 /t "http://timestamp.verisign.com/scripts/timestamp.dll" /f "c:\mackerel-secure\cert.p12" /p "%%w" /v "..\build\mackerel-agent-k%MSI_SUFFIX%.msi"
)
