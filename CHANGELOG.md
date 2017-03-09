# Changelog

## 0.41.1 (2017-03-09)

* add check-tcp on pluginlist.txt #351 (daiksy)
* use new bot token #352 (daiksy)
* use new bot token #354 (daiksy)


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
