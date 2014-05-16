mackerel-agent
===============

![agent-si](docs/images/agent-si.png "mackerel-agent")

mackerel-agent is an agent program to post your hosts' metrics to [Mackerel](https://mackerel.io/).

mackerel-agent executes the following tasks in foreground:
- Register your host to Mackerel
- Collect and post specs and metrics of your host to Mackerel periodically

You can see information of the host on [Mackerel](https://mackerel.io/).

For now, mackerel-agent is guaranteed to run only on CentOS 5/6 and Debian 6/7.

SYNOPSIS
--------

Build and Run mackerel-agent.

```
make build
make run
```

The `apikey` is required to run the agent.

Create an organization in [Mackerel](https://mackerel.io/) and configure `apikey` in `mackerel-agent.conf`.


The following commands can be used instead of `make`.

```
go get -d github.com/mackerelio/mackerel-agent
go build -o build/mackerel-agent \
  -ldflags="\
    -X github.com/mackerelio/mackerel-agent/version.GITCOMMIT `git rev-parse --short HEAD` \
    -X github.com/mackerelio/mackerel-agent/version.VERSION   `git describe --tags --abbrev=0 | sed 's/^v//' | sed 's/\+.*$$//'` " \
  github.com/mackerelio/mackerel-agent
./build/mackerel-agent -conf=mackerel-agent.conf
```


Test
----------

Test mackerel-agent.

The agent collects information about a host which the agent run.

```
make test
```

