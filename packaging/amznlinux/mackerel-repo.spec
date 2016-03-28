Name: mackerel-repo
Version: 1
Release: 0
Summary: Mackerel Repo
License: Commercial
URL: http://mackerel.io
Source0: mackerel.repo
BuildRoot: %{_tmppath}/%{name}-%{version}-%{release}-root
BuildArch: noarch

%description
mackerel.io repository

%prep

%build

%install
rm -rf $RPM_BUILD_ROOT

install -dm 755 $RPM_BUILD_ROOT%{_sysconfdir}/yum.repos.d
install -pm 644 %{SOURCE0} $RPM_BUILD_ROOT%{_sysconfdir}/yum.repos.d

%clean
rm -rf $RPM_BUILD_ROOT

%files
%defattr(-,root,root,-)
%config(noreplace) /etc/yum.repos.d/*

%changelog
* Mon Mar 28 2016 Shinji Tanaka <stanaka@hatena.ne.jp>
– Create Package
