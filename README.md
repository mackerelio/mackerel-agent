mackerel-agent
===============

![agent-si](docs/images/agent-si.png "mackerel-agent")

mackerel-agent is a piece of software which is installed on your hosts to collect metrics and events and send them to [Mackerel](https://mackerel.io/) where they can be visualized and monitored.

mackerel-agent executes the following tasks in the foreground:
- registering your hosts with Mackerel
- collecting specs and metrics from your hosts and posting them to Mackerel

Your hosts' information will be viewable on [Mackerel](https://mackerel.io/).

SYNOPSIS
--------

Build and Run the mackerel-agent.

```
make build
make run
```

The `apikey` will be required in order to run the agent.

An organization must first be created in [Mackerel](https://mackerel.io/), then the `apikey` can be configured in `mackerel-agent.conf`.

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

To build on Windows use the command ```build.bat```

To run on Windows use the command ```run.bat```


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
