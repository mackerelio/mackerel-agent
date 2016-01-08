# Changelog

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
