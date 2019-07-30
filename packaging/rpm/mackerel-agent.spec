# sudo yum -y install rpmdevtools go && rpmdev-setuptree
# rpmbuild -ba ~/rpmbuild/SPECS/mackerel-agent.spec

%define _binaries_in_noarch_packages_terminate_build   0

Name:      mackerel-agent
Version:   %{_version}
Release:   1
License:   ASL 2.0
Summary:   mackerel.io agent
URL:       https://mackerel.io
Group:     Hatena Co., Ltd.
Source0:   %{name}.initd
Source1:   %{name}.sysconfig
Source2:   %{name}.logrotate
Source3:   %{name}.conf
Packager:  Hatena Co., Ltd.
BuildArch: %{buildarch}
BuildRoot: %{_tmppath}/%{name}-%{version}-%{release}-root
Requires(post): /sbin/chkconfig
Requires(preun): /sbin/chkconfig, /sbin/service
Requires(postun): /sbin/service

%description
mackerel.io agent

%prep

%build

%install
%{__rm} -rf %{buildroot}
%{__install} -Dp -m0755 %{_builddir}/%{name}             %{buildroot}%{_bindir}/%{name}
%{__install} -d  -m0755                                  %{buildroot}/%{_localstatedir}/log/
%{__install} -Dp -m0755 %{_sourcedir}/%{name}.initd      %{buildroot}/%{_initrddir}/%{name}
%{__install} -Dp -m0644 %{_sourcedir}/%{name}.sysconfig  %{buildroot}/%{_sysconfdir}/sysconfig/%{name}
%{__install} -Dp -m0644 %{_sourcedir}/%{name}.logrotate  %{buildroot}/%{_sysconfdir}/logrotate.d/%{name}
%{__install} -Dp -m0644 %{_sourcedir}/%{name}.conf       %{buildroot}/%{_sysconfdir}/%{name}/%{name}.conf
%{__install} -Dp -m0755 %{_sourcedir}/%{name}.deprecated %{buildroot}/usr/local/bin/%{name}

%clean
%{__rm} -rf %{buildroot}

%pre

%post
/sbin/chkconfig --add %{name}

%preun
if [ $1 = 0 ]; then
  service %{name} stop > /dev/null 2>&1
  chkconfig --del %{name}
fi

%files
%defattr(-,root,root)
%{_initrddir}/%{name}
%{_bindir}/%{name}
%config(noreplace) %{_sysconfdir}/sysconfig/%{name}
%config(noreplace) %{_sysconfdir}/%{name}/%{name}.conf
%{_sysconfdir}/logrotate.d/%{name}
/usr/local/bin/%{name}

%changelog
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

* Wed Feb 22 2017 <mackerel-developers@hatena.ne.jp> - 0.40.0-1
- support metadata plugins in configuration (by itchyny)
- Add metadata plugin feature (by itchyny)
- Use Named Result Parameters as document (by haya14busa)
- Set large number of file descriptors for the safety sake in init scripts (by Songmu)
- Improve darwin cpu spec (by astj)
- Fix format verb: use '%v' (by haya14busa)

* Wed Feb 08 2017 <mackerel-developers@hatena.ne.jp> - 0.39.4-1
- prepare windows eventlog (by daiksy)
- Refactor plugin configurations (by itchyny)
- Execute less `go build`s on deploy (by astj)
- treat xmlns (by mattn)
- Fix xmlns (by mattn)

* Wed Jan 25 2017 <mackerel-developers@hatena.ne.jp> - 0.39.3-1
- Fix segfault when loading a bad config file (by hanazuki)
- fix windows eventlog level when "verbose=true" (by daiksy)

* Mon Jan 16 2017 <mackerel-developers@hatena.ne.jp> - 0.39.2-1
- Test wix/pluginlist.txt on AppVeyor ci (by astj)
- Revert "remove windows plugins on pluginslist" (by daiksy)

* Thu Jan 12 2017 <mackerel-developers@hatena.ne.jp> - 0.39.1-1
- support filesystems.Ignore on windows (by Songmu)
- remove windows plugins on pluginslist (by daiksy)

* Wed Jan 11 2017 <mackerel-developers@hatena.ne.jp> - 0.39.0-1
- implement `pluginGenerators` for windows (by daiksy)
- add check-windows-eventlog on pluginlist (by daiksy)
- Remove duplicated generator in Windows (by astj)
- add mackerel-plugin-windows-server-sessions on pluginlist (by daiksy)

* Wed Dec 21 2016 <mackerel-developers@hatena.ne.jp> - 0.38.0-1
- fix typo (by ts-3156)
- Add Copyright (by yuuki)
- Separate interfaceGenerator from specGenerators (by motemen)
- Timout http reuquest in 30 sec (requries go 1.3) (by hakobe)
- specify command arguments in mackerel-agent.conf (by Songmu)
- several improvements for Windows (by daiksy)
- Avoid time.Tick and use time.NewTicker instead (by haya14busa)

* Tue Nov 29 2016 <mackerel-developers@hatena.ne.jp> - 0.37.1-1
- fix pluginlist (by daiksy)
- Suppress ec2 metadata warnings (by itchyny)
- Uncapitalize error messages (by itchyny)

* Thu Oct 27 2016 <mackerel-developers@hatena.ne.jp> - 0.37.0-1
- improve Windows support (by daiksy)

* Tue Oct 18 2016 <mackerel-developers@hatena.ne.jp> - 0.36.0-1
- don't use HTTP_PROXY when requesting cloud instance metadata APIs (by Songmu)
- Add an option to output filesystem-related metrics with key by mountpoint (by astj)

* Thu Sep 29 2016 <mackerel-developers@hatena.ne.jp> - 0.35.1-1
- support MACKEREL_PLUGIN_WORKDIR in init scripts (by Songmu)
- Add platform metadata for Darwin (by astj)
- Disable http2 for now (by Songmu)

* Wed Sep 07 2016 <mackerel-developers@hatena.ne.jp> - 0.35.0-1
- built with Go 1.7 (by Songmu)
- remove `func (vs *Values) Merge(other Values)` (by Songmu)
- [incompatible] consider df  (used + available) as size of filesystem (by Songmu)
- Remove DigitalOcean related comment/definition from spec/cloud.go (by astj)
- Fix golint is not working on ci, and add some comment to pass golint (by astj)
- Add linux distribution information to kernel spec (by ak1t0)
- http_proxy configuration (by Songmu)
- set PATH and LANG only in unix environment (by Songmu)
- Ignore docker mapper storage in spec as well (by itchyny)

* Thu Aug 18 2016 <mackerel-developers@hatena.ne.jp> - 0.34.0-1
- Reduce retry count on finding a host by the custom identifier (by itchyny)
- suppress checker flooding when resuming from sleep mode (by Songmu)
- truncate checker message up to 1024 characters (by Songmu)
- commonalize spec.FilesystemGenerator around unix OSs (by Songmu)
- define type DfStat,	remove dfColumnSpecs and refactor (by Songmu)

* Mon Aug 08 2016 <mackerel-developers@hatena.ne.jp> - 0.33.0-1
- Fill the customIdentifier in EC2 (by itchyny)

* Thu Jul 14 2016 <mackerel-developers@hatena.ne.jp> - 0.32.2-1
- fix GOMAXPROCS to 1 for avoiding rare panics (by Songmu)

* Thu Jul 07 2016 <mackerel-developers@hatena.ne.jp> - 0.32.1-1
- Add user for executing a plugin (by y-kuno)

* Thu Jun 30 2016 <mackerel-developers@hatena.ne.jp> - 0.32.0-1
- Added plugin check interval option (by karupanerura)

* Thu Jun 23 2016 <mackerel-developers@hatena.ne.jp> - 0.31.2-1
- Refactor around metrics/linux/memory (by Songmu)
- Don't stop mackerel-agent process on upgrading by debian package (by karupanerura)
- add `silent` configuration key for suppressing log output (by Songmu)
- change log level ERROR to WARNING in spec/spec.go (by Songmu)
- remove /usr/local/bin from sample.conf (by Songmu)

* Wed May 25 2016 <mackerel-developers@hatena.ne.jp> - 0.31.0-1
- Post the custom metrics to the hosts specified by custom identifiers (by itchyny)
- refactor FilesystemGenerator (by Songmu)
- Refactor metrics/linux/interface.go (by Songmu)
- remove regexp from spec/linux/cpu (by Songmu)
- Fix missing printf args (by shogo82148)

* Tue May 10 2016 <mackerel-developers@hatena.ne.jp> - 0.30.4-1
- Recover from panic while processing generators (by stanaka)
- check length of cols just to be safe in metrics/linux/disk.go (by Songmu)

* Mon May 02 2016 <mackerel-developers@hatena.ne.jp> - 0.30.3-1
- Remove usr local bin again (by Songmu)
- Fix typo (by yukiyan)
- Fix comments (by stefafafan)
- Remove go get cmd/vet (by itchyny)
- retry retirement when api request failed (by Songmu)
- output plugin stderr to log (by Songmu)

* Fri Apr 08 2016 <mackerel-developers@hatena.ne.jp> - 0.30.5-1
- Feature some3 (by stanaka)

* Fri Apr 08 2016 <mackerel-developers@hatena.ne.jp> - 0.30.4-1
- update (by stanaka)
- update (by stanaka)
- Feature some2 (by stanaka)
- update (by stanaka)

* Fri Apr 08 2016 <mackerel-developers@hatena.ne.jp> - 0.30.3-1
- update README.md (by stanaka)
- update (by stanaka)

* Fri Mar 25 2016 <y.songmu@gmail.com> - 0.30.2-1
- Revert "Merge pull request #211 from mackerelio/usr-bin" (by Songmu)

* Fri Mar 25 2016 <y.songmu@gmail.com> - 0.30.1-1
- deprecate /usr/local/bin (by Songmu)
- use GOARCH=amd64 for now (by Songmu)

* Thu Mar 17 2016 <y.songmu@gmail.com> - 0.30.0-1
- remove uptime metrics generator (by Songmu)
- Remove deprecated-sensu feature (by Songmu)
- Send all IP addresses of each interface (linux only) (by mechairoi)
- add `init` subcommand (by Songmu)
- Refactor net interface (multi ip support and bugfix) (by Songmu)
- Stop to fetch flags of cpu in spec/linux/cpu (by Songmu)

* Mon Mar 07 2016 <y.songmu@gmail.com> - 0.29.2-1
- Don't overwrite mackerel-agent.conf when updating deb package (Fix deb packaging) (by Songmu)

* Fri Mar 04 2016 <y.songmu@gmail.com> - 0.29.1-1
- maintenance release

* Wed Mar 02 2016 <y.songmu@gmail.com> - 0.29.0-1
- remove deprecated command line options (-version and -once) (by Songmu)
- Report checker execution timeout as Unknown status (by hanazuki)

* Thu Feb 18 2016 <stefafafan@hatena.ne.jp> - 0.28.1-1
- fix the exit status on stopping the agent in the init script of debian (by itchyny)

* Thu Feb 04 2016 <y.songmu@gmail.com> - 0.28.0-1
- add a configuration to ignore filesystems (by stanaka)
- fix the code of extending the process's environment (by itchyny)
- s{code.google.com/p/winsvc}{golang.org/x/sys/windows/svc} (by Songmu)
- Max check attempts option for check plugin (by mechairoi)

* Fri Jan 08 2016 <y.songmu@gmail.com> - 0.27.1-1
- [bugfix] fix timeout interval when calling `df` (by Songmu)

* Wed Jan 06 2016 <y.songmu@gmail.com> - 0.27.0-1
- use timeout when calling `df` (by Songmu)
- Notification Interval for check monitoring (by itchyny)

* Thu Dec 10 2015 <y.songmu@gmail.com> - 0.26.2-1
- output success message to stderr when configtest succeed (by Songmu)

* Wed Dec 09 2015 <y.songmu@gmail.com> - 0.26.1-1
- fix deprecate message (by Songmu)

* Tue Dec 08 2015 <y.songmu@gmail.com> - 0.26.0-1
- Make HostID storage replacable (by motemen)
- Publicize command.Context's fields (by motemen)
- Configtest (by fujiwara)
- Refactor config loading and check if Apikey exists in configtest (by Songmu)
- fix exit status of debian init script. (by fujiwara)
- Deprecate version and once option (by Songmu)

* Wed Nov 25 2015 <y.songmu@gmail.com> - 0.25.1-1
- Go 1.5.1 (by Songmu)
- logging STDERR of checker command (by Songmu)

* Thu Nov 12 2015 <y.songmu@gmail.com> - 0.25.0-1
- Retrieve interfaces on Darwin (by itchyny)
- add NetBSD support. (by miwarin)

* Thu Nov 05 2015 <y.songmu@gmail.com> - 0.24.1-1
- We are Mackerel (by itchyny)

* Mon Oct 26 2015 <daiksy@hatena.ne.jp> - 0.24.0-1
- define config.agentName and set proper config path (by Songmu)
- /proc/cpuinfo parser for old ARM Linux kernels (by hanazuki)
- os.MkdirAll() before creating pidfile (by Songmu)

* Wed Sep 30 2015 <ttsujikawa@gmail.com> - 0.23.1-1
- Code signing for windows installer (by mechairoi)

* Mon Sep 14 2015 <itchyny@hatena.ne.jp> - 0.23.0-1
- send check monitor report to server when check script failed even if the monitor result is not changed (by Songmu)
- Correct sample nginx comment. (by kamatama41)

* Wed Sep 02 2015 <tomohiro68@gmail.com> - 0.22.0-1
- add `reload` to init scripts (by Songmu)

* Wed Sep 02 2015 <tomohiro68@gmail.com> - 0.21.0-1
- Exclude mkr binary from deb/rpm package (by Sixeight)

* Thu Aug 13 2015 <tomohiro68@gmail.com> - 0.20.1-1
- use C struct for accessing Windows APIs (by stanaka)
- Fix bug that checks is not removed when no checks. (by Sixeight)

* Wed Jul 29 2015 <y.songmu@gmail.com> - 0.20.0-1
- support subcommand (by Songmu)
- remove trailing newline chars when loading hostID (by Songmu)
- add sub-command `retire` and support $AUTO_RETIREMENT in initd (by Songmu)
- add postinst to register mackerel-agent to start-up (deb package) (by stanaka)
- bump bundled mkr version to 0.3.1 (by Songmu)

* Wed Jul 22 2015 <y.songmu@gmail.com> - 0.19.0-1
- Support gce meta (by Songmu)
- Valid pidfile handling (fix on darwin) (by Songmu)
- -once only takes one second (by Songmu)
- fix shutdown priority in rpm/src/mackerel-agent.initd (by Songmu)

* Thu Jul 16 2015 <y.songmu@gmail.com> - 0.18.1-1
- s/ami_id/ami-id/ in spec/cloud.go (by Songmu)
- remove `UpdateHost()` process from `prepareHost()` for simplicity (by Songmu)
- filter invalid roleFullNames with warning logs (by Songmu)
- allow using spaces as delimiter for custom metric values (by Songmu)

* Wed Jul 08 2015 <tomohiro68@gmail.com> - 0.18.0-1
- Retry in prepare (by Songmu)
- [WORKAROUND] downgrade golang version for windows (by Sixeight)

* Wed Jun 17 2015 <tomohiro68@gmail.com> - 0.17.1-1
- Update to go 1.4.2 for windows build (by mechairoi)

* Wed Jun 10 2015 <tomohiro68@gmail.com> - 0.17.0-1
- Set `displayName` via agent (by Sixeight)
- refactoring around api access (by Songmu)
- Configurable host status on start/stop agent (by Songmu)
- Add an agent memory usage metrics generator for diagnostic use (by hakobe)
- Add mkr to deb/rpm package (by Sixeight)

* Tue May 12 2015 <y.songmu@gmail.com> - 0.16.1-1
- Code sharing around dfValues (by Songmu)
- [FreeBSD] Fix 'panic: runtime error: index out of range'. (by iwadon)
- separete out metrics/darwin/swap.go from memory.go (by Songmu)

* Fri May 08 2015 <y.songmu@gmail.com> - 0.16.0-1
- suppress logging (by stanaka)
- "Check" functionality (by motemen)
- update for windows (by daiksy)
- collect memory metrics of osx (by Songmu)
- Send plugin.check._name_s list on `updateHost` (by mechairoi)

* Thu Apr 02 2015 <y.songmu@gmail.com> - 0.15.0-1
- building packages (by Songmu)
- Only skip device mapper created by docker (Resolve #70) (by mechairoi)
- Run once and output results to stdout (by stanaka)
- introduce Songmu/timeout for interrupting long time plugin execution (by Songmu)
- add config.apibase (by Songmu)
- output GOOS GOARCH runtime.Version() when -version option is specified (by Songmu)
* Mon Mar 23 2015 Songmu <songmu@hatena.ne.jp> 0.14.3-1
- [enhancement] add collector for ec2 metadata (stanaka)
* Tue Jan 20 2015 Songmu <songmu@hatena.ne.jp> 0.14.1-1
- [fix] skip device mapper metrics
- [fix] filter invalid float values
- [enhancement] testing
- [enhancement] collect more metrics about darwin and freebsd
* Thu Dec 25 2014 Songmu <songmu@hatena.ne.jp> 0.14.0-1
- [improve] wait for termination until postQueue is empty up to 30 seconds.
- [improve] wait up to 30 seconds before initial posting
- [feature] work on Windows darwin FreeBSD (unofficial support)
* Tue Nov 18 2014 y_uuki <y_uuki@hatena.ne.jp> 0.13.0-1
- [feature] Support `-version` flag
- [improve] Do bulk posting metrics when retrying metrics sending
- [feature] Support darwin
* Wed Oct  1 2014 skozawa <skozawa@hatena.ne.jp> 0.12.3-1
- [fix] Fixed index out of rage for diskstats
- [improve] Update hostname on updating host specs
* Tue Sep 16 2014 y_uuki <y_uuki@hatena.ne.jp> 0.12.2-3
- [fix] Add validation if pidfile is invalid
* Mon Sep  8 2014 skozawa <skozawa@hatena.ne.jp> 0.12.2-2
- [fix] Add a process name to killproc
* Fri Sep  5 2014 skozawa <skozawa@hatena.ne.jp> 0.12.2-1
- [fix] change retry and dequeue delay time
* Thu Aug 21 2014 motemen <motemen@hatena.ne.jp> 0.12.1-1
- Extended retry queue
* Tue Aug 19 2014 motemen <motemen@hatena.ne.jp> 0.12.0-1
- [breaking] Changed custom metric plugins' meta information format to JSON instead of TOML
- [feature] Added filesystem metrics
* Wed Aug  6 2014 motemen <motemen@hatena.ne.jp> 0.11.1-1
- [fix] Fixed non-critical log message when plugin meta loading
* Wed Aug  6 2014 motemen <motemen@hatena.ne.jp> 0.11.0-1
- [feature] Including config files with 'include' key
* Tue Aug  5 2014 motemen <motemen@hatena.ne.jp> 0.10.1-1
- [fix] Fixed issue that environment variable was not set
* Tue Aug  5 2014 motemen <motemen@hatena.ne.jp> 0.10.0-1
- [feature] Added support for custom metric schemata
* Wed Jul  9 2014 skozawa <skozawa@hatena.ne.jp> 0.9.0-2
- [fix] Removed unused metrics #20
- [feature] Add configurations for posting metrics #19
- [fix] Prevent exiting without cleaning pidfile #18
* Tue Jun 24 2014 hakobe <hakobe@hatena.ne.jp> 0.8.0-1
- [improve] Using go 1.3
- [feature] Periodically update host specs #15
- [fix] Http request now have timeout #17
* Fri Jun  6 2014 motemen <motemen@hatena.ne.jp> 0.7.0-1
- [fix] Replace invalid characters (e.g. '.') in disk and interface names with underscores
- [fix] Removed deprecated metrics
* Fri May 23 2014 hakobe <hakobe@hatena.ne.jp> 0.6.1-1
- [breaking change] Automatically add 'custom.' prefix to the name of custom metrics
- [change] Change the key to configure custom metrics from "sensu.checks." to "plugin.metrics." in the config file
- [improve] More friendly and consistent error messages
- [fix] Change the permission of /var/lib/mackerel-agent directory to 755
- [fix] Change the permission of /etc/init.d/mackerel-agent to 755
* Wed May 14 2014 motemen <motemen@hatena.ne.jp> 0.5.1-3
- [fix] Fixed init script not to use APIKEY if empty
* Tue May 13 2014 motemen <motemen@hatena.ne.jp> 0.5.1-2
- Updated version string
* Tue May 13 2014 motemen <motemen@hatena.ne.jp> 0.5.1-1
- [improve] Warn and exit on startup if no API key given
- [fix] Support parsing large disk sizes
- [fix] Trap SIGHUP not to die
- [fix] Continue running even if failed to collect host specs
- [fix] Use binaries under /sbin/ and /bin/ to generate specs/metrics
* Thu May  8 2014 hakobe <hakobe@hatena.ne.jp> 0.5.0-1
- [improve] Verbose option now prints debug information
- [misc] Changed license from Test-use only to Commercial
* Wed May  7 2014 hakobe <hakobe@hatena.ne.jp> 0.4.3-1
- [fix] Changed sleep time for buffered requests
* Wed Apr 30 2014 hakobe <hakobe@hatena.ne.jp> 0.4.2-1
- [fix] Fixed a memory leak when metrics collection unexpectedly blocked
* Mon Apr 28 2014 mechairoi <mechairoi@hatena.ne.jp> 0.4.1-1
- [fix] Fixed a crash when increasing or decreasing disks or interfaces
* Fri Apr 25 2014 skozawa <skozawa@hatena.ne.jp> 0.4.0-1
- [improve] Change interval for disk, cpu and interface metrics
* Wed Apr 23 2014 hakobe <hakobe@hatena.ne.jp> 0.3.0-2
- [fix] Exclude log files from package
- [fix] Remove an unncecessary setting sample
* Tue Apr 22 2014 mechairoi <mechairoi@hatena.ne.jp> 0.3.0-1
- [improve] Update interfaces information each start
- [improve] Set nice 'User-Agent' header
- [improve] Add 'memory.used' metrics
- [improve] Execute sensu command through 'sh -c'
- [fix] Fix interval of collecting metrics
- [fix] Fix crashes when collecting disk usage
* Thu Apr 17 2014 skozawa <skozawa@hatena.ne.jp> 0.2.0-2
- Fix config file comments
* Wed Apr 16 2014 motemen <motemen@hatena.ne.jp> 0.2.0-1
- [feature] Add support for sensu plugins
- [feature] Buffer metric values in case of request error
* Wed Apr 9 2014 motemen <motemen@hatena.ne.jp> 0.1.1-2
- Add mackerel-agent.conf
- Use 32-bit binary
* Wed Apr 9 2014 mechairoi <mechairoi@hatena.ne.jp> 0.1.1-1
- New features
* Fri Apr 4 2014 hakobe932 <hakobe932@hatena.ne.jp> 0.1.0-1
- New features
* Tue Mar 31 2014 y_uuki <y_uuki@hatena.ne.jp> 0.0.2-2
- Add logrotate.
* Tue Mar 25 2014 y_uuki <y_uuki@hatena.ne.jp> 0.0.2-1
- New features.
* Fri Mar 7 2014 y_uuki <y_uuki@hatena.ne.jp> 0.0.1-1
- Initial spec.
