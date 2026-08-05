# opnsenseHealth Roadmap

OPNsense health plugin. Feature parity with monokit1 `opnsenseHealth/`.

## Scaffolding

- [X] Plugin skeleton from pluginTemplate (go.mod, justfile, Containerfile, test harness)
- [X] Example config `config/opnsense.yml`
- [ ] Config struct + `opnsense.yml` case wired into monokit_lib
- [ ] Podman integration tests (TLS endpoint fixture)

## Features

- [ ] OPNsense web domain resolution (auto-detect when not configured)
- [ ] TLS handshake against the web UI and certificate inspection
  - [ ] Alarm + Redmine issue: SSL connection failure
  - [ ] Alarm + Redmine issue: no certificate presented
  - [ ] Alarm + Redmine issue: certificate expired
  - [ ] Alarm + Redmine issue: certificate expiring soon (expiration-limit days)
  - [ ] Recovery closes the Redmine issue
  - [ ] Separate alarm when Redmine issue creation itself fails
- [ ] OPNsense version reporting (monokit1: via versionCheck `opnsense-version`; monokit2: osHealth vlib `opnsense.go`)
