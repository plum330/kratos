package main

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/plum330/kratos/cmd/kratos/v2/internal/change"
	"github.com/plum330/kratos/cmd/kratos/v2/internal/create"
	"github.com/plum330/kratos/cmd/kratos/v2/internal/project"
	"github.com/plum330/kratos/cmd/kratos/v2/internal/proto"
	"github.com/plum330/kratos/cmd/kratos/v2/internal/run"
	"github.com/plum330/kratos/cmd/kratos/v2/internal/upgrade"
	"github.com/plum330/kratos/cmd/kratos/v2/version"
)

var rootCmd = &cobra.Command{
	Use:     "kratos",
	Short:   "Kratos: An elegant toolkit for Go microservices.",
	Long:    `Kratos: An elegant toolkit for Go microservices.`,
	Version: version.Release,
}

func init() {
	rootCmd.AddCommand(project.CmdNew)
	rootCmd.AddCommand(proto.CmdProto)
	rootCmd.AddCommand(upgrade.CmdUpgrade)
	rootCmd.AddCommand(change.CmdChange)
	rootCmd.AddCommand(run.CmdRun)
	rootCmd.AddCommand(create.Cmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
