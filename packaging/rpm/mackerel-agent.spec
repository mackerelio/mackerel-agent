# sudo yum -y install rpmdevtools go && rpmdev-setuptree
# rpmbuild -ba ~/rpmbuild/SPECS/mackerel-agent.spec

%define _binaries_in_noarch_packages_terminate_build   0

Name:      mackerel-agent
Version:   %{_version}
Release:   1
License:   Commercial
Summary:   mackerel.io agent
URL:       https://mackerel.io
Group:     Hatena
Source0:   %{name}.initd
Source1:   %{name}.sysconfig
Source2:   %{name}.logrotate
Source3:   %{name}.conf
Packager:  Hatena
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
