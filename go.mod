module github.com/flatcar-linux/locksmith

go 1.14

replace github.com/coreos/etcd => ./etcd

replace github.com/coreos/etcd/client => ./etcd/client

require (
	github.com/coreos/etcd v3.3.25+incompatible // indirect
	github.com/coreos/etcd/client v0.0.0-00010101000000-000000000000
	github.com/coreos/go-systemd v0.0.0-20141015001424-e3e4f602334e
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f
	github.com/godbus/dbus v0.0.0-20141007185835-25a4b8ca48c6
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/rkt/rkt v1.30.0
	golang.org/x/net v0.0.0-20201031054903-ff519b6c9102
)
