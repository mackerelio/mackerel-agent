mackerel-agent
===============

mackerel-agentは[Mackerel](https://mackerel.io/)にホストを登録するためのエージェントです。

mackerel-agentはフォアグラウンドで以下の動作をします。
- ホストの情報をMackerelに登録
- ホスト上で定期的にリソース情報を収集し、Mackerelに投稿

投稿されたホスト情報、リソース情報はMackerelのウェブインターフェースで確認できます。

いまのところ[一部のLinux](http://help-ja.mackerel.io/entry/overview)上での動作しか保証していません。


SYNOPSIS
--------

エージェントをビルドして、実行します。

```
make build
make run
```

エージェントの実行には`apikey`が必要です。

[Mackerel](https://mackerel.io/)でオーガニゼーションを作成し、`apikey`を`mackrel-agent.conf`に記述してください。


makeを使わず、直接起動することもできます。

```
go get -d github.com/mackerelio/mackerel-agent
go build -o build/mackerel-agent \
  -ldflags="\
    -X github.com/mackerelio/mackerel-agent/version.GITCOMMIT `git rev-parse --short HEAD` \
    -X github.com/mackerelio/mackerel-agent/version.VERSION   `git describe --tags --abbrev=0 | sed 's/^v//' | sed 's/\+.*$$//'` " \
  github.com/mackerelio/mackerel-agent
./build/mackerel-agent -conf=mackerel-agent.conf
```

詳しくは[mackerel-agent仕様](http://help-ja.mackerel.io/entry/spec/agent)をご覧ください。

Windowsのビルドの場合は```build.bat```を実行してください。

Windowsのエージェントの実行には```run.bat```を実行してください。

Test
------

エージェントのテストを行います。

実行したホストの情報を収集します。

```
make test
```

