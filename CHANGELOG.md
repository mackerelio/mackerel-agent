# Changelog

## 0.85.0 (2025-05-16)

* update mackerelio/* on wix/ #1085 (yseto)
* Remove old build tags. #1079 (yseto)
* introduce interfaces/use_adapter flag to take metrics from Network Adapter #1077 (kmuto)


## 0.84.3 (2025-03-31)

* update github.com/mackerelio/go-check-plugins, mkr on /wix #1072 (yseto)
* Bump github.com/mackerelio/mackerel-client-go from 0.35.0 to 0.36.0 #1070 (dependabot[bot])


## 0.84.2 (2025-03-31)

* update Go #1068 (yseto)
* replace to newer runner-images #1067 (yseto)
* Bump golang.org/x/net from 0.33.0 to 0.36.0 in /wix #1066 (dependabot[bot])
* Bump golang.org/x/text from 0.14.0 to 0.23.0 #1065 (dependabot[bot])
* Bump github.com/mackerelio/mkr from 0.59.0 to 0.59.2 in /wix #1064 (dependabot[bot])
* Bump golang.org/x/sys from 0.29.0 to 0.31.0 in /wix #1063 (dependabot[bot])
* Bump github.com/mackerelio/mackerel-agent-plugins from 0.86.0 to 0.88.0 in /wix #1062 (dependabot[bot])
* use ghcr.io/mackerelio/mackerel-rpm-builder #1060 (yseto)
* Bump mackerelio/workflows from 1.3.0 to 1.4.0 #1059 (dependabot[bot])
* Bump github.com/mackerelio/mackerel-client-go from 0.34.0 to 0.35.0 #1056 (dependabot[bot])


## 0.84.1 (2025-03-03)

* Set the default path by "mkr plugin install" to the PATH environment. #1058 (fujiwara)
* Bump mackerelio/workflows from 1.2.0 to 1.3.0 #1052 (dependabot[bot])


## 0.84.0 (2025-01-27)

* Bump golang.org/x/net from 0.25.0 to 0.33.0 in /wix #1050 (dependabot[bot])
* Bump azure/trusted-signing-action from 0.3.20 to 0.5.1 #1047 (dependabot[bot])
* Bump golang.org/x/sys from 0.21.0 to 0.29.0 in /wix #1046 (dependabot[bot])
* Pined build actions runner version #1045 (appare45)
* Bump golang.org/x/crypto from 0.24.0 to 0.31.0 in /wix #1044 (dependabot[bot])
* use mackerelio/workflows@v1.2.0 #1043 (yseto)
* Bump github.com/vishvananda/netlink from 1.1.0 to 1.3.0 #1037 (dependabot[bot])
* Bump github.com/fatih/color from 1.16.0 to 1.18.0 #1036 (dependabot[bot])
* Add disable_http_keep_alive option #1035 (appare45)
* check status code in go.yaml #1033 (rmatsuoka)
* Bump github.com/mackerelio/mackerel-client-go from 0.31.0 to 0.34.0 #1023 (dependabot[bot])
* Bump github.com/Songmu/gocredits from 0.3.0 to 0.3.1 #1020 (dependabot[bot])
* Bump github.com/mackerelio/go-osstat from 0.2.4 to 0.2.5 #1011 (dependabot[bot])
* Bump github.com/shirou/gopsutil/v3 from 3.24.3 to 3.24.5 #1008 (dependabot[bot])


## 0.83.0 (2024-11-26)

* Bump github.com/mackerelio/mkr from 0.58.0 to 0.59.0 in /wix #1029 (dependabot[bot])
* Bump github.com/mackerelio/mackerel-agent-plugins from 0.82.1 to 0.86.0 in /wix #1028 (dependabot[bot])
* add mackerel-plugin-snmp and check-ping to Windows build. #1027 (kmuto)
* fix lint error #1025 (rmatsuoka)
* Write out error message when plugin seems have a bug #1024 (rmatsuoka)
* Bump github.com/mackerelio/go-check-plugins from 0.46.3 to 0.47.0 in /wix #1022 (dependabot[bot])
* Use Azure Trusted Signing for Windows installer code signing #1021 (yohfee)


## 0.82.0 (2024-06-13)

* go mod tidy on wix/ #1015 (yseto)
* go get -u github.com/mackerelio/mkr #1014 (yseto)
* Update mackerel/docker-mackerel-rpm-builder #1009 (yseto)


## 0.81.0 (2024-04-25)

* Update wix go-check-plugins v0.46.3, mackerel-agent-plugins v0.82.1 #996 (ne-sachirou)
* Bump golang.org/x/net from 0.22.0 to 0.23.0 in /wix #995 (dependabot[bot])
* Ignore errors of getting information of a disk on Windows #994 (Arthur1)
* Bump github.com/mackerelio/mackerel-client-go from 0.30.0 to 0.31.0 #993 (dependabot[bot])
* Bump github.com/shirou/gopsutil/v3 from 3.24.2 to 3.24.3 #992 (dependabot[bot])
* Bump github.com/mackerelio/mackerel-agent-plugins from 0.81.0 to 0.82.0 in /wix #991 (dependabot[bot])
* Bump golang.org/x/sys from 0.18.0 to 0.19.0 in /wix #990 (dependabot[bot])
* Bump peter-evans/repository-dispatch from 2 to 3 #965 (dependabot[bot])
* Bump actions/cache from 3 to 4 #964 (dependabot[bot])
* Bump actions/download-artifact from 3 to 4 #956 (dependabot[bot])
* Bump actions/upload-artifact from 3 to 4 #955 (dependabot[bot])


## 0.80.0 (2024-03-19)

* Bump github.com/mackerelio/mkr from 0.55.0 to 0.57.0 in /wix #987 (dependabot[bot])
* Bump github.com/mackerelio/mackerel-client-go from 0.29.0 to 0.30.0 #986 (dependabot[bot])
* Bump google.golang.org/protobuf from 1.31.0 to 1.33.0 in /wix #985 (dependabot[bot])
* refactor: lint and test workflows #973 (lufia)


## 0.79.0 (2024-03-06)

* Bump github.com/mackerelio/mackerel-agent-plugins from 0.79.0 to 0.81.0 in /wix #981 (dependabot[bot])
* Bump github.com/shirou/gopsutil/v3 from 3.23.8 to 3.24.2 #979 (dependabot[bot])
* Bump github.com/mackerelio/go-check-plugins from 0.46.1 to 0.46.2 in /wix #978 (dependabot[bot])
* Bump github.com/mackerelio/mkr from 0.54.0 to 0.55.0 in /wix #977 (dependabot[bot])
* fix build errors on *BSD platforms #975 (lufia)
* command: drop `go test -short` mode #972 (lufia)
* replace interface{} to any #971 (Arthur1)
* Use shirou/gopsutil for getting CPU info on Linux #970 (Arthur1)


## 0.78.3 (2023-12-25)

* Bump github.com/mackerelio/mkr from 0.53.0 to 0.54.0 in /wix #959 (dependabot[bot])
* Bump github.com/mackerelio/mackerel-client-go from 0.26.0 to 0.29.0 #958 (dependabot[bot])
* Bump actions/setup-go from 4 to 5 #954 (dependabot[bot])
* Bump actions/github-script from 6 to 7 #948 (dependabot[bot])
* Bump github.com/fatih/color from 1.15.0 to 1.16.0 #944 (dependabot[bot])
* Bump golang.org/x/text from 0.13.0 to 0.14.0 #943 (dependabot[bot])


## 0.78.2 (2023-12-01)

* Avoid iowait percentage's overflow when counter is reset #947 (Arthur1)


## 0.78.1 (2023-11-15)

* Bump github.com/mackerelio/go-check-plugins from 0.45.0 to 0.46.1 in /wix #945 (dependabot[bot])
* Bump golang.org/x/sys from 0.12.0 to 0.14.0 in /wix #941 (dependabot[bot])
* Bump golang.org/x/net from 0.15.0 to 0.17.0 in /wix #940 (dependabot[bot])


## 0.78.0 (2023-09-22)

* Bump github.com/mackerelio/mkr from 0.51.0 to 0.53.0 in /wix #936 (dependabot[bot])
* Bump github.com/mackerelio/mackerel-agent-plugins from 0.78.0 to 0.79.0 in /wix #935 (dependabot[bot])
* Bump github.com/mackerelio/go-check-plugins from 0.44.0 to 0.45.0 in /wix #934 (dependabot[bot])
* update go version 1.19.x -> 1.20.x on GitHub Action #933 (rmatsuoka)
* Bump actions/checkout from 3 to 4 #932 (dependabot[bot])
* Bump github.com/shirou/gopsutil/v3 from 3.23.6 to 3.23.8 #930 (dependabot[bot])
* Bump golang.org/x/text from 0.9.0 to 0.13.0 #929 (dependabot[bot])
* Bump golang.org/x/sys from 0.7.0 to 0.12.0 in /wix #928 (dependabot[bot])
* remove comment that differ from current plugin implementation #927 (kga)
* Remove old rpm packaging #926 (yseto)
* Bump github.com/BurntSushi/toml from 1.3.0 to 1.3.2 #910 (dependabot[bot])
* Bump actions/setup-go from 3 to 4 #887 (dependabot[bot])


## 0.77.1 (2023-07-13)

* Bump github.com/shirou/gopsutil/v3 from 3.23.3 to 3.23.6 #915 (dependabot[bot])
* Bump github.com/mackerelio/mkr from 0.49.2 to 0.51.0 in /wix #908 (dependabot[bot])
* Bump github.com/mackerelio/mackerel-client-go from 0.24.0 to 0.26.0 #900 (dependabot[bot])


## 0.77.0 (2023-06-14)

* fixed to log the key when metric parsing fails #909 (tukaelu)
* Bump github.com/BurntSushi/toml from 1.2.1 to 1.3.0 #907 (dependabot[bot])


## 0.76.0 (2023-05-31)

* added check-file-age on windows #903 (yseto)
* update Actions Runner Image of windows. #899 (yseto)
* add feature to ignore disk labels by regexp in 'disk' of system metrics #898 (kmuto)


## 0.75.2 (2023-04-12)

* Bump golang.org/x/sys from 0.5.0 to 0.7.0 in /wix #891 (dependabot[bot])
* Bump golang.org/x/text from 0.7.0 to 0.9.0 #890 (dependabot[bot])
* Bump github.com/shirou/gopsutil/v3 from 3.23.1 to 3.23.3 #889 (dependabot[bot])
* Bump github.com/mackerelio/mackerel-agent-plugins from 0.75.0 to 0.78.0 in /wix #888 (dependabot[bot])
* Bump github.com/mackerelio/go-osstat from 0.2.3 to 0.2.4 #886 (dependabot[bot])
* Bump github.com/fatih/color from 1.14.1 to 1.15.0 #885 (dependabot[bot])
* Bump github.com/mackerelio/mkr from 0.49.0 to 0.49.2 in /wix #883 (dependabot[bot])
* Bump github.com/mackerelio/go-check-plugins from 0.43.0 to 0.44.0 in /wix #882 (dependabot[bot])


## 0.75.1 (2023-02-15)

* Bump golang.org/x/sys from 0.4.0 to 0.5.0 in /wix #871 (dependabot[bot])
* Bump golang.org/x/text from 0.6.0 to 0.7.0 #870 (dependabot[bot])
* fix msi package installscope. #869 (yseto)
* Bump github.com/mackerelio/mkr from 0.48.0 to 0.49.0 in /wix #868 (dependabot[bot])
* Bump github.com/mackerelio/go-check-plugins from 0.42.4 to 0.43.0 in /wix #867 (dependabot[bot])
* Bump github.com/shirou/gopsutil/v3 from 3.22.12 to 3.23.1 #865 (dependabot[bot])


## 0.75.0 (2023-02-01)

* Bump peter-evans/repository-dispatch from 1 to 2 #863 (dependabot[bot])
* Bump golang.org/x/sys from 0.3.0 to 0.4.0 in /wix #862 (dependabot[bot])
* Bump github.com/fatih/color from 1.13.0 to 1.14.1 #861 (dependabot[bot])
* Enables Dependabot version updates for GitHub Actions #860 (Arthur1)
* Upgrade reusable actions 2 #859 (Arthur1)
* Upgrade reusable actions #858 (Arthur1)
* Remove build command for apt v1 #857 (Arthur1)
* Bump github.com/mackerelio/mackerel-agent-plugins from 0.72.1 to 0.75.0 in /wix #856 (dependabot[bot])
* Bump github.com/mackerelio/mkr from 0.47.1 to 0.48.0 in /wix #855 (dependabot[bot])
* Bump github.com/mackerelio/go-check-plugins from 0.42.1 to 0.42.4 in /wix #854 (dependabot[bot])
* Get osname, version, and release by WMI #852 (Arthur1)
* Bump github.com/mackerelio/mackerel-client-go from 0.22.0 to 0.24.0 #849 (dependabot[bot])
* Bump golang.org/x/text from 0.3.7 to 0.6.0 #848 (dependabot[bot])
* Bump github.com/BurntSushi/toml from 1.2.0 to 1.2.1 #825 (dependabot[bot])


## 0.74.1 (2023-01-18)

* Replace WMI library #851 (Arthur1)
* Update some libraries #847 (yseto)
* Remove deb builder #845 (yseto)


## 0.74.0 (2022-12-20)

* [ci]fix timeout #843 (yseto)
* use ubuntu-20.04 #842 (yseto)
* refactor(config/validate): use `candidates` instead of `parentConfKeys` #839 (wafuwafu13)
* improve purge stage of Debian package to remove id file and keep custom files #837 (kmuto)
* fix(configtest): don't detect child key if parent key has already detected #832 (wafuwafu13)


## 0.73.3 (2022-11-09)

* Fix config test #830 (ryuichi1208)


## 0.73.2 (2022-11-04)

* Replace linter #821 (yseto)
* Bump github.com/mackerelio/mackerel-client-go from 0.21.2 to 0.22.0 #819 (dependabot[bot])
* Improve `mackerel-agent configtest`: Add suggestion to unexpected keys #818 (wafuwafu13)
* Bump github.com/Songmu/gocredits from 0.2.0 to 0.3.0 #817 (dependabot[bot])
* Bump github.com/mackerelio/mkr from 0.47.0 to 0.47.1 in /wix #816 (dependabot[bot])
* go.mod 1.17 -> 1.18 #814 (yseto)
* Improve `mackerel-agent configtest`: detect unexpected key #813 (wafuwafu13)
* fix deprecated function. #808 (yseto)
* Bump github.com/mackerelio/go-osstat from 0.2.2 to 0.2.3 #802 (dependabot[bot])
* Bump github.com/Songmu/goxz from 0.8.2 to 0.9.1 #800 (dependabot[bot])
* Bump github.com/BurntSushi/toml from 1.1.0 to 1.2.0 #795 (dependabot[bot])


## 0.73.1 (2022-09-14)

* Bump github.com/mackerelio/mkr from 0.46.9 to 0.47.0 in /wix #812 (dependabot[bot])
* config_test: Add the case of LoadConfigWithInvalidToml #811 (wafuwafu13)
* In the test, if the Fatal if the result is nil. #807 (yseto)
* remove unused codes #806 (yseto)
* replace io/ioutil #805 (yseto)
* Bump github.com/mackerelio/mackerel-client-go from 0.21.1 to 0.21.2 #803 (dependabot[bot])
* get interface information via netlink on linux. #801 (yseto)


## 0.73.0 (2022-07-27)

* Bump github.com/mackerelio/mkr from 0.46.8 to 0.46.9 in /wix #797 (dependabot[bot])
* Bump github.com/mackerelio/mackerel-agent-plugins from 0.71.0 to 0.72.1 in /wix #796 (dependabot[bot])
* Loosen the conditions of delaying report of check monitering. #792 (sugy)


## 0.72.15 (2022-07-20)

* Bump github.com/mackerelio/mkr from 0.46.7 to 0.46.8 in /wix #791 (dependabot[bot])
* Bump github.com/mackerelio/go-check-plugins from 0.42.0 to 0.42.1 in /wix #790 (dependabot[bot])
* Bump github.com/mackerelio/mackerel-client-go from 0.21.0 to 0.21.1 #788 (dependabot[bot])


## 0.72.14 (2022-06-22)

* add s, ms, bps to metric units #786 (Arthur1)
* Bump github.com/mackerelio/mkr from 0.46.6 to 0.46.7 in /wix #784 (dependabot[bot])


## 0.72.13 (2022-06-08)

* Bump github.com/mackerelio/mkr from 0.46.5 to 0.46.6 in /wix #782 (dependabot[bot])
* Bump github.com/Songmu/prompter from 0.5.0 to 0.5.1 #781 (dependabot[bot])


## 0.72.12 (2022-05-26)

* Bump github.com/Songmu/goxz from 0.8.1 to 0.8.2 #779 (dependabot[bot])
* Bump github.com/mackerelio/mackerel-agent-plugins from 0.70.6 to 0.71.0 in /wix #778 (dependabot[bot])


## 0.72.11 (2022-04-14)

* Bump github.com/mackerelio/mackerel-agent-plugins from 0.70.4 to 0.70.6 in /wix #776 (dependabot[bot])
* Bump github.com/mackerelio/go-check-plugins from 0.41.7 to 0.42.0 in /wix #775 (dependabot[bot])
* Bump github.com/mackerelio/mkr from 0.46.4 to 0.46.5 in /wix #774 (dependabot[bot])
* Bump github.com/BurntSushi/toml from 1.0.0 to 1.1.0 #772 (dependabot[bot])


## 0.72.10 (2022-03-30)

* Bump github.com/mackerelio/go-osstat from 0.2.1 to 0.2.2 #770 (dependabot[bot])
* Bump github.com/mackerelio/go-check-plugins from 0.41.6 to 0.41.7 in /wix #769 (dependabot[bot])
* Bump github.com/mackerelio/mackerel-agent-plugins from 0.70.3 to 0.70.4 in /wix #768 (dependabot[bot])
* Bump github.com/mackerelio/mkr from 0.46.3 to 0.46.4 in /wix #767 (dependabot[bot])


## 0.72.9 (2022-03-15)

* Bump github.com/mackerelio/go-check-plugins from 0.41.5 to 0.41.6 in /wix #765 (dependabot[bot])
* Bump github.com/mackerelio/mkr from 0.46.2 to 0.46.3 in /wix #764 (dependabot[bot])
* Bump github.com/mackerelio/mackerel-agent-plugins from 0.70.2 to 0.70.3 in /wix #763 (dependabot[bot])


## 0.72.8 (2022-02-16)

* upgrade Go 1.16 -> 1.17 #760 (lufia)
* Bump github.com/mackerelio/mkr from 0.46.1 to 0.46.2 in /wix #759 (dependabot[bot])
* Bump github.com/mackerelio/go-check-plugins from 0.41.4 to 0.41.5 in /wix #758 (dependabot[bot])


## 0.72.7 (2022-02-02)

* Bump github.com/mackerelio/mkr from 0.46.0 to 0.46.1 in /wix #756 (dependabot[bot])
* Bump github.com/mackerelio/mackerel-agent-plugins from 0.69.1 to 0.70.2 in /wix #755 (dependabot[bot])
* Bump github.com/mackerelio/go-check-plugins from 0.41.1 to 0.41.4 in /wix #754 (dependabot[bot])
* Bump github.com/BurntSushi/toml from 0.3.1 to 1.0.0 #753 (dependabot[bot])
* Bump github.com/Songmu/goxz from 0.7.0 to 0.8.1 #749 (dependabot[bot])
* Bump github.com/mackerelio/mackerel-client-go from 0.19.0 to 0.21.0 #748 (dependabot[bot])


## 0.72.6 (2022-01-12)

* Bump github.com/mackerelio/mkr from v0.45.3 to 0.46.0 #746 (susisu)


## 0.72.5 (2021-12-01)

* Bump github.com/mackerelio/mackerel-client-go from 0.17.0 to 0.19.0 #738 (dependabot[bot])


## 0.72.4 (2021-11-18)

* spec: reuse http.Client #740 (lufia)
* Bump github.com/mackerelio/go-osstat from 0.2.0 to 0.2.1 #737 (dependabot[bot])
* Bump github.com/mattn/goveralls from 0.0.9 to 0.0.11 #736 (dependabot[bot])
* make wix/ a submodule #735 (susisu)
* Add arm64/darwin build to GitHub release #734 (astj)
* read a response body even if status is not good #733 (lufia)


## 0.72.3 (2021-10-20)

* Bump github.com/mackerelio/mackerel-agent-plugins from v0.65.0 to v0.69.1 #731 (susisu)
* Bump github.com/mackerelio/go-check-plugins from v0.39.5 to v0.41.1 #730 (ne-sachirou)


## 0.72.2 (2021-09-06)

* Update Code Signing Certificates #727 (Krout0n)
* Bump github.com/mackerelio/go-check-plugins from 0.39.3 to 0.39.5 #722 (dependabot[bot])
* Bump github.com/mackerelio/mackerel-agent-plugins from 0.64.2 to 0.65.0 #721 (dependabot[bot])
* Bump github.com/mackerelio/mkr from 0.45.1 to 0.45.2 #718 (dependabot[bot])


## 0.72.1 (2021-06-23)

* Bump github.com/mackerelio/go-check-plugins from 0.39.2 to 0.39.3 #711 (dependabot[bot])
* Bump github.com/mattn/goveralls from 0.0.8 to 0.0.9 #709 (dependabot[bot])
* Bump github.com/mackerelio/mackerel-agent-plugins from 0.64.0 to 0.64.2 #714 (dependabot[bot])
* Bump github.com/mackerelio/mkr from 0.45.0 to 0.45.1 #712 (dependabot[bot])


## 0.72.0 (2021-06-17)

* fix http_proxy option in v0.71.2 #715 (yseto)


## 0.71.2 (2021-05-26)

* Bump github.com/mackerelio/go-osstat from 0.1.0 to 0.2.0 #707 (dependabot[bot])
* Bump golang.org/x/text from 0.3.5 to 0.3.6 #703 (dependabot[bot])
* Bump github.com/Songmu/goxz from 0.6.0 to 0.7.0 #704 (dependabot[bot])
* upgrade Go 1.14 to 1.16 #705 (lufia)
* [ci] avoid additional Go installation #698 (lufia)
* Bump github.com/Songmu/prompter from 0.4.0 to 0.5.0 #702 (dependabot[bot])
* Bump github.com/mackerelio/mackerel-client-go from 0.15.0 to 0.16.0 #701 (dependabot[bot])
* Bump github.com/mattn/goveralls from 0.0.7 to 0.0.8 #697 (dependabot[bot])
* Bump github.com/mackerelio/mackerel-client-go from 0.13.0 to 0.15.0 #700 (dependabot[bot])
* Bump github.com/mackerelio/golib from 1.1.0 to 1.2.0 #696 (dependabot[bot])
* [ci] fix option at repository-dispatch #699 (yseto)
* [ci] added repository_dispatch to homebrew-mackerel-agent #694 (yseto)
* [ci] replace token #695 (yseto)
* [ci] replace mackerel-github-release #692 (yseto)
* Changed the test method from TestDiskGenerator to TestParseDiskStats  because the test results are flaky. #693 (yseto)


## 0.71.1 (2021-01-21)

* remove .circleci/config.yml #689 (yseto)
* [ci] bump Windows i386 Golang to 1.14.14 #688 (astj)
* Build Windows package on GitHub Actions #684 (yseto)
* Bump github.com/mackerelio/golib from 1.0.0 to 1.1.0 #686 (dependabot[bot])
* Bump github.com/mackerelio/mackerel-client-go from 0.12.0 to 0.13.0 #687 (dependabot[bot])
* Bump golang.org/x/text from 0.3.4 to 0.3.5 #685 (dependabot[bot])


## 0.71.0 (2020-12-14)

* Bump github.com/mackerelio/mackerel-client-go from 0.11.0 to 0.12.0 #680 (dependabot[bot])
* Network interface exclusion feature #681 (yseto)
* Bump golang.org/x/text from 0.3.3 to 0.3.4 #666 (dependabot-preview[bot])


## 0.70.3 (2020-11-25)

* Build with Go 1.14 in CI (was 1.15 by mistake) #678 (astj)


## 0.70.2 (2020-11-19)

* Fix artifact filename pattern again to include mackerel-agent_{os}_{arch}.tar.gz to GitHub Release artifacts #676 (astj)


## 0.70.1 (2020-11-19)

* include mackerel-agent_{os}_{arch}.tar.gz to GitHub Release artifacts #674 (astj)


## 0.70.0 (2020-11-19)

* replace Travis CI workflow with GitHub Actions #670 (astj)
* Retry once immediately on posting metrics and check reports when error is caused by net/http (network error) #669 (astj)


## 0.69.3 (2020-10-28)

* Bump github.com/shirou/gopsutil from 2.20.8+incompatible to 2.20.9+incompatible #664 (dependabot-preview[bot])


## 0.69.2 (2020-10-01)

* Bump github.com/mackerelio/mackerel-client-go from 0.10.1 to 0.11.0 #662 (dependabot-preview[bot])


## 0.69.1 (2020-09-15)

* kcps, stage: add --target option to build rpm #660 (lufia)


## 0.69.0 (2020-09-15)

* revert changing filename unexpectedly #658 (lufia)
* Bump github.com/shirou/gopsutil from 2.20.6+incompatible to 2.20.8+incompatible #656 (dependabot-preview[bot])
* Bump github.com/mattn/goveralls from 0.0.6 to 0.0.7 #657 (dependabot-preview[bot])
* Bump github.com/Songmu/prompter from 0.3.0 to 0.4.0 #655 (dependabot-preview[bot])
* revert mkdir with shell expansion #654 (lufia)
* add arm64 RPM packages, and change deb architecture to be correct #653 (lufia)
* update go: 1.12 -> 1.14 #648 (lufia)


## 0.68.2 (2020-07-29)

* Bump github.com/shirou/gopsutil from 2.20.4+incompatible to 2.20.6+incompatible #647 (dependabot-preview[bot])


## 0.68.1 (2020-07-20)

* Bump github.com/mackerelio/mackerel-client-go from 0.10.0 to 0.10.1 #646 (dependabot-preview[bot])
* Bump golang.org/x/text from 0.3.2 to 0.3.3 #645 (dependabot-preview[bot])
* Bump github.com/mackerelio/mackerel-client-go from 0.9.1 to 0.10.0 #643 (dependabot-preview[bot])


## 0.68.0 (2020-05-14)

* Bump github.com/shirou/gopsutil from 2.20.3+incompatible to 2.20.4+incompatible #641 (dependabot-preview[bot])
* Improve FreeBSD rc script #640 (metalefty)
* Bump github.com/shirou/gopsutil from 2.20.2+incompatible to 2.20.3+incompatible #639 (dependabot-preview[bot])
* [Windows]support x64 installation #637 (lufia)


## 0.67.1 (2020-04-03)

* Bump github.com/shirou/gopsutil from 2.20.1+incompatible to 2.20.2+incompatible #636 (dependabot-preview[bot])
* fix too late closing the response body #635 (shogo82148)
* Bump github.com/mackerelio/mackerel-client-go from 0.9.0 to 0.9.1 #634 (dependabot-preview[bot])


## 0.67.0 (2020-02-05)

* Bump github.com/shirou/gopsutil from 2.19.12+incompatible to 2.20.1+incompatible #631 (dependabot-preview[bot])
* Bump github.com/mackerelio/mackerel-client-go from 0.8.0 to 0.9.0 #629 (dependabot-preview[bot])
* Allow formatted duration in config #624 (itchyny)
* rename: github.com/motemen/gobump -> github.com/x-motemen/gobump #630 (lufia)
* Support IMDSv2 for AWS EC2 #627 (shogo82148)
* `%q` verb of fmt is invalid for map[string]float64 types #628 (shogo82148)


## 0.66.0 (2020-01-22)

* Bump github.com/pkg/errors from 0.8.1 to 0.9.1 #623 (dependabot-preview[bot])
* Bump github.com/shirou/gopsutil from 2.19.11+incompatible to 2.19.12+incompatible #620 (dependabot-preview[bot])
* Bump github.com/Songmu/prompter from 0.2.0 to 0.3.0 #617 (dependabot-preview[bot])
* Implement GCEGenerator.SuggestCustomIdentifier #618 (tanatana)
* fix how to get self executable path for autoshutdown option #616 (lufia)


## 0.65.0 (2019-12-05)

* add -private-autoshutdown option #612 (lufia)
* Fix Windows Edition name #614 (mattn)
* Bump github.com/shirou/gopsutil from 2.19.10+incompatible to 2.19.11+incompatible #611 (dependabot-preview[bot])
* update go-osstat and golang.org/x #610 (lufia)
* refactor: improve interface and testing for spec/cloud #609 (astj)
* refactor: Inject CloudMetaGenerators to Suggester in order to test them in safer way #608 (astj)


## 0.64.1 (2019-11-21)

* Install development tools in module-aware mode #606 (lufia)
* Bump github.com/shirou/gopsutil from 2.18.12+incompatible to 2.19.10+incompatible #602 (dependabot-preview[bot])
* Add armhf Debian package to release #599 (hnw)


## 0.64.0 (2019-10-24)

* Build with Go 1.12.12
* stop building 32bit Darwin artifacts #600 (astj)
* Fix wix/mackerel-agent.sample.conf #597 (ryosms)
* Pass the check monitoring result message to "action" by env #598 (a-know)
* Bump github.com/mackerelio/mackerel-client-go from 0.6.0 to 0.8.0 #595 (dependabot-preview[bot])
* add .dependabot/config.yml #594 (lufia)


## 0.63.0 (2019-09-11)

* avoid to use unnamed NICs for registering hosts on Windows #580 (lufia)
* Fixed to create configuration directory directory when executing init command if not exist directory #592 (homedm)


## 0.62.1 (2019-08-29)

* Update dependencies #590 (astj)


## 0.62.0 (2019-07-30)

* Allow working directory configuration in env of metadata plugins #585 (itchyny)
* Remove tempdir in tests #588 (astj)
* Remove memory.active and inactive metrics #584 (itchyny)
* Check command name on pid check for pid confliction after OS restart #583 (itchyny)
* change the owner of files created in docker #587 (hayajo)
* Fix to fit #557 into our workflow. #577 (hayajo)
* Add mips and arm64 architecture debian packaging support #557 (tnishinaga)


## 0.61.1 (2019-07-23)

* Set rpm dist to ".el7.centos", not ".el7" in rpm-v2 #581 (astj)


## 0.61.0 (2019-07-22)

* Generate and include CREDITS file in the release artifacts #575 (itchyny)
* Migrate docker repository #572 (hayajo)
* [check-plugin] Support custom_identifier  #571 (astj)
* Stop unnecessary builds #569 (lufia)
* Care newer busybox #570 (astj)
* migrate to mackerel.Client #566 (lufia)


## 0.60.0 (2019-06-11)

* migrate CreatingMetricsValue to mackerel.HostMetricValue #565 (lufia)
* migrate to use mkr.GraphDefsParam instead of CreateGraphDefsPayload #564 (lufia)
* migrate to use mkr.CheckReports instead of monitoringChecksPayload #563 (lufia)
* update appveyor.yml to build 64bit binaries #561 (lufia)
* migrate to use mkr.XxxHostParam instead of mackerel.HostSpec #554 (lufia)
* support Go Modules #549 (lufia)


## 0.59.3 (2019-05-08)

* Add rc script for FreeBSD #559 (owatan)
* migrate to use mkr.Interface instead of mackerel.NetInterface #553 (lufia)
* migrate to use mkr.Host instead of mackerel.Host #552 (lufia)


## 0.59.2 (2019-03-27)

* trim trailing newlines from command string on windows #548 (Songmu)
* Improve Makefile #547 (itchyny)


## 0.59.1 (2019-02-13)

* fix counter naming problem on Windows #544 (lufia)


## 0.59.0 (2019-01-10)

* Fix decoding error message of executables on Windows #539 (mattn)
* Fix detecting EC2 instance on Windows #540 (mattn)
* add check-disk plugin for Windows #541 (susisu)


## 0.58.2 (2018-11-27)

* [windows] Bump mkr to latest  #537 (astj)


## 0.58.1 (2018-11-26)

* Fix disk metrics for Linux kernel 4.19 #535 (itchyny)


## 0.58.0 (2018-11-12)

* To work in BusyBox #526 (Songmu)
* [incompatible] CollectDfValues only from local file systems on linux #532 (Songmu)


## 0.57.0 (2018-09-14)

* update Code Signing Certificate. #524 (hayajo)
* Build with Go 1.11 #522 (astj)
* [darwin] Fix iostat output parsing in CPU usage generator #520 (itchyny)
* [darwin] fix filesystem metrics for APFS vm partition volume #517 (itchyny)
* add loadavg1 and loadavg15 #519 (itchyny)


## 0.56.1 (2018-08-30)

* Do HTTP retry on determining cloud platform and suggesting customIdentifier #516 (astj)
* [windows] Add timeout to WMI query for disk metrics #511 (astj)


## 0.56.0 (2018-07-25)

* Fix starting order of Windows Service #506 (mattn)
* Auto retire with shutdown on Windows #505 (mattn)
* Use RunWithEnv instead of os.Setenv to avoid environment variable races #507 (itchyny)
* Improve debug messages for check monitoring actions #510 (itchyny)
* add mssql-plugin in windows msi #509 (daiksy)
* Replace GCE metadata endpoint with absolute FQDN #508 (i2tsuki)


## 0.55.0 (2018-06-20)

* improve PATH handling #501 (astj)
* Build with Go 1.10 #500 (astj)


## 0.54.1 (2018-03-28)

* Support UUID in little-endian format on EC2 detection #496 (hayajo)
* change the message level from WARNING to INFO when customIdentifier is not registered #493 (hayajo)


## 0.54.0 (2018-03-20)

* fix isEC2 #494 (Songmu)
* care `MemAvailable` in collecting metrics around memory on linux #491 (Songmu)


## 0.53.0 (2018-03-15)

* Stop collecting memory.available for now #490 (Songmu)
* omit `/Volumes/` from collected `df` values on darwin #489 (Songmu)
* Enhance diagnostic mode #486 (Songmu)
* Fix EC2 check for KVM based EC2 instance (e.g. c5 instance) #488 (hayajo)


## 0.52.1 (2018-03-01)

* context support in cmdutil #485 (Songmu)
* Improve error handling when executing commands #484 (Songmu)
* extend timeout for retrieving cloud metadata #483 (hayajo)


## 0.52.0 (2018-02-08)

* Refine metrics collector #442 (mechairoi)
*  Add `memo` option to check plugin config #480 (mechairoi)


## 0.51.0 (2018-01-23)

* Fix metric values of pagefile total and pagefile free on Windows #456 (itchyny)
* update rpm-v2 task for building Amazon Linux 2 package #475 (hayajo)
* Care plugins that handle timeout signal(SIGTERM) #476 (Songmu)


## 0.50.1 (2018-01-15)

* Add mkr to dependencies to include it into windows msi #478 (shibayu36)


## 0.50.0 (2018-01-15)

* use supervisor mode in sysvinit script for crash recovery #472 (Songmu)
* include mkr into windows msi #465 (Songmu)
* pass returned value from command.RunOnce so that `mackerel-agent once… #474 (astj)


## 0.49.0 (2018-01-10)

* cut out `cmdutil` package from `util` and interface adjustment #470 (Songmu)
* Ignore connection configurations in mackerel-agent.conf #463 (itchyny)
* fix error check in TestStart of start_test.go #471 (Ken2mer)
* [fix] `action` command in `checks` is able to have an individual timeout settings #469 (Songmu)
* Add an option of timeout duration for executing command #460 (taku-k)
* Adjust appveyor.yml #466 (Songmu)
* introduce goxz #468 (Songmu)
* using os.Executable() for getting executable path on windows environment #464 (Songmu)
* include commands_gen.go in repo for go-gettability #467 (Songmu)
* Ignore veth in network I/O metrics on Linux. (Docker creats a lot) #462 (hayajo)
* Ignore device-mapper in disk I/O metrics on Linux. (Docker creats a lot) #461 (hayajo)
* Ignore devicemapper #459 (hayajo)
* Ignore empty hostid file #458 (astj)
* add check-uptime.exe on msi #455 (Songmu)
* fix the retry of check reports #453 (hayajo)


## 0.48.2 (2017-12-20)

* Fix network interface spec collector on Windows #452 (itchyny)


## 0.48.1 (2017-12-13)

* fix a bug when action of check-plugin was not specified #450 (hayajo)


## 0.48.0 (2017-12-12)

* Set environment variables for plugins #448 (hayajo)
* Add an option to declare cloud platform explicitly #447 (astj)


## 0.47.3 (2017-11-28)

* Fix interface metrics of large counter values on Linux #445 (itchyny)
* Refine license notice #444 (itchyny)
* Improve plugin command parsing error message #443 (itchyny)
* Log stderr and err of check action #432 (mechairoi)
* Commonize interface generators for Linux, Darwin and add support for BSD systems #441 (itchyny)


## 0.47.2 (2017-11-09)

* Use go 1.9.2 #437 (astj)
* Commonize loadavg5 generators for Linux, Darwin and BSD systems #435 (itchyny)
* Change log level in device generator if /sys/block does not exist #424 (itchyny)


## 0.47.1 (2017-10-26)

* Use go-osstat library on linux #428 (itchyny)


## 0.47.0 (2017-10-19)

* Trigger action command after check plugin running. #425 (mechairoi)
* Ensure returned value of retrieveAzureVMMetadata is not null #429 (astj)
* Use go-osstat library on darwin #422 (itchyny)
* Subtract cpu.guest from cpu.user on Linux #423 (itchyny)
* Improve kernel spec generator performance for Linux #427 (itchyny)
* Improve implementation for memory spec on Linux #426 (itchyny)
* Do not send too many reports in one API request. #420 (astj)


## 0.46.0 (2017-10-04)

* Use new API BaseURL #417 (astj)
* Filter plugin metrics value by include_pattern and exclude_pattern option #416 (astj)


## 0.45.0 (2017-09-27)

* build with Go 1.9 #414 (astj)


## 0.44.2 (2017-08-30)

* Change the log level for failure of posting metric values #409 (itchyny)
* Show CPU/SoC model name on Linux/MIPS #408 (hnw)


## 0.44.1 (2017-08-23)

* Fail to start when custom identifiers are mismatched #405 (mechairoi)
* Fix the Azure VM check #404 (stefafafan)
* Adjust the Azure Virtual Machine metadata keys #403 (stefafafan)


## 0.44.0 (2017-07-26)

* Adjust isEC2 check  #401 (stefafafan)
* Support Azure VM Metadata #399 (stefafafan)
* FreeBSD: don't collect nullfs disk stat #400 (kyontan)
* Improve the EC2 Instance check #398 (stefafafan)


## 0.43.2 (2017-06-14)

* Revert "Enable HTTP/2" #393 (Songmu)
* [refactoring] remove version package and adjust internal dependencies #391 (Songmu)


## 0.43.1 (2017-05-17)

* rename command.Context to command.App #384 (Songmu)
* Add `prevent_alert_auto_close` option for check plugins #387 (mechairoi)
* Remove supported OS section from README. #388 (astj)


## 0.43.0 (2017-05-09)

* Use DiskReadsPerSec/DiskWritesPerSec instead of DiskReadBytesPersec/DiskWriteBytesPersec (on Windows) #382 (mattn)
* Enable HTTP/2 #383 (astj)


## 0.42.3 (2017-04-27)

* Output error logs of mackerel-agent as warning log of windows event log #380 (Songmu)


## 0.42.2 (2017-04-19)

* Adjust config package #375 (Songmu)
* use CRLF in mackerel-agent.conf on windows #377 (Songmu)


## 0.42.1 (2017-04-11)

* LC_ALL=C on initialization #373 (Songmu)


## 0.42.0 (2017-04-06)

* Logs that are not via the mackerel-agent's logger are also output to the eventlog #367 (Songmu)
* Change package License to Apache 2.0 #368 (astj)
* Release systemd deb packages to github releases #369 (astj)
* Change systemd deb package architecture to amd64 #370 (astj)


## 0.41.3 (2017-03-27)

* build with Go 1.8 #342 (astj)
* [EXPERIMENTAL] Add systemd support for deb packages #360 (astj)
* Timeout for command execution on Windows #361 (mattn)
* It need to read output from command continuously. #364 (mattn)
* remove util/util_windows.go and commonalize util.RunCommand #365 (Songmu)


## 0.41.2 (2017-03-22)

* Don't raise error when creating pidfile if the contents of pidfile is same as own pid #357 (Songmu)
* Exclude _tools from package #358 (itchyny)
* Add workaround for docker0 interface in docker-enabled Travis #359 (astj)


## 0.41.1 (2017-03-09)

* add check-tcp on pluginlist.txt #351 (daiksy)


## 0.41.0 (2017-03-07)

* [EXPERIMENTAL] systemd support for CentOS 7 #317 (astj)
* add `supervise` subcommand (supervisor mode) #327 (Songmu)
* Build RPM packages with Docker #330 (astj)
* run test with -race in CI #339 (haya14busa)
* Use hw.physmem64 instead of hw.physmem in NetBSD #343 (miwarin, astj)
* Build RPM files on CentOS5 on Docker #344 (astj)
* Keep environment variables when Agent runs commands with sudo #346 (astj)
* Release systemd RPMs to github releases #347 (astj)
* Fix disk metrics on Windows #348 (mattn)


## 0.40.0 (2017-02-22)

* support metadata plugins in configuration #331 (itchyny)
* Add metadata plugin feature #333 (itchyny)
* Use Named Result Parameters as document #334 (haya14busa)
* Set large number of file descriptors for the safety sake in init scripts #337 (Songmu)
* Improve darwin cpu spec #338 (astj)
* Fix format verb: use '%v' #340 (haya14busa)


## 0.39.4 (2017-02-08)

* prepare windows eventlog #319 (daiksy)
* Refactor plugin configurations #322 (itchyny)
* Execute less `go build`s on deploy #323 (astj)
* treat xmlns #324 (mattn)
* Fix xmlns #326 (mattn)


## 0.39.3 (2017-01-25)

* Fix segfault when loading a bad config file #316 (hanazuki)
* fix windows eventlog level when "verbose=true" #318 (daiksy)


## 0.39.2 (2017-01-16)

* Test wix/pluginlist.txt on AppVeyor ci #313 (astj)
* Revert "remove windows plugins on pluginslist" #314 (daiksy)


## 0.39.1 (2017-01-12)

* support filesystems.Ignore on windows #303 (Songmu)
* remove windows plugins on pluginslist #309 (daiksy)


## 0.39.0 (2017-01-11)

* implement `pluginGenerators` for windows #301 (daiksy)
* add check-windows-eventlog on pluginlist #302 (daiksy)
* Remove duplicated generator in Windows #305 (astj)
* add mackerel-plugin-windows-server-sessions on pluginlist #306 (daiksy)


## 0.38.0 (2016-12-21)

* fix typo #12 (ts-3156)
* Add Copyright #13 (yuuki)
* Separate interfaceGenerator from specGenerators #14 (motemen)
* Timout http reuquest in 30 sec (requries go 1.3) #17 (hakobe)
* specify command arguments in mackerel-agent.conf #293 (Songmu)
* several improvements for Windows #298 (daiksy)
* Avoid time.Tick and use time.NewTicker instead #299 (haya14busa)


## 0.37.1 (2016-11-29)

* fix pluginlist #291 (daiksy)
* Suppress ec2 metadata warnings #294 (itchyny)
* Uncapitalize error messages #295 (itchyny)


## 0.37.0 (2016-10-27)

* improve Windows support #289 (daiksy)


## 0.36.0 (2016-10-18)

* don't use HTTP_PROXY when requesting cloud instance metadata APIs #285 (Songmu)
* Add an option to output filesystem-related metrics with key by mountpoint #286 (astj)


## 0.35.1 (2016-09-29)

* support MACKEREL_PLUGIN_WORKDIR in init scripts #277 (Songmu)
* Add platform metadata for Darwin #280 (astj)
* Disable http2 for now #283 (Songmu)


## 0.35.0 (2016-09-07)

* built with Go 1.7 #266 (Songmu)
* remove `func (vs *Values) Merge(other Values)` #268 (Songmu)
* [incompatible] consider df  (used + available) as size of filesystem #271 (Songmu)
* Remove DigitalOcean related comment/definition from spec/cloud.go #272 (astj)
* Fix golint is not working on ci, and add some comment to pass golint #273 (astj)
* Add linux distribution information to kernel spec #274 (ak1t0)
* http_proxy configuration #275 (Songmu)
* set PATH and LANG only in unix environment #276 (Songmu)
* Ignore docker mapper storage in spec as well #278 (itchyny)


## 0.34.0 (2016-08-18)

* Reduce retry count on finding a host by the custom identifier #258 (itchyny)
* suppress checker flooding when resuming from sleep mode #260 (Songmu)
* truncate checker message up to 1024 characters #261 (Songmu)
* commonalize spec.FilesystemGenerator around unix OSs #262 (Songmu)
* define type DfStat,	remove dfColumnSpecs and refactor #263 (Songmu)


## 0.33.0 (2016-08-08)

* Fill the customIdentifier in EC2 #255 (itchyny)


## 0.32.2 (2016-07-14)

* fix GOMAXPROCS to 1 for avoiding rare panics #253 (Songmu)


## 0.32.1 (2016-07-07)

* Add user for executing a plugin #250 (y-kuno)


## 0.32.0 (2016-06-30)

* Added plugin check interval option #245 (karupanerura)


## 0.31.2 (2016-06-23)

* Refactor around metrics/linux/memory #242 (Songmu)
* Don't stop mackerel-agent process on upgrading by debian package #243 (karupanerura)
* add `silent` configuration key for suppressing log output #244 (Songmu)
* change log level ERROR to WARNING in spec/spec.go #246 (Songmu)
* remove /usr/local/bin from sample.conf #248 (Songmu)


## 0.31.0 (2016-05-25)

* Post the custom metrics to the hosts specified by custom identifiers #231 (itchyny)
* refactor FilesystemGenerator #233 (Songmu)
* Refactor metrics/linux/interface.go #234 (Songmu)
* remove regexp from spec/linux/cpu #235 (Songmu)
* Fix missing printf args #237 (shogo82148)


## 0.30.4 (2016-05-10)

* Recover from panic while processing generators #228 (stanaka)
* check length of cols just to be safe in metrics/linux/disk.go #229 (Songmu)


## 0.30.3 (2016-05-02)

* Remove usr local bin again #217 (Songmu)
* Fix typo #221 (yukiyan)
* Fix comments #222 (stefafafan)
* Remove go get cmd/vet #223 (itchyny)
* retry retirement when api request failed #224 (Songmu)
* output plugin stderr to log #226 (Songmu)


## 0.30.2 (2016-03-25)

* Revert "Merge pull request #211 from mackerelio/usr-bin" #215 (Songmu)


## 0.30.1 (2016-03-25)

* deprecate /usr/local/bin #211 (Songmu)
* use GOARCH=amd64 for now #213 (Songmu)


## 0.30.0 (2016-03-17)

* remove uptime metrics generator #161 (Songmu)
* Remove deprecated-sensu feature #202 (Songmu)
* Send all IP addresses of each interface (linux only) #205 (mechairoi)
* add `init` subcommand #207 (Songmu)
* Refactor net interface (multi ip support and bugfix) #208 (Songmu)
* Stop to fetch flags of cpu in spec/linux/cpu #209 (Songmu)


## 0.29.2 (2016-03-07)

* Don't overwrite mackerel-agent.conf when updating deb package (Fix deb packaging) #199 (Songmu)


## 0.29.1 (2016-03-04)

* maintenance release

## 0.29.0 (2016-03-02)

* remove deprecated command line options (-version and -once) #194 (Songmu)
* Report checker execution timeout as Unknown status #197 (hanazuki)


## 0.28.1 (2016-02-18)

* fix the exit status on stopping the agent in the init script of debian #192 (itchyny)


## 0.28.0 (2016-02-04)

* add a configuration to ignore filesystems #186 (stanaka)
* fix the code of extending the process's environment #187 (itchyny)
* s{code.google.com/p/winsvc}{golang.org/x/sys/windows/svc} #188 (Songmu)
* Max check attempts option for check plugin #189 (mechairoi)


## 0.27.1 (2016-01-08)

* [bugfix] fix timeout interval when calling `df` #184 (Songmu)


## 0.27.0 (2016-01-06)

* use timeout when calling `df` #180 (Songmu)
* Notification Interval for check monitoring #181 (itchyny)


## 0.26.2 (2015-12-10)

* output success message to stderr when configtest succeed #178 (Songmu)


## 0.26.1 (2015-12-09)

* fix deprecate message #176 (Songmu)


## 0.26.0 (2015-12-08)

* Make HostID storage replacable #167 (motemen)
* Publicize command.Context's fields #168 (motemen)
* Configtest #169 (fujiwara)
* Refactor config loading and check if Apikey exists in configtest #171 (Songmu)
* fix exit status of debian init script. #172 (fujiwara)
* Deprecate version and once option #173 (Songmu)


## 0.25.1 (2015-11-25)

* Go 1.5.1 #164 (Songmu)
* logging STDERR of checker command #165 (Songmu)


## 0.25.0 (2015-11-12)

* Retrieve interfaces on Darwin #158 (itchyny)
* add NetBSD support. #162 (miwarin)


## 0.24.1 (2015-11-05)

* We are Mackerel #156 (itchyny)


## 0.24.0 (2015-10-26)

* define config.agentName and set proper config path #150 (Songmu)
* /proc/cpuinfo parser for old ARM Linux kernels #152 (hanazuki)
* os.MkdirAll() before creating pidfile #153 (Songmu)


## 0.23.1 (2015-09-30)

* Code signing for windows installer #148 (mechairoi)


## 0.23.0 (2015-09-14)

* send check monitor report to server when check script failed even if the monitor result is not changed #143 (Songmu)
* Correct sample nginx comment. #144 (kamatama41)


## 0.22.0 (2015-09-02)

* add `reload` to init scripts #139 (Songmu)


## 0.21.0 (2015-09-02)

* Exclude mkr binary from deb/rpm package #137 (Sixeight)


## 0.20.1 (2015-08-13)

* use C struct for accessing Windows APIs #134 (stanaka)
* Fix bug that checks is not removed when no checks. #135 (Sixeight)


## 0.20.0 (2015-07-29)

* support subcommand #122 (Songmu)
* remove trailing newline chars when loading hostID #129 (Songmu)
* add sub-command `retire` and support $AUTO_RETIREMENT in initd #130 (Songmu)
* add postinst to register mackerel-agent to start-up (deb package) #131 (stanaka)
* bump bundled mkr version to 0.3.1 #132 (Songmu)


## 0.19.0 (2015-07-22)

* Support gce meta #115 (Songmu)
* Valid pidfile handling (fix on darwin) #123 (Songmu)
* -once only takes one second #126 (Songmu)
* fix shutdown priority in rpm/src/mackerel-agent.initd #127 (Songmu)


## 0.18.1 (2015-07-16)

* s/ami_id/ami-id/ in spec/cloud.go #112 (Songmu)
* remove `UpdateHost()` process from `prepareHost()` for simplicity #116 (Songmu)
* filter invalid roleFullNames with warning logs #117 (Songmu)
* allow using spaces as delimiter for custom metric values #119 (Songmu)


## 0.18.0 (2015-07-08)

* Retry in prepare #108 (Songmu)
* [WORKAROUND] downgrade golang version for windows #109 (Sixeight)


## 0.17.1 (2015-06-17)

* Update to go 1.4.2 for windows build #105 (mechairoi)


## 0.17.0 (2015-06-10)

* Set `displayName` via agent #92 (Sixeight)
* refactoring around api access #97 (Songmu)
* Configurable host status on start/stop agent #100 (Songmu)
* Add an agent memory usage metrics generator for diagnostic use #101 (hakobe)
* Add mkr to deb/rpm package #102 (Sixeight)


## 0.16.1 (2015-05-12)

* Code sharing around dfValues #85 (Songmu)
* [FreeBSD] Fix 'panic: runtime error: index out of range'. #89 (iwadon)
* separete out metrics/darwin/swap.go from memory.go #90 (Songmu)


## 0.16.0 (2015-05-08)

* suppress logging #78 (stanaka)
* "Check" functionality #80 (motemen)
* update for windows #81 (daiksy)
* collect memory metrics of osx #84 (Songmu)
* Send plugin.check._name_s list on `updateHost` #86 (mechairoi)


## 0.15.0 (2015-04-02)

* Only skip device mapper created by docker
* Run once and output results to stdout with -once option
* introduce Songmu/timeout for interrupting long time plugin execution
* add config.apibase
* output GOOS GOARCH runtime.Version() when -version option is specified

## 0.14.3 (2015-03-23)

- [enhancement] add collector for ec2 metadata

## 0.14.1 (2015-01-20)

* [fix] skip device mapper metrics
* [fix] filter invalid float values
* [enhancement] testing
* [enhancement] collect more metrics about darwin and freebsd

## 0.14.0 (2014-12-25)

* [improve] wait for termination until postQueue is empty up to 30 seconds.
* [improve] wait up to 30 seconds before initial posting
* [feature] work on Windows darwin FreeBSD (unofficial support)

## 0.8.0 (2014-06-26)

* [improve] Using go 1.3
* [feature] Periodically update host specs (#15)
* [fix] Http request now have timeout (#17)

## 0.7.0 (2014-06-06)

* [feature] Windows port (not officially supported), thanks to @mattn (#8)
* [fix] Replace invalid characters (e.g. '.') in disk and interface names with underscores (#10)
* [fix] Removed deprecated metrics (#11)
