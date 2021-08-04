module github.com/flatcar-linux/locksmith

go 1.14

require (
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f
	github.com/godbus/dbus v4.1.0+incompatible // indirect
	github.com/godbus/dbus/v5 v5.0.4
	github.com/hashicorp/errwrap v1.1.0 // indirect
	github.com/rkt/rkt v1.30.0
	go.etcd.io/etcd/api/v3 v3.5.0
	go.etcd.io/etcd/client/v3 v3.5.0
	golang.org/x/net v0.0.0-20210405180319-a5a99cb37ef4
)

// Most recent etcd version is not compatible with grpc v1.31.x.
replace google.golang.org/grpc => google.golang.org/grpc v1.29.1
