mackerel-agent for FreeBSD
==========================

1. Get binary for FreeBSD at release page. https://github.com/mackerelio/mackerel-agent/releases
    - `amd64`: `mackerel-agent_freebsd_amd64.tar.gz`
    - `i386` : `mackerel-agent_freebsd_386.tar.gz`
    - (`arm` : `mackerel-agent_freebsd_arm.tar.gz`)
2. Extract file, and `cd`.
    - `tar -xzvf mackerel-agent_freebsd_*.tar.gz`
3. Copy `mackerel-agent` to `/usr/local/bin/`.
    - `(sudo) cp mackerel-agent /usr/local/bin/`
4. Edit `mackerel-agent.conf`, then copy it to `/usr/local/etc/`.
    - `${EDITOR} mackerel-agent.conf`
    - `(sudo) cp mackerel-agent.conf /usr/local/etc/`
5. Copy `mackerel_agent` (this directory file; rc script) to `/usr/local/etc/rc.d/`.
    - `(sudo) cp mackerel_agent /usr/local/etc/rc.d/`
6. Add `mackerel_agent_enable="YES"` at `/etc/rc.conf` (or use `sysrc`).
    - `(sudo) sysrc mackerel_agent_enable="YES"`
    - (or `sudoedit /etc/rc.conf`)
7. Start `mackerel_agent`.
    - `(sudo) service mackerel_agent start`

