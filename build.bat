echo on

FOR /F "usebackq" %%w IN (`git rev-parse --short HEAD`) DO SET COMMIT=%%w

go build -o build/mackerel-agent.exe -ldflags="-X main.gitcommit=%COMMIT% " github.com/mackerelio/mackerel-agent
