module github.com/go-kratos/kratos/contrib/server/kafka/v2

go 1.16

require (
	github.com/Shopify/sarama v1.31.1
	github.com/go-kratos/kratos/v2 v2.1.5
	golang.org/x/sync v0.0.0-20210220032951-036812b2e83c // indirect
)

replace github.com/go-kratos/kratos/v2 => ../../../
