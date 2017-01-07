package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	Command = &cobra.Command{
		Use:        "psn",
		Short:      "psn inspects Linux processes, sockets (ps, ss, netstat).",
		SuggestFor: []string{"pssn", "psns", "snp"},
	}
)

func init() {
	Command.AddCommand(dsCommand)
}

func init() {
	cobra.EnablePrefixMatching = true
}

func main() {
	if err := Command.Execute(); err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
}
