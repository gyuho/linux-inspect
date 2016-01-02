package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	cliName        = "ssn"
	cliDescription = "ssn is an utility to investigate sockets."
)

// GlobalFlags contains all the flags defined globally
// and that are to be inherited to all sub-commands.
type GlobalFlags struct {
}

var (
	globalFlags = GlobalFlags{}

	rootCmd = &cobra.Command{
		Use:        cliName,
		Short:      cliDescription,
		SuggestFor: []string{"sssn", "sns", "sn"},
		RunE:       rootCommandFunc,
	}
)

func init() {
	// ...
}

func init() {
	cobra.EnablePrefixMatching = true
}

func rootCommandFunc(cmd *cobra.Command, args []string) error {
	fmt.Println("not ready")
	return nil
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
