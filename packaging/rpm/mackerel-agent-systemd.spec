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

