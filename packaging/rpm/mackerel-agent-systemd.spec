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
BuildArch: %{buildarch}
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

