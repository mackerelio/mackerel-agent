# sudo yum -y install rpmdevtools go && rpmdev-setuptree
# rpmbuild -ba ~/rpmbuild/SPECS/mackerel-agent.spec

%define _binaries_in_noarch_packages_terminate_build   0

Name:      mackerel-agent
Version:   %{_version}
Release:   1%{?dist}
License:   Commercial
Summary:   mackerel.io agent
URL:       https://mackerel.io
Group:     Hatena
Source0:   %{name}.sysconfig
Source1:   %{name}.conf
Source2:   %{name}.service
Packager:  Hatena
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
* Thu Mar 09 2017 <mackerel-developers@hatena.ne.jp> - 0.41.1-1
- add check-tcp on pluginlist.txt (by daiksy)
- use new bot token (by daiksy)
- use new bot token (by daiksy)

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

