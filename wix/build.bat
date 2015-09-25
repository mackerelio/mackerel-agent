go get -d -v -t ./...

pushd %0\..\..

call build.bat

pushd %0\..

go get github.com/mackerelio/mackerel-agent/wix/wrapper
go get github.com/mackerelio/mackerel-agent/wix/replace

go build -o ..\build\wrapper.exe wrapper\wrapper_windows.go
go build -o ..\build\replace.exe replace\replace_windows.go

REM retrieve numeric version from git tag
FOR /F "usebackq" %%w IN (`git describe --tags --abbrev^=0`) DO SET VERSION=%%w
set VERSION=%VERSION:v=%
FOR /F "tokens=1 delims=-+" %%w IN ('ECHO %%VERSION%%') DO SET VERSION=%%w

del /F mackerel-agent.wxs
..\build\replace.exe mackerel-agent.wxs.template mackerel-agent.wxs "___VERSION___" "%VERSION%"

"%WIX%bin\candle.exe" mackerel-agent.wxs
"%WIX%bin\light.exe" -ext WixUIExtension -out "..\build\mackerel-agent.msi" mackerel-agent.wixobj

REM code signing if build on tags
REM if defined APPVEYOR_REPO_TAG_NAME (
  certutil -f -decode c:\mackerel-secure\cert.p12.base64 c:\mackerel-secure\cert.p12

  REM ref. https://support.comodo.com/index.php?/Default/Knowledgebase/Article/View/68/7/time-stamping-server
  FOR /F "usebackq" %%w IN (`type c:\mackerel-secure\certpass`) DO "%SIGNTOOL%" sign /t "http://timestamp.comodoca.com/authenticode" /f "c:\mackerel-secure\cert.p12" /p "%%w" /v "..\build\mackerel-agent.msi"
REM )
