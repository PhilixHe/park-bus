package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"park-bus/pkg/version"
)

const ErrExitCode = 1

func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "print version",
		RunE: func(cmd *cobra.Command, _ []string) error {
			fmt.Println(version.Get().String())
			return nil
		},
	}
}

func NewServiceCmd() *cobra.Command {
	return &cobra.Command{
		Use:     "service",
		Short:   "run service",
		Version: version.Get().String(),
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println("service running...")
			return nil
		},
	}
}

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "park-bus",
		Short:   "park-bus online binary",
		Version: version.Get().String(),
	}

	cmd.AddCommand(
		NewVersionCmd(),
		NewServiceCmd(),
	)

	return cmd
}

func main() {
	if err := NewRootCmd().Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(ErrExitCode)
	}
}
