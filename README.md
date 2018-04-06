mackerel-agent
===============

![agent-si](docs/images/agent-si.png "mackerel-agent")

`mackerel-agent` is a client software for [Mackerel](https://mackerel.io/).
[Mackerel](https://mackerel.io/) is an online visualization and monitoring service for servers.

Once `mackerel-agent` is installed, it runs the following tasks on the installed host in foreground:
- register your hosts to Mackerel
- collect specs and metrics of those hosts and post them to Mackerel

Collected information will be visualized on [Mackerel](https://mackerel.io/).

PREREQUISITES
-------------

You have to create an organization on [Mackerel](https://mackerel.io/) at first.
After that, specify `apikey` value in `mackerel-agent.conf` with the following command.

```
% mackerel-agent init -apikey {{YOUR_APIKEY}}
```

SYNOPSIS
--------

Build and Run the mackerel-agent.

```console
% make build
% make run
```

You can run the following commands instead of using `make`.

```console
% go get -d github.com/mackerelio/mackerel-agent
% go build -o build/mackerel-agent \
  -ldflags="\
    -X github.com/mackerelio/mackerel-agent/version.GITCOMMIT `git rev-parse --short HEAD` \
    -X github.com/mackerelio/mackerel-agent/version.VERSION   `git describe --tags --abbrev=0 | sed 's/^v//' | sed 's/\+.*$$//'` " \
  github.com/mackerelio/mackerel-agent
./build/mackerel-agent -conf=mackerel-agent.conf
```

### On Windows

Use `.bat` files instead of `make` commands.

```console
% build.bat
```

Test
----------

Test mackerel-agent to confirm it's working properly.

The agent will collect information about the host on which it has been installed.

```console
% make test
```

License
----------
```
Copyright 2014 Hatena Co., Ltd.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
