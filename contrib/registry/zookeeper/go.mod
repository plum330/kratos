module github.com/plum330/kratos/contrib/registry/zookeeper/v2

go 1.19

require (
	github.com/go-zookeeper/zk v1.0.4
	github.com/plum330/kratos/v2 v2.8.0
	golang.org/x/sync v0.7.0
)

replace github.com/plum330/kratos/v2 => ../../../
