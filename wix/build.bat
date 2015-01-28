pushd %0\..\..

call build.bat

pushd %0\..

go get github.com/mackerelio/mackerel-agent/wix/wrapper
go get github.com/mackerelio/mackerel-agent/wix/replace

go build -o ..\build\wrapper.exe wrapper\wrapper.go
go build -o ..\build\replace.exe replace\replace.go

"%WIX%bin\candle.exe" mackerel-agent.wxs
"%WIX%bin\light.exe" -out "..\build\mackerel-agent.msi" mackerel-agent.wixobj

pause
