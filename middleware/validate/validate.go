package validate

import (
	"context"
	"fmt"

	"github.com/bufbuild/protovalidate-go"

	"github.com/plum330/kratos/v2/errors"
	"github.com/plum330/kratos/v2/middleware"

	"google.golang.org/protobuf/proto"
)

type validator interface {
	Validate() error
	ValidateAll() error
}

var val *protovalidate.Validator

func init() {
	var err error
	val, err = protovalidate.New()
	if err != nil {
		panic(fmt.Sprintf("init protovalidator error:%+v", err))
	}
}

// Validator is a validator middleware.
func Validator() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (any, error) {
			var (
				ok  bool
				v   validator
				err error
				m   proto.Message
			)
			reason := "VALIDATOR"
			// to compatible with the [old validator](https://github.com/envoyproxy/protoc-gen-validate)
			if v, ok = req.(validator); ok {
				if err = v.ValidateAll(); err != nil {
					return nil, errors.BadRequest(reason, err.Error()).WithCause(err)
				}
			}

			if m, ok = req.(proto.Message); ok {
				if err = val.Validate(m); err != nil {
					return nil, errors.BadRequest(reason, err.Error()).WithCause(err)
				}
			}
			return handler(ctx, req)
		}
	}
}
