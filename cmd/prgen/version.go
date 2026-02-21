package main

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/alann-estrada-KSH/ai-pr-generator/internal/version"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show the prgen version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(version.String())
		},
	}
}
