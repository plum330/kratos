module github.com/plum330/kratos/contrib/registry/discovery/v2

go 1.19

require (
	github.com/go-resty/resty/v2 v2.11.0
	github.com/pkg/errors v0.9.1
	github.com/plum330/kratos/v2 v2.8.0
)

require golang.org/x/net v0.25.0 // indirect

replace github.com/plum330/kratos/v2 => ../../../
