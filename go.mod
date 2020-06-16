module github.com/flatcar-linux/locksmith

go 1.14

replace github.com/coreos/etcd => ./etcd

replace github.com/coreos/etcd/client => ./etcd/client

require (
	github.com/coreos/etcd v2.2.5+incompatible // indirect
	github.com/coreos/etcd/client v0.0.0-00010101000000-000000000000
	github.com/coreos/go-systemd v0.0.0-20141015001424-e3e4f602334e
	github.com/coreos/pkg v0.0.0-20160210003529-549bd7890e35
	github.com/godbus/dbus v0.0.0-20141007185835-25a4b8ca48c6
	github.com/hashicorp/errwrap v1.0.0 // indirect
	github.com/rkt/rkt v1.27.0
	golang.org/x/net v0.0.0-20160201052856-d513e58596cd
)