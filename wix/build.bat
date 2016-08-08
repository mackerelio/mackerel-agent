echo on

go get -d -v -t ./...

pushd %0\..\..

call build.bat
call build-k.bat

pushd %0\..

go get github.com/mackerelio/mackerel-agent/wix/wrapper
go get github.com/mackerelio/mackerel-agent/wix/generate_wxs

go build -o ..\build\wrapper.exe wrapper\wrapper_windows.go
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

del /F mackerel-agent.wxs
..\build\generate_wxs.exe -templateFile mackerel-agent.wxs.template -outputFile mackerel-agent.wxs -pluginDir ..\..\go-check-plugins\build\ -productVersion "%VERSION%"

"%WIX%bin\candle.exe" mackerel-agent.wxs
"%WIX%bin\light.exe" -ext WixUIExtension -out "..\build\mackerel-agent.msi" mackerel-agent.wixobj
copy ..\build\mackerel-agent-kcps.exe ..\build\mackerel-agent.exe
"%WIX%bin\light.exe" -ext WixUIExtension -out "..\build\mackerel-agent-k.msi" mackerel-agent.wixobj

REM code signing if build on tags
if defined APPVEYOR_REPO_TAG_NAME (
  certutil -f -decode c:\mackerel-secure\cert.p12.base64 c:\mackerel-secure\cert.p12

  FOR /F "usebackq" %%w IN (`type c:\mackerel-secure\certpass`) DO "%SIGNTOOL%" sign /fd sha256 /t "http://timestamp.verisign.com/scripts/timestamp.dll" /f "c:\mackerel-secure\cert.p12" /p "%%w" /v "..\build\mackerel-agent.msi"

  FOR /F "usebackq" %%w IN (`type c:\mackerel-secure\certpass`) DO "%SIGNTOOL%" sign /fd sha256 /t "http://timestamp.verisign.com/scripts/timestamp.dll" /f "c:\mackerel-secure\cert.p12" /p "%%w" /v "..\build\mackerel-agent-k.msi"
)
