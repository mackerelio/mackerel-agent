# sudo yum -y install rpmdevtools go && rpmdev-setuptree
# rpmbuild -ba ~/rpmbuild/SPECS/mackerel-agent.spec

%define _binaries_in_noarch_packages_terminate_build   0

Name:      mackerel-agent
Version:   %{_version}
Release:   1%{?dist}
License:   ASL 2.0
Summary:   mackerel.io agent
URL:       https://mackerel.io
Group:     Hatena Co., Ltd.
Source0:   %{name}.sysconfig
Source1:   %{name}.conf
Source2:   %{name}.service
Packager:  Hatena Co., Ltd.
BuildRoot: %{_tmppath}/%{name}-%{version}-%{release}-root

%{?systemd_requires}
BuildRequires: systemd

%description
mackerel.io agent

%prep

%build

%install
%{__rm} -rf %{buildroot}
%{__install} -Dp -m0755 %{_builddir}/%{name}             %{buildroot}%{_bindir}/%{name}
%{__install} -Dp -m0644 %{_sourcedir}/%{name}.sysconfig  %{buildroot}/%{_sysconfdir}/sysconfig/%{name}
%{__install} -Dp -m0644 %{_sourcedir}/%{name}.conf       %{buildroot}/%{_sysconfdir}/%{name}/%{name}.conf
%{__install} -Dp -m0644 %{_sourcedir}/%{name}.service    %{buildroot}%{_unitdir}/%{name}.service

%clean
%{__rm} -rf %{buildroot}

%post
%systemd_post %{name}.service
systemctl enable %{name}.service

%preun
%systemd_preun %{name}.service

%postun
%systemd_postun %{name}.service

%files
%defattr(-,root,root)
%{_unitdir}/%{name}.service
%{_bindir}/%{name}
%config(noreplace) %{_sysconfdir}/sysconfig/%{name}
%config(noreplace) %{_sysconfdir}/%{name}/%{name}.conf

%changelog
* Wed Dec 14 2022 <mackerel-developers@hatena.ne.jp> - 0.74.0
- refactor(config/validate): use `candidates` instead of `parentConfKeys` (by wafuwafu13)
- improve purge stage of Debian package to remove id file and keep custom files (by kmuto)
- fix(configtest): don't detect child key if parent key has already detected (by wafuwafu13)

* Wed Nov 9 2022 <mackerel-developers@hatena.ne.jp> - 0.73.3
- Fix config test (by ryuichi1208)

* Fri Nov 4 2022 <mackerel-developers@hatena.ne.jp> - 0.73.2
- Replace linter (by yseto)
- Bump github.com/mackerelio/mackerel-client-go from 0.21.2 to 0.22.0 (by dependabot[bot])
- Improve `mackerel-agent configtest`: Add suggestion to unexpected keys (by wafuwafu13)
- Bump github.com/Songmu/gocredits from 0.2.0 to 0.3.0 (by dependabot[bot])
- Bump github.com/mackerelio/mkr from 0.47.0 to 0.47.1 in /wix (by dependabot[bot])
- go.mod 1.17 -> 1.18 (by yseto)
- Improve `mackerel-agent configtest`: detect unexpected key (by wafuwafu13)
- fix deprecated function. (by yseto)
- Bump github.com/mackerelio/go-osstat from 0.2.2 to 0.2.3 (by dependabot[bot])
- Bump github.com/Songmu/goxz from 0.8.2 to 0.9.1 (by dependabot[bot])
- Bump github.com/BurntSushi/toml from 1.1.0 to 1.2.0 (by dependabot[bot])

* Wed Sep 14 2022 <mackerel-developers@hatena.ne.jp> - 0.73.1
- Bump github.com/mackerelio/mkr from 0.46.9 to 0.47.0 in /wix (by dependabot[bot])
- config_test: Add the case of LoadConfigWithInvalidToml (by wafuwafu13)
- In the test, if the Fatal if the result is nil. (by yseto)
- remove unused codes (by yseto)
- replace io/ioutil (by yseto)
- Bump github.com/mackerelio/mackerel-client-go from 0.21.1 to 0.21.2 (by dependabot[bot])
- get interface information via netlink on linux. (by yseto)

* Wed Jul 27 2022 <mackerel-developers@hatena.ne.jp> - 0.73.0
- Bump github.com/mackerelio/mkr from 0.46.8 to 0.46.9 in /wix (by dependabot[bot])
- Bump github.com/mackerelio/mackerel-agent-plugins from 0.71.0 to 0.72.1 in /wix (by dependabot[bot])
- Loosen the conditions of delaying report of check monitering. (by sugy)

* Wed Jul 20 2022 <mackerel-developers@hatena.ne.jp> - 0.72.15
- Bump github.com/mackerelio/mkr from 0.46.7 to 0.46.8 in /wix (by dependabot[bot])
- Bump github.com/mackerelio/go-check-plugins from 0.42.0 to 0.42.1 in /wix (by dependabot[bot])
- Bump github.com/mackerelio/mackerel-client-go from 0.21.0 to 0.21.1 (by dependabot[bot])

* Wed Jun 22 2022 <mackerel-developers@hatena.ne.jp> - 0.72.14
- add s, ms, bps to metric units (by Arthur1)
- Bump github.com/mackerelio/mkr from 0.46.6 to 0.46.7 in /wix (by dependabot[bot])

* Wed Jun 8 2022 <mackerel-developers@hatena.ne.jp> - 0.72.13
- Bump github.com/mackerelio/mkr from 0.46.5 to 0.46.6 in /wix (by dependabot[bot])
- Bump github.com/Songmu/prompter from 0.5.0 to 0.5.1 (by dependabot[bot])

* Thu May 26 2022 <mackerel-developers@hatena.ne.jp> - 0.72.12
- Bump github.com/Songmu/goxz from 0.8.1 to 0.8.2 (by dependabot[bot])
- Bump github.com/mackerelio/mackerel-agent-plugins from 0.70.6 to 0.71.0 in /wix (by dependabot[bot])

* Thu Apr 14 2022 <mackerel-developers@hatena.ne.jp> - 0.72.11
- Bump github.com/mackerelio/mackerel-agent-plugins from 0.70.4 to 0.70.6 in /wix (by dependabot[bot])
- Bump github.com/mackerelio/go-check-plugins from 0.41.7 to 0.42.0 in /wix (by dependabot[bot])
- Bump github.com/mackerelio/mkr from 0.46.4 to 0.46.5 in /wix (by dependabot[bot])
- Bump github.com/BurntSushi/toml from 1.0.0 to 1.1.0 (by dependabot[bot])

* Wed Mar 30 2022 <mackerel-developers@hatena.ne.jp> - 0.72.10
- Bump github.com/mackerelio/go-osstat from 0.2.1 to 0.2.2 (by dependabot[bot])
- Bump github.com/mackerelio/go-check-plugins from 0.41.6 to 0.41.7 in /wix (by dependabot[bot])
- Bump github.com/mackerelio/mackerel-agent-plugins from 0.70.3 to 0.70.4 in /wix (by dependabot[bot])
- Bump github.com/mackerelio/mkr from 0.46.3 to 0.46.4 in /wix (by dependabot[bot])

* Tue Mar 15 2022 <mackerel-developers@hatena.ne.jp> - 0.72.9
- Bump github.com/mackerelio/go-check-plugins from 0.41.5 to 0.41.6 in /wix (by dependabot[bot])
- Bump github.com/mackerelio/mkr from 0.46.2 to 0.46.3 in /wix (by dependabot[bot])
- Bump github.com/mackerelio/mackerel-agent-plugins from 0.70.2 to 0.70.3 in /wix (by dependabot[bot])

* Wed Feb 16 2022 <mackerel-developers@hatena.ne.jp> - 0.72.8
- upgrade Go 1.16 -> 1.17 (by lufia)
- Bump github.com/mackerelio/mkr from 0.46.1 to 0.46.2 in /wix (by dependabot[bot])
- Bump github.com/mackerelio/go-check-plugins from 0.41.4 to 0.41.5 in /wix (by dependabot[bot])

* Wed Feb 2 2022 <mackerel-developers@hatena.ne.jp> - 0.72.7
- Bump github.com/mackerelio/mkr from 0.46.0 to 0.46.1 in /wix (by dependabot[bot])
- Bump github.com/mackerelio/mackerel-agent-plugins from 0.69.1 to 0.70.2 in /wix (by dependabot[bot])
- Bump github.com/mackerelio/go-check-plugins from 0.41.1 to 0.41.4 in /wix (by dependabot[bot])
- Bump github.com/BurntSushi/toml from 0.3.1 to 1.0.0 (by dependabot[bot])
- Bump github.com/Songmu/goxz from 0.7.0 to 0.8.1 (by dependabot[bot])
- Bump github.com/mackerelio/mackerel-client-go from 0.19.0 to 0.21.0 (by dependabot[bot])

* Wed Jan 12 2022 <mackerel-developers@hatena.ne.jp> - 0.72.6
- Bump github.com/mackerelio/mkr from v0.45.3 to 0.46.0 (by susisu)

* Wed Dec 1 2021 <mackerel-developers@hatena.ne.jp> - 0.72.5
- Bump github.com/mackerelio/mackerel-client-go from 0.17.0 to 0.19.0 (by dependabot[bot])

* Thu Nov 18 2021 <mackerel-developers@hatena.ne.jp> - 0.72.4
- spec: reuse http.Client (by lufia)
- Bump github.com/mackerelio/go-osstat from 0.2.0 to 0.2.1 (by dependabot[bot])
- Bump github.com/mattn/goveralls from 0.0.9 to 0.0.11 (by dependabot[bot])
- make wix/ a submodule (by susisu)
- Add arm64/darwin build to GitHub release (by astj)
- read a response body even if status is not good (by lufia)

* Wed Oct 20 2021 <mackerel-developers@hatena.ne.jp> - 0.72.3
- Bump github.com/mackerelio/mackerel-agent-plugins from v0.65.0 to v0.69.1 (by susisu)
- Bump github.com/mackerelio/go-check-plugins from v0.39.5 to v0.41.1 (by ne-sachirou)

* Mon Sep 6 2021 <mackerel-developers@hatena.ne.jp> - 0.72.2
- Update Code Signing Certificates (by Krout0n)
- Bump github.com/mackerelio/go-check-plugins from 0.39.3 to 0.39.5 (by dependabot[bot])
- Bump github.com/mackerelio/mackerel-agent-plugins from 0.64.2 to 0.65.0 (by dependabot[bot])
- Bump github.com/mackerelio/mkr from 0.45.1 to 0.45.2 (by dependabot[bot])

* Wed Jun 23 2021 <mackerel-developers@hatena.ne.jp> - 0.72.1
- Bump github.com/mackerelio/go-check-plugins from 0.39.2 to 0.39.3 (by dependabot[bot])
- Bump github.com/mattn/goveralls from 0.0.8 to 0.0.9 (by dependabot[bot])
- Bump github.com/mackerelio/mackerel-agent-plugins from 0.64.0 to 0.64.2 (by dependabot[bot])
- Bump github.com/mackerelio/mkr from 0.45.0 to 0.45.1 (by dependabot[bot])

* Thu Jun 17 2021 <mackerel-developers@hatena.ne.jp> - 0.72.0
- fix http_proxy option in v0.71.2 (by yseto)

* Wed May 26 2021 <mackerel-developers@hatena.ne.jp> - 0.71.2
- Bump github.com/mackerelio/go-osstat from 0.1.0 to 0.2.0 (by dependabot[bot])
- Bump golang.org/x/text from 0.3.5 to 0.3.6 (by dependabot[bot])
- Bump github.com/Songmu/goxz from 0.6.0 to 0.7.0 (by dependabot[bot])
- upgrade Go 1.14 to 1.16 (by lufia)
- [ci] avoid additional Go installation (by lufia)
- Bump github.com/Songmu/prompter from 0.4.0 to 0.5.0 (by dependabot[bot])
- Bump github.com/mackerelio/mackerel-client-go from 0.15.0 to 0.16.0 (by dependabot[bot])
- Bump github.com/mattn/goveralls from 0.0.7 to 0.0.8 (by dependabot[bot])
- Bump github.com/mackerelio/mackerel-client-go from 0.13.0 to 0.15.0 (by dependabot[bot])
- Bump github.com/mackerelio/golib from 1.1.0 to 1.2.0 (by dependabot[bot])
- [ci] fix option at repository-dispatch (by yseto)
- [ci] added repository_dispatch to homebrew-mackerel-agent (by yseto)
- [ci] replace token (by yseto)
- [ci] replace mackerel-github-release (by yseto)
- Changed the test method from TestDiskGenerator to TestParseDiskStats  because the test results are flaky. (by yseto)

* Thu Jan 21 2021 <mackerel-developers@hatena.ne.jp> - 0.71.1
- remove .circleci/config.yml (by yseto)
- [ci] bump Windows i386 Golang to 1.14.14 (by astj)
- Build Windows package on GitHub Actions (by yseto)
- Bump github.com/mackerelio/golib from 1.0.0 to 1.1.0 (by dependabot[bot])
- Bump github.com/mackerelio/mackerel-client-go from 0.12.0 to 0.13.0 (by dependabot[bot])
- Bump golang.org/x/text from 0.3.4 to 0.3.5 (by dependabot[bot])

* Mon Dec 14 2020 <mackerel-developers@hatena.ne.jp> - 0.71.0
- Bump github.com/mackerelio/mackerel-client-go from 0.11.0 to 0.12.0 (by dependabot[bot])
- Network interface exclusion feature (by yseto)
- Bump golang.org/x/text from 0.3.3 to 0.3.4 (by dependabot-preview[bot])

* Wed Nov 25 2020 <mackerel-developers@hatena.ne.jp> - 0.70.3
- Build with Go 1.14 in CI (was 1.15 by mistake) (by astj)

* Thu Nov 19 2020 <mackerel-developers@hatena.ne.jp> - 0.70.2
- Fix artifact filename pattern again to include mackerel-agent_{os}_{arch}.tar.gz to GitHub Release artifacts (by astj)

* Thu Nov 19 2020 <mackerel-developers@hatena.ne.jp> - 0.70.1
- include mackerel-agent_{os}_{arch}.tar.gz to GitHub Release artifacts (by astj)

* Thu Nov 19 2020 <mackerel-developers@hatena.ne.jp> - 0.70.0
- replace Travis CI workflow with GitHub Actions (by astj)
- Retry once immediately on posting metrics and check reports when error is caused by net/http (network error) (by astj)

* Wed Oct 28 2020 <mackerel-developers@hatena.ne.jp> - 0.69.3
- Bump github.com/shirou/gopsutil from 2.20.8+incompatible to 2.20.9+incompatible (by dependabot-preview[bot])

* Thu Oct 01 2020 <mackerel-developers@hatena.ne.jp> - 0.69.2
- Bump github.com/mackerelio/mackerel-client-go from 0.10.1 to 0.11.0 (by dependabot-preview[bot])

* Tue Sep 15 2020 <mackerel-developers@hatena.ne.jp> - 0.69.1
- kcps, stage: add --target option to build rpm (by lufia)

* Tue Sep 15 2020 <mackerel-developers@hatena.ne.jp> - 0.69.0
- revert changing filename unexpectedly (by lufia)
- Bump github.com/shirou/gopsutil from 2.20.6+incompatible to 2.20.8+incompatible (by dependabot-preview[bot])
- Bump github.com/mattn/goveralls from 0.0.6 to 0.0.7 (by dependabot-preview[bot])
- Bump github.com/Songmu/prompter from 0.3.0 to 0.4.0 (by dependabot-preview[bot])
- revert mkdir with shell expansion (by lufia)
- add arm64 RPM packages, and change deb architecture to be correct (by lufia)
- update go: 1.12 -> 1.14 (by lufia)

* Wed Jul 29 2020 <mackerel-developers@hatena.ne.jp> - 0.68.2
- Bump github.com/shirou/gopsutil from 2.20.4+incompatible to 2.20.6+incompatible (by dependabot-preview[bot])

* Mon Jul 20 2020 <mackerel-developers@hatena.ne.jp> - 0.68.1
- Bump github.com/mackerelio/mackerel-client-go from 0.10.0 to 0.10.1 (by dependabot-preview[bot])
- Bump golang.org/x/text from 0.3.2 to 0.3.3 (by dependabot-preview[bot])
- Bump github.com/mackerelio/mackerel-client-go from 0.9.1 to 0.10.0 (by dependabot-preview[bot])

* Thu May 14 2020 <mackerel-developers@hatena.ne.jp> - 0.68.0
- Bump github.com/shirou/gopsutil from 2.20.3+incompatible to 2.20.4+incompatible (by dependabot-preview[bot])
- Improve FreeBSD rc script (by metalefty)
- Bump github.com/shirou/gopsutil from 2.20.2+incompatible to 2.20.3+incompatible (by dependabot-preview[bot])
- [Windows]support x64 installation (by lufia)

* Fri Apr 03 2020 <mackerel-developers@hatena.ne.jp> - 0.67.1
- Bump github.com/shirou/gopsutil from 2.20.1+incompatible to 2.20.2+incompatible (by dependabot-preview[bot])
- fix too late closing the response body (by shogo82148)
- Bump github.com/mackerelio/mackerel-client-go from 0.9.0 to 0.9.1 (by dependabot-preview[bot])

* Wed Feb 05 2020 <mackerel-developers@hatena.ne.jp> - 0.67.0
- Bump github.com/shirou/gopsutil from 2.19.12+incompatible to 2.20.1+incompatible (by dependabot-preview[bot])
- Bump github.com/mackerelio/mackerel-client-go from 0.8.0 to 0.9.0 (by dependabot-preview[bot])
- Allow formatted duration in config (by itchyny)
- rename: github.com/motemen/gobump -> github.com/x-motemen/gobump (by lufia)
- Support IMDSv2 for AWS EC2 (by shogo82148)
- `%q` verb of fmt is invalid for map[string]float64 types (by shogo82148)

* Wed Jan 22 2020 <mackerel-developers@hatena.ne.jp> - 0.66.0
- Bump github.com/pkg/errors from 0.8.1 to 0.9.1 (by dependabot-preview[bot])
- Bump github.com/shirou/gopsutil from 2.19.11+incompatible to 2.19.12+incompatible (by dependabot-preview[bot])
- Bump github.com/Songmu/prompter from 0.2.0 to 0.3.0 (by dependabot-preview[bot])
- Implement GCEGenerator.SuggestCustomIdentifier (by tanatana)
- fix how to get self executable path for autoshutdown option (by lufia)

* Thu Dec 05 2019 <mackerel-developers@hatena.ne.jp> - 0.65.0
- add -private-autoshutdown option (by lufia)
- Fix Windows Edition name (by mattn)
- Bump github.com/shirou/gopsutil from 2.19.10+incompatible to 2.19.11+incompatible (by dependabot-preview[bot])
- update go-osstat and golang.org/x (by lufia)
- refactor: improve interface and testing for spec/cloud (by astj)
- refactor: Inject CloudMetaGenerators to Suggester in order to test them in safer way (by astj)

* Thu Nov 21 2019 <mackerel-developers@hatena.ne.jp> - 0.64.1
- Install development tools in module-aware mode (by lufia)
- Bump github.com/shirou/gopsutil from 2.18.12+incompatible to 2.19.10+incompatible (by dependabot-preview[bot])
- Add armhf Debian package to release (by hnw)

* Thu Oct 24 2019 <mackerel-developers@hatena.ne.jp> - 0.64.0
- Build with Go 1.12.12
- stop building 32bit Darwin artifacts (by astj)
- Fix wix/mackerel-agent.sample.conf (by ryosms)
- Pass the check monitoring result message to "action" by env (by a-know)
- Bump github.com/mackerelio/mackerel-client-go from 0.6.0 to 0.8.0 (by dependabot-preview[bot])
- add .dependabot/config.yml (by lufia)

* Wed Sep 11 2019 <mackerel-developers@hatena.ne.jp> - 0.63.0
- avoid to use unnamed NICs for registering hosts on Windows (by lufia)
- Fixed to create configuration directory directory when executing init command if not exist directory (by homedm)

* Thu Aug 29 2019 <mackerel-developers@hatena.ne.jp> - 0.62.1
- Update dependencies (by astj)

* Tue Jul 30 2019 <mackerel-developers@hatena.ne.jp> - 0.62.0
- Allow working directory configuration in env of metadata plugins (by itchyny)
- Remove tempdir in tests (by astj)
- Remove memory.active and inactive metrics (by itchyny)
- Check command name on pid check for pid confliction after OS restart (by itchyny)
- change the owner of files created in docker (by hayajo)
- Fix to fit #557 into our workflow. (by hayajo)

* Tue Jul 23 2019 <mackerel-developers@hatena.ne.jp> - 0.61.1
- Set rpm dist to ".el7.centos", not ".el7" in rpm-v2 (by astj)

* Mon Jul 22 2019 <mackerel-developers@hatena.ne.jp> - 0.61.0
- Generate and include CREDITS file in the release artifacts (by itchyny)
- Migrate docker repository (by hayajo)
- [check-plugin] Support custom_identifier  (by astj)
- Stop unnecessary builds (by lufia)
- Care newer busybox (by astj)
- migrate to mackerel.Client (by lufia)

* Tue Jun 11 2019 <mackerel-developers@hatena.ne.jp> - 0.60.0
- migrate CreatingMetricsValue to mackerel.HostMetricValue (by lufia)
- migrate to use mkr.GraphDefsParam instead of CreateGraphDefsPayload (by lufia)
- migrate to use mkr.CheckReports instead of monitoringChecksPayload (by lufia)
- update appveyor.yml to build 64bit binaries (by lufia)
- migrate to use mkr.XxxHostParam instead of mackerel.HostSpec (by lufia)
- support Go Modules (by lufia)

* Wed May 08 2019 <mackerel-developers@hatena.ne.jp> - 0.59.3
- Add rc script for FreeBSD (by owatan)
- migrate to use mkr.Interface instead of mackerel.NetInterface (by lufia)
- migrate to use mkr.Host instead of mackerel.Host (by lufia)

* Wed Mar 27 2019 <mackerel-developers@hatena.ne.jp> - 0.59.2
- trim trailing newlines from command string on windows (by Songmu)
- Improve Makefile (by itchyny)

* Wed Feb 13 2019 <mackerel-developers@hatena.ne.jp> - 0.59.1
- fix counter naming problem on Windows (by lufia)

* Thu Jan 10 2019 <mackerel-developers@hatena.ne.jp> - 0.59.0
- Fix decoding error message of executables on Windows (by mattn)
- Fix detecting EC2 instance on Windows (by mattn)
- add check-disk plugin for Windows (by susisu)

* Tue Nov 27 2018 <mackerel-developers@hatena.ne.jp> - 0.58.2
- [windows] Bump mkr to latest  (by astj)

* Mon Nov 26 2018 <mackerel-developers@hatena.ne.jp> - 0.58.1
- Fix disk metrics for Linux kernel 4.19 (by itchyny)

* Mon Nov 12 2018 <mackerel-developers@hatena.ne.jp> - 0.58.0
- To work in BusyBox (by Songmu)
- [incompatible] CollectDfValues only from local file systems on linux (by Songmu)

* Fri Sep 14 2018 <mackerel-developers@hatena.ne.jp> - 0.57.0
- update Code Signing Certificate. (by hayajo)
- Build with Go 1.11 (by astj)
- [darwin] Fix iostat output parsing in CPU usage generator (by itchyny)
- [darwin] fix filesystem metrics for APFS vm partition volume (by itchyny)
- add loadavg1 and loadavg15 (by itchyny)

* Thu Aug 30 2018 <mackerel-developers@hatena.ne.jp> - 0.56.1
- Do HTTP retry on determining cloud platform and suggesting customIdentifier (by astj)
- [windows] Add timeout to WMI query for disk metrics (by astj)

* Wed Jul 25 2018 <mackerel-developers@hatena.ne.jp> - 0.56.0
- Fix starting order of Windows Service (by mattn)
- Auto retire with shutdown on Windows (by mattn)
- Use RunWithEnv instead of os.Setenv to avoid environment variable races (by itchyny)
- Improve debug messages for check monitoring actions (by itchyny)
- add mssql-plugin in windows msi (by daiksy)
- Replace GCE metadata endpoint with absolute FQDN (by i2tsuki)

* Wed Jun 20 2018 <mackerel-developers@hatena.ne.jp> - 0.55.0
- improve PATH handling (by astj)
- Build with Go 1.10 (by astj)

* Wed Mar 28 2018 <mackerel-developers@hatena.ne.jp> - 0.54.1
- Support UUID in little-endian format on EC2 detection (by hayajo)
- change the message level from WARNING to INFO when customIdentifier is not registered (by hayajo)

* Tue Mar 20 2018 <mackerel-developers@hatena.ne.jp> - 0.54.0
- fix isEC2 (by Songmu)
- care `MemAvailable` in collecting metrics around memory on linux (by Songmu)

* Thu Mar 15 2018 <mackerel-developers@hatena.ne.jp> - 0.53.0
- Stop collecting memory.available for now (by Songmu)
- omit `/Volumes/` from collected `df` values on darwin (by Songmu)
- Enhance diagnostic mode (by Songmu)
- Fix EC2 check for KVM based EC2 instance (e.g. c5 instance) (by hayajo)

* Thu Mar 01 2018 <mackerel-developers@hatena.ne.jp> - 0.52.1
- context support in cmdutil (by Songmu)
- Improve error handling when executing commands (by Songmu)
- extend timeout for retrieving cloud metadata (by hayajo)

* Thu Feb 08 2018 <mackerel-developers@hatena.ne.jp> - 0.52.0
- Refine metrics collector (by mechairoi)
-  Add `memo` option to check plugin config (by mechairoi)

* Tue Jan 23 2018 <mackerel-developers@hatena.ne.jp> - 0.51.0
- Fix metric values of pagefile total and pagefile free on Windows (by itchyny)
- update rpm-v2 task for building Amazon Linux 2 package (by hayajo)
- Care plugins that handle timeout signal(SIGTERM) (by Songmu)

* Mon Jan 15 2018 <mackerel-developers@hatena.ne.jp> - 0.50.1
- Add mkr to dependencies to include it into windows msi (by shibayu36)

* Mon Jan 15 2018 <mackerel-developers@hatena.ne.jp> - 0.50.0
- use supervisor mode in sysvinit script for crash recovery (by Songmu)
- include mkr into windows msi (by Songmu)
- pass returned value from command.RunOnce so that `mackerel-agent once… (by astj)

* Wed Jan 10 2018 <mackerel-developers@hatena.ne.jp> - 0.49.0
- cut out `cmdutil` package from `util` and interface adjustment (by Songmu)
- Ignore connection configurations in mackerel-agent.conf (by itchyny)
- fix error check in TestStart of start_test.go (by Ken2mer)
- [fix] `action` command in `checks` is able to have an individual timeout settings (by Songmu)
- Add an option of timeout duration for executing command (by taku-k)
- Adjust appveyor.yml (by Songmu)
- introduce goxz (by Songmu)
- using os.Executable() for getting executable path on windows environment (by Songmu)
- include commands_gen.go in repo for go-gettability (by Songmu)
- Ignore veth in network I/O metrics on Linux. (Docker creats a lot) (by hayajo)
- Ignore device-mapper in disk I/O metrics on Linux. (Docker creats a lot) (by hayajo)
- Ignore devicemapper (by hayajo)
- Ignore empty hostid file (by astj)
- add check-uptime.exe on msi (by Songmu)
- fix the retry of check reports (by hayajo)

* Wed Dec 20 2017 <mackerel-developers@hatena.ne.jp> - 0.48.2
- Fix network interface spec collector on Windows (by itchyny)

* Wed Dec 13 2017 <mackerel-developers@hatena.ne.jp> - 0.48.1
- fix a bug when action of check-plugin was not specified (by hayajo)

* Tue Dec 12 2017 <mackerel-developers@hatena.ne.jp> - 0.48.0
- Set environment variables for plugins (by hayajo)
- Add an option to declare cloud platform explicitly (by astj)

* Tue Nov 28 2017 <mackerel-developers@hatena.ne.jp> - 0.47.3
- Fix interface metrics of large counter values on Linux (by itchyny)
- Refine license notice (by itchyny)
- Improve plugin command parsing error message (by itchyny)
- Log stderr and err of check action (by mechairoi)
- Commonize interface generators for Linux, Darwin and add support for BSD systems (by itchyny)

* Thu Nov 09 2017 <mackerel-developers@hatena.ne.jp> - 0.47.2
- Use go 1.9.2 (by astj)
- Commonize loadavg5 generators for Linux, Darwin and BSD systems (by itchyny)
- Change log level in device generator if /sys/block does not exist (by itchyny)

* Thu Oct 26 2017 <mackerel-developers@hatena.ne.jp> - 0.47.1
- Use go-osstat library on linux (by itchyny)

* Thu Oct 19 2017 <mackerel-developers@hatena.ne.jp> - 0.47.0
- Trigger action command after check plugin running. (by mechairoi)
- Ensure returned value of retrieveAzureVMMetadata is not null (by astj)
- Use go-osstat library on darwin (by itchyny)
- Subtract cpu.guest from cpu.user on Linux (by itchyny)
- Improve kernel spec generator performance for Linux (by itchyny)
- Improve implementation for memory spec on Linux (by itchyny)
- Do not send too many reports in one API request. (by astj)

* Wed Oct 04 2017 <mackerel-developers@hatena.ne.jp> - 0.46.0
- Use new API BaseURL (by astj)
- Filter plugin metrics value by include_pattern and exclude_pattern option (by astj)

* Wed Sep 27 2017 <mackerel-developers@hatena.ne.jp> - 0.45.0
- build with Go 1.9 (by astj)

* Wed Aug 30 2017 <mackerel-developers@hatena.ne.jp> - 0.44.2
- Change the log level for failure of posting metric values (by itchyny)
- Show CPU/SoC model name on Linux/MIPS (by hnw)

* Wed Aug 23 2017 <mackerel-developers@hatena.ne.jp> - 0.44.1
- Fail to start when custom identifiers are mismatched (by mechairoi)
- Fix the Azure VM check (by stefafafan)
- Adjust the Azure Virtual Machine metadata keys (by stefafafan)

* Wed Jul 26 2017 <mackerel-developers@hatena.ne.jp> - 0.44.0
- Adjust isEC2 check  (by stefafafan)
- Support Azure VM Metadata (by stefafafan)
- FreeBSD: don't collect nullfs disk stat (by kyontan)
- Improve the EC2 Instance check (by stefafafan)

* Wed Jun 14 2017 <mackerel-developers@hatena.ne.jp> - 0.43.2
- Revert "Enable HTTP/2" (by Songmu)
- [refactoring] remove version package and adjust internal dependencies (by Songmu)

* Wed May 17 2017 <mackerel-developers@hatena.ne.jp> - 0.43.1-1
- rename command.Context to command.App (by Songmu)
- Add `prevent_alert_auto_close` option for check plugins (by mechairoi)
- Remove supported OS section from README. (by astj)

* Tue May 09 2017 <mackerel-developers@hatena.ne.jp> - 0.43.0-1
- Use DiskReadsPerSec/DiskWritesPerSec instead of DiskReadBytesPersec/DiskWriteBytesPersec (on Windows) (by mattn)
- Enable HTTP/2 (by astj)

* Thu Apr 27 2017 <mackerel-developers@hatena.ne.jp> - 0.42.3-1
- Output error logs of mackerel-agent as warning log of windows event log (by Songmu)

* Wed Apr 19 2017 <mackerel-developers@hatena.ne.jp> - 0.42.2-1
- Adjust config package (by Songmu)
- use CRLF in mackerel-agent.conf on windows (by Songmu)

* Tue Apr 11 2017 <mackerel-developers@hatena.ne.jp> - 0.42.1-1
- LC_ALL=C on initialization (by Songmu)

* Thu Apr 06 2017 <mackerel-developers@hatena.ne.jp> - 0.42.0-1
- Logs that are not via the mackerel-agent's logger are also output to the eventlog (by Songmu)
- Change package License to Apache 2.0 (by astj)
- Release systemd deb packages to github releases (by astj)
- Change systemd deb package architecture to amd64 (by astj)

* Mon Mar 27 2017 <mackerel-developers@hatena.ne.jp> - 0.41.3-1
- build with Go 1.8 (by astj)
- [EXPERIMENTAL] Add systemd support for deb packages (by astj)
- Timeout for command execution on Windows (by mattn)
- It need to read output from command continuously. (by mattn)
- remove util/util_windows.go and commonalize util.RunCommand (by Songmu)

* Wed Mar 22 2017 <mackerel-developers@hatena.ne.jp> - 0.41.2-1
- Don't raise error when creating pidfile if the contents of pidfile is same as own pid (by Songmu)
- Exclude _tools from package (by itchyny)
- Add workaround for docker0 interface in docker-enabled Travis (by astj)

* Thu Mar 09 2017 <mackerel-developers@hatena.ne.jp> - 0.41.1-1
- add check-tcp on pluginlist.txt (by daiksy)

* Tue Mar 07 2017 <mackerel-developers@hatena.ne.jp> - 0.41.0-1
- [EXPERIMENTAL] systemd support for CentOS 7 (by astj)
- add `supervise` subcommand (supervisor mode) (by Songmu)
- Build RPM packages with Docker (by astj)
- run test with -race in CI (by haya14busa)
- Use hw.physmem64 instead of hw.physmem in NetBSD (by miwarin, astj)
- Build RPM files on CentOS5 on Docker (by astj)
- Keep environment variables when Agent runs commands with sudo (by astj)
- Release systemd RPMs to github releases (by astj)
- Fix disk metrics on Windows (by mattn)

