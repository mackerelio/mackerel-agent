mackerel-agent for FreeBSD
==========================

## Install from ports (recommended)

mackerel-agent is now available  FreeBSD Ports Collection.

### Install
```
$ sudo make -C /usr/ports/sysutils/mackerel-agent
or
$ sudo portmaster sysutils/mackerel-agent
or
$ sudo pkg install mackerel-agent
```

### Configure
```
$ sudoedit /usr/local/etc/mackerel-agent/mackerel-agent.conf
```

### Register as startup and start agent

```
$ sudo sysrc mackerel_agent_enable=YES
$ sudo service mackerel_agent start
```

## Install from tarball

### Fetch and extract

Fetch the release tarball matches your architecture.

```
$ fetch https://github.com/mackerelio/mackerel-agent/releases/download/v0.67.1/mackerel-agent_freebsd_amd64.tar.gz
$ tar zxfv mackerel-agent_freebsd_amd64.tar.gz
$ cd mackerel-agent_freebsd_amd64
```

The release tarball doesn't include rc script so far. Fetch it from the git repository.

```
$ fetch https://raw.githubusercontent.com/mackerelio/mackerel-agent/master/packaging/freebsd/mackerel_agent
```

### Install

```
$ sudo install -d /usr/local/etc/mackerel-agent
$ sudo install -m 0600 mackerel-agent.conf /usr/local/etc/mackerel-agent
$ sudo install -m 555 mackerel-agent /usr/local/bin
$ sudo install -m mackerel_agent /usr/local/etc/rc.d
```

### Configure

```
$ sudoedit /usr/local/etc/mackerel-agent/mackerel-agent.conf
```

### Register as startup and start agent

```
$ sudo sysrc mackerel_agent_enable=YES
$ sudo service mackerel_agent start
```
