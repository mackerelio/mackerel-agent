echo on

setlocal enabledelayedexpansion
FOR /F %%w in (.\wix\pluginlist.txt) DO (
  go list %%w
  if not "!ERRORLEVEL!" == "0" (
    exit /b 1
  )
)
