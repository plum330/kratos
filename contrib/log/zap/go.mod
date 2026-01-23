module github.com/plum330/kratos/contrib/log/zap/v2

go 1.22

require (
	github.com/plum330/kratos/v2 v2.9.2
	go.uber.org/zap v1.26.0
)

require go.uber.org/multierr v1.11.0 // indirect

replace github.com/plum330/kratos/v2 => ../../../
