package validate

import (
	"context"

	"github.com/bufbuild/protovalidate-go"

	"github.com/plum330/kratos/v2/errors"
	"github.com/plum330/kratos/v2/middleware"

	"google.golang.org/protobuf/proto"
)

type validator interface {
	Validate() error
	ValidateAll() error
}

// ProtoValidate is a validator middleware.
func ProtoValidate() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req any) (any, error) {
			if msg, ok := req.(proto.Message); ok {
				if err := protovalidate.Validate(msg); err != nil {
					return nil, errors.BadRequest("VALIDATOR", err.Error()).WithCause(err)
				}
			}

			// to compatible with the [old validator](https://github.com/envoyproxy/protoc-gen-validate)
			if v, ok := req.(validator); ok {
				if err := v.ValidateAll(); err != nil {
					return nil, errors.BadRequest("VALIDATOR", err.Error()).WithCause(err)
				}
			}
			return handler(ctx, req)
		}
	}
}
