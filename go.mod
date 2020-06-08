module github.com/flatcar-linux/locksmith

replace github.com/coreos/bbolt => go.etcd.io/bbolt v1.3.4

go 1.14

require (
	github.com/coreos/etcd v3.3.22+incompatible
	github.com/coreos/go-semver v0.3.0 // indirect
	github.com/coreos/go-systemd v0.0.0-20191104093116-d3cd4ed1dbcf
	github.com/coreos/pkg v0.0.0-20180928190104-399ea9e2e55f
	github.com/godbus/dbus v4.1.0+incompatible
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/json-iterator/go v1.1.9 // indirect
	github.com/modern-go/reflect2 v1.0.1 // indirect
	github.com/rkt/rkt v1.30.0
	golang.org/x/net v0.0.0-20200602114024-627f9648deb9
)
