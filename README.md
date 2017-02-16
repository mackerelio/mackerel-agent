mackerel-agent
===============

![agent-si](docs/images/agent-si.png "mackerel-agent")

`mackerel-agent` is a client software for [Mackerel](https://mackerel.io/).
[Mackerel](https://mackerel.io/) is an online visualization and monitoring service for servers.

Once `mackerel-agent` is installed, it runs following tasks on the installed host in foreground:
- register your hosts to Mackerel
- collect specs and metrics of the host and post them to Mackerel

Collected information will be visualized on [Mackerel](https://mackerel.io/).

As of now, mackerel-agent is officially supported to run on Amazon Linux, CentOS 5/6/7, Ubuntu 12.04LTS/14.04LTS, Debian 6/7 or Windows Server 2008 R2 and later 32-bit/64-bit environments.

PREREQUISITES
-------------

You have to create an organization on [Mackerel](https://mackerel.io/) at first.
After that, set `apikey` in `mackerel-agent.conf`.


SYNOPSIS
--------

Build and Run the mackerel-agent.

```
make build
make run
```

You can run following commands instead of `make`.

```
go get -d github.com/mackerelio/mackerel-agent
go build -o build/mackerel-agent \
  -ldflags="\
    -X github.com/mackerelio/mackerel-agent/version.GITCOMMIT `git rev-parse --short HEAD` \
    -X github.com/mackerelio/mackerel-agent/version.VERSION   `git describe --tags --abbrev=0 | sed 's/^v//' | sed 's/\+.*$$//'` " \
  github.com/mackerelio/mackerel-agent
./build/mackerel-agent -conf=mackerel-agent.conf
```

### on Windows

Use `.bat` files instead of `make` commands.

```
build.bat
run.bat
```

Test
----------

Test mackerel-agent to confirm it's working properly.

The agent will collect information about the host on which it has been installed.

```
make test
```

License
----------

Copyright 2014 Hatena Co., Ltd.

Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License. You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the License for the specific language governing permissions and limitations under the License.
