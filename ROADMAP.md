# opnsenseHealth Roadmap

OPNsense health plugin. Feature parity with monokit1 `opnsenseHealth/`.

## Scaffolding

- [X] Plugin skeleton from pluginTemplate (go.mod, justfile, Containerfile, test harness)
- [X] Example config `config/opnsense.yml`
- [ ] Config struct + `opnsense.yml` case wired into monokit_lib
  - [ ] Per-subsystem `enabled` toggles (wireguard, ipsec, gateway, dns, carp) where an
        omitted toggle means enabled, so adding a check does not require editing every
        existing opnsense.yml
- [ ] Podman integration tests (TLS endpoint fixture)
  - [ ] Fixtures for the parsers: `ifconfig` CARP output, `wg show all dump`,
        `list_status.py` JSON, `gateway_status.php` JSON, `config.xml`
- [ ] Auto-detection: require `opnsense.yml` **and** `/usr/local/opnsense` or
      `/usr/local/sbin/opnsense-version`, so deleting the config disables the plugin

## Features

### SSL / web UI

- [ ] OPNsense web domain resolution (auto-detect when not configured)
- [ ] TLS handshake against the web UI and certificate inspection
  - [ ] Port fallback: try the configured `port`, else 9443 then 443
  - [ ] SNI set from the domain parsed out of config.xml
  - [ ] Alarm + Redmine issue: SSL connection failure
  - [ ] Alarm + Redmine issue: no certificate presented
  - [ ] Alarm + Redmine issue: certificate expired
  - [ ] Alarm + Redmine issue: certificate expiring soon (expiration-limit days)
  - [ ] Subject / issuer / expiry / days-remaining reported (CN preferred over full DN)
  - [ ] Recovery closes the Redmine issue
  - [ ] Separate alarm when Redmine issue creation itself fails

### config.xml parsing

The parsed config supplies the SNI hostname, the WireGuard peer and IPSec
connection display names, the CARP VIP→vhid map and Unbound's state, so it is
read once per run and shared by every other check.

- [ ] Web domain / hostname extraction
- [ ] WireGuard server list with `carp_depend_on` per instance, and peer names by public key
- [ ] IPSec connection descriptions keyed by UUID
- [ ] CARP VIP UUID→vhid map, and whether CARP is configured at all
- [ ] Unbound enabled/disabled state and its configured port

### CARP (HA)

- [ ] CARP VIP state per interface from `ifconfig` (MASTER / BACKUP / INIT / DISABLED, vhid)
- [ ] `IsBackupForTunnels`: true only when every VIP a tunnel depends on is BACKUP —
      a box that is MASTER for some groups and BACKUP for others must not be treated
      as globally BACKUP
- [ ] Fail safe: an unreadable CARP state reports not-backup so tunnel alarms keep firing
      rather than a broken lookup silencing the fleet

### WireGuard

- [ ] Per-interface health from `wg show all dump` (interface lines 5 fields, peer lines 9)
- [ ] Interface up/RUNNING state from `ifconfig` flags, raw flag list in the alarm body
- [ ] Missing-interface detection: enabled in OPNsense but no device present
- [ ] Peer handshake age vs `handshake_timeout` (default 300s), only for peers this box
      dials out to (`persistent-keepalive` set)
- [ ] Idle dial-in peers reported but never alarmed — an offline client is not a fault,
      and `wg` resets latest-handshake on every device recreation
- [ ] Problem classification per peer: NO_HANDSHAKE / STALE
- [ ] `excluded_peers` (matched by name or public key), clearing any alarm already held
- [ ] One alarm per interface covering link state and peers together, so a downed
      interface does not fan out an alarm per peer
- [ ] CARP suppression per interface via the instance's `carp_depend_on`, releasing
      (not just skipping) the alarm so a box dropping to BACKUP clears its own faults
- [ ] Skip entirely when WireGuard is switched off in config.xml — the binary ships
      regardless, so binary presence is not the signal

### IPSec

- [ ] Phase 1 state per connection from OPNsense's `list_status.py` (child SAs are
      deliberately not alarmed — far too noisy)
- [ ] Most-favourable SA state wins, so a stale SA left by a rekey does not read as down
- [ ] Connection ID as the alarm key (GUI descriptions change), display name resolved
      from config.xml
- [ ] `excluded_connections` (description or UUID), clearing any alarm already held
- [ ] Script failure only counts as a fault when a tunnel is actually configured —
      `list_status.py` ships on every box and fails when charon is not running
- [ ] Empty output treated as "no connections loaded", not a broken script
- [ ] Per-connection decode so one unexpected field cannot discard every other tunnel
- [ ] Partial-unreadability alarm rather than silently dropping entries
- [ ] CARP suppression sharing the tunnel-wide verdict (IPSec has no `carp_depend_on`)

### Gateways

- [ ] Per-gateway health from `gateway_status.php` (address, monitor, status, loss, delay, stddev)
- [ ] Packet loss vs `loss_limit` (default 20%)
- [ ] `force_down` treated as an admin action: release the alarm, do not fault
- [ ] `Pending` (dpinger has no samples yet) leaves the existing alarm state untouched —
      calling recovery here would emit a false clear for a gateway that is still down
- [ ] Gateways without a monitor have no loss figure and are not compared
- [ ] Readability alarm for the script itself, separate from per-gateway alarms

### DNS

- [ ] Real resolver query against the local Unbound — the only way to tell "answering"
      from "merely running"
- [ ] Resolver address: explicit `dns.server`, else Unbound's own configured port from
      config.xml, else 127.0.0.1:53
- [ ] Configurable query name (default `opnsense.org`), 5s timeout, duration reported
- [ ] Skip when Unbound is disabled in config.xml
- [ ] Distinct failure modes: lookup error vs zero records returned

### Shared execution / parsing concerns

- [ ] Binary resolution preferring `$PATH` with OPNsense absolute-path fallbacks
      (`ifconfig`, `php`, `wg`), where absence means "feature not present"
- [ ] Command timeout with stderr folded into the error so a failing script explains
      itself in the alarm body
- [ ] Strip non-JSON preamble before decoding — the PHP CLI and Python both write
      notices to stdout, and a strict decode would turn an OPNsense upgrade into a
      fleet-wide false alarm
- [ ] Tolerate every payload shape the status scripts have emitted across releases:
      object keyed by id, bare list, and `{"items": [...]}`
- [ ] Fields that arrive as string, number or list depending on strongSwan version
- [ ] Alarm-key sanitisation for names with spaces/punctuation

### Other

- [ ] OPNsense version reporting (monokit1: via versionCheck `opnsense-version`; monokit2: osHealth vlib `opnsense.go`)
- [ ] Health summary box output (depends on the lib renderer)
- [ ] Health data POST to the server API (depends on base client/server API)
