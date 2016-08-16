echo on

FOR /F "usebackq" %%w IN (`git rev-parse --short HEAD`) DO SET COMMIT=%%w

FOR /F "usebackq" %%w IN (`git tag -l --sort=-version:refname "v*"`) DO (
  IF NOT DEFINED VERSION (
    SET VERSION=%%w
  )
)

set VERSION=%VERSION:v=%

echo %VERSION%

go build -o build/mackerel-agent.exe -ldflags="-X github.com/mackerelio/mackerel-agent/version.GITCOMMIT=%COMMIT% -X github.com/mackerelio/mackerel-agent/version.VERSION=%VERSION% " github.com/mackerelio/mackerel-agent
