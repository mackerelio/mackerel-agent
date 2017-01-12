echo on

go get -d github.com/mackerelio/go-check-plugins/...
go get -d github.com/mackerelio/mackerel-agent-plugins/...

FOR /F %%w in (.\wix\pluginlist.txt) DO (
  go list %%W
  if not "%ERRORLEVEL%" == "0" (
    exit /b 1
  )
)
