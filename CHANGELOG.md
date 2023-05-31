# Changelog
All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [unreleased]
### Security
### Added
### Changed
- Replace deprecated CPUShares and MemoryLimit control options in the systemd unit with CPUWeight and MemoryHigh ([#20](https://github.com/flatcar/locksmith/pull/20))

### Removed

## [v0.7.0](https://github.com/flatcar/locksmith/releases/tag/v0.7.0)- 30/11/2021

### Security
- Update Go dependencies to latest versions ([#6](https://github.com/flatcar/locksmith/pull/6), [#8](https://github.com/flatcar/locksmith/pull/8))

### Changed
- Use Go modules to handle dependencies ([#4](https://github.com/flatcar/locksmith/pull/4))
- Use `github.com/flatcar/locksmith` as Go module name ([#2](https://github.com/flatcar/locksmith/pull/2))
- Node will not reboot outside of reboot window in case lock is acquired outside of it ([#10](https://github.com/flatcar/locksmith/pull/10))

### Fixed
- Possible panic on shutdown ([#1](https://github.com/flatcar/locksmith/pull/1))


### v0.6.2

Locksmith now internally handles DBus signals in a safer way, with additional
checks on unexpected/malformed ones. This avoids a panic in the daemon component,
which was frequently occurring on shutdown.

### v0.6.1

Locksmith now sequences itself after update_engine, avoiding an exit and
relaunch if update_engine has not yet started.

### v0.6.0

Locksmith no longer supports the deprecated `best-effort` strategy. The default
strategy is now `reboot`.

### v0.5.0

Locksmith now writes an informational metadata file at the well known path
`/etc/update-engine/coordinator.conf`.

Other tools that fulfill a similar role may also safely write to that file by
ensuring an exclusive lock is held on it for the duration of them running.

### v0.4.2

Locksmith now uses Restart=on-failure in the systemd unit, so that if the
reboot strategy is off, it is not constantly restarted.

### v0.4.1

Locksmith no longer requires update-engine.service and does not have an
ordering dependency on user-config.target and system-config.target due to the
dependency loop when using coreos-cloudinit.

### v0.4.0

Locksmith now uses github.com/coreos/etcd/client, instead of the deprecated
github.com/coreos/go-etcd.

Locksmith now uses github.com/coreos/pkg/capnslog for logging.

The reboot strategy `best-effort` is deprecated, and locksmithd will complain
loudly if it is used. Please use an explicit `reboot` or `etcd-lock` strategy
instead.

Locksmith logs some information about the configured reboot window.

Locksmith supports etcd basic auth.

Locksmith again requires update-engine.service, and will start after
user-config.target system-config.target are reached.

### v0.3.4

The environment variables controlling reboot windows (`REBOOT_WINDOW_START`,
`REBOOT_WINDOW_LENGTH`) have been renamed to include the prefix `LOCKSMITHD_`
to maintain consistency with other locksmithd environment variables.

The old environment variables are still read to maintain compatibility with
locksmithd v0.3.1 to v0.3.3.

### v0.3.3

Remove dependency on update-engine.service from locksmithd.service. If
update-engine failed to start, systemd wouldn't start locksmith and the restart
logic only applies if the service can be started.

### v0.3.2

Set GOMAXPROCS=1 in the locksmithd systemd service to keep behavior consistent between builds using Go 1.5 and previous versions.

### v0.3.1

v0.3.1 is the first release with a changelog :-)

There are also a number of new features in this release, including [groups](README.md#groups), an [`off` strategy](README.md#configuration), and [reboot windows](README.md#reboot-windows).

Full list of changes since v0.3.0:
- New features
  - "groups" feature, facilitating partitioned co-ordinating of reboots (#70)
  - "off" strategy, which will cause locksmith to perform no action and shut itself down (#79)
  - reboot windows, allowing control over when reboots occur (#80)
- Bug fixes
  - daemon now considers strategy when attempting to unlock, rather than just blindly checking the local etcd's activeness (#86)
  - updateengine client no longer attempts to close a dbus connection if the connecion failed (#83)
- Other changes
  - greater verbosity of error messages in the case of unlocking failures (#82)
