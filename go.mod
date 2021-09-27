module github.com/flatcar-linux/locksmith

go 1.14

require (
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f
	github.com/flatcar-linux/fleetlock v0.0.0-20210922150917-05e572675abd
	github.com/godbus/dbus v4.1.0+incompatible // indirect
	github.com/godbus/dbus/v5 v5.0.4
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/rkt/rkt v1.30.0
	go.etcd.io/etcd v0.0.0-00010101000000-000000000000
	go.uber.org/tools v0.0.0-20190618225709-2cfd321de3ee // indirect
	golang.org/x/net v0.0.0-20210405180319-a5a99cb37ef4
)

replace (
	// Force updating etcd to most recent version.
	go.etcd.io/etcd => go.etcd.io/etcd v0.5.0-alpha.5.0.20200824191128-ae9734ed278b
	// Most recent etcd version is not compatible with grpc v1.31.x.
	google.golang.org/grpc => google.golang.org/grpc v1.29.1
)
