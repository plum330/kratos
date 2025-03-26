module github.com/plum330/kratos/contrib/registry/zookeeper/v2

go 1.21

require (
	github.com/go-zookeeper/zk v1.0.3
	github.com/plum330/kratos/v2 v2.8.4
	golang.org/x/sync v0.8.0
)

replace github.com/plum330/kratos/v2 => ../../../
