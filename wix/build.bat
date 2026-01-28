echo on

IF NOT DEFINED WIX (
  ECHO Environment variable "WIX" not set
  EXIT /B
)

SET VERSION="%1"

CD %~dp0

if exist mackerel-agent.wxs del /F mackerel-agent.wxs
..\build\generate_wxs.exe -templateFile mackerel-agent.wxs.template -outputFile mackerel-agent.wxs -buildDir ..\build\ -productVersion %VERSION% -platform %PLATFORM_ID%
"%WIX%bin\candle.exe" -ext WixUIExtension -ext WixUtilExtension -arch %PLATFORM_ID% mackerel-agent.wxs
"%WIX%bin\light.exe" -ext WixUIExtension -ext WixUtilExtension -out "..\build\mackerel-agent%MSI_SUFFIX%.msi" mackerel-agent.wixobj

if exist mackerel-agent-kcps.wxs del /F mackerel-agent-kcps.wxs
..\build\generate_wxs.exe -templateFile mackerel-agent.wxs.template -outputFile mackerel-agent-kcps.wxs -buildDir ..\build\ -productVersion %VERSION% -platform %PLATFORM_ID% -configFile mackerel-agent-kcps.sample.conf
"%WIX%bin\candle.exe" -ext WixUIExtension -ext WixUtilExtension -arch %PLATFORM_ID% mackerel-agent-kcps.wxs
"%WIX%bin\light.exe" -ext WixUIExtension -ext WixUtilExtension -out "..\build\mackerel-agent-kcps%MSI_SUFFIX%.msi" mackerel-agent-kcps.wixobj

copy ..\build\mackerel-agent-kcps.exe ..\build\mackerel-agent.exe
"%WIX%bin\light.exe" -ext WixUIExtension -ext WixUtilExtension -out "..\build\mackerel-agent-k%MSI_SUFFIX%.msi" mackerel-agent.wixobj
