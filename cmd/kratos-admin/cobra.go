package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/go-kratos/kratos/v2"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   Name,
		Short: "A kratos-layout service",
		Long: `A kratos-layout service. 
				Complete documentation is available at https://omalloc.com/kratos-layout`,
		Run: func(cmd *cobra.Command, args []string) {},
	}

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Print the version number",
		Long:  `All software has versions. My has one too.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf(`Version: %s --HEAD
GitCommit: %s
Kratos Version: %s
Go Version: %s
OS/Arch: %s/%s
`,
				Version,
				GitHash,
				kratos.Release,
				runtime.Version(),
				runtime.GOOS, runtime.GOARCH)
			os.Exit(0)
		},
	}
)
