package upgrade

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/plum330/kratos/cmd/kratos/v2/internal/base"
)

// CmdUpgrade represents the upgrade command.
var CmdUpgrade = &cobra.Command{
	Use:   "upgrade",
	Short: "Upgrade the kratos tools",
	Long:  "Upgrade the kratos tools. Example: kratos upgrade",
	Run:   Run,
}

// Run upgrade the kratos tools.
func Run(_ *cobra.Command, _ []string) {
	plugins := []string{
		"github.com/plum330/kratos/cmd/protoc-gen-go-http/v2@latest",
		"github.com/plum330/kratos/cmd/protoc-gen-go-errors/v2@latest",
		"google.golang.org/protobuf/cmd/protoc-gen-go@latest",
		"google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest",
		"github.com/google/gnostic/cmd/protoc-gen-openapi@latest",
		"github.com/envoyproxy/protoc-gen-validate@latest",
		"github.com/favadi/protoc-go-inject-tag@latest",
		"github.com/google/wire/cmd/wire@latest",
		"github.com/rakyll/statik@latest",
		"github.com/mikefarah/yq/v4@latest",
	}
	err := base.GoInstall(plugins...)
	if err != nil {
		fmt.Println(err)
	}
}
