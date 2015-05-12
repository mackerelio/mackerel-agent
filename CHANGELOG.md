# Changelog

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
