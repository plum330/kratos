package validate

import (
	"context"
	"fmt"
	"strings"

	"google.golang.org/protobuf/proto"

	"github.com/bufbuild/protovalidate-go"

	"github.com/plum330/kratos/v2/errors"
	"github.com/plum330/kratos/v2/middleware"
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
			if v, ok := req.(validator); ok {
				if err := v.ValidateAll(); err != nil {
					return nil, errors.BadRequest("VALIDATOR", err.Error()).WithCause(err)
				}
			} else {
				if m, ok := req.(proto.Message); ok {
					if err := val.Validate(m); err != nil {
						es := err.Error()
						str := strings.Split(es, " ")
						// nolint:mnd
						if len(str) == 6 {
							return nil, errors.BadRequest(str[4], es).WithCause(err)
						}
						return nil, errors.BadRequest(es, es).WithCause(err)
					}
				}
			}
			return handler(ctx, req)
		}
	}
}
