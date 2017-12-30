echo on

go get -d github.com/mackerelio/go-check-plugins/...
go get -d github.com/mackerelio/mackerel-agent-plugins/...
go get -d github.com/mackerelio/mkr

setlocal enabledelayedexpansion
FOR /F %%w in (.\wix\pluginlist.txt) DO (
  go list %%w
  if not "!ERRORLEVEL!" == "0" (
    exit /b 1
  )
)
