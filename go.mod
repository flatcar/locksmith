module github.com/flatcar/locksmith

go 1.14

require (
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f
	github.com/godbus/dbus v4.1.0+incompatible // indirect
	github.com/godbus/dbus/v5 v5.0.3
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/rkt/rkt v1.30.0
	go.etcd.io/etcd v0.0.0-00010101000000-000000000000
	go.uber.org/zap v1.16.0 // indirect
	golang.org/x/net v0.0.0-20201031054903-ff519b6c9102
	google.golang.org/grpc v1.33.2 // indirect
)

replace (
	// Force updating etcd to most recent version.
	go.etcd.io/etcd => go.etcd.io/etcd v0.5.0-alpha.5.0.20200824191128-ae9734ed278b
	// Most recent etcd version is not compatible with grpc v1.31.x.
	google.golang.org/grpc => google.golang.org/grpc v1.29.1
)
