package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	command = &cobra.Command{
		Use:        "psn",
		Short:      "psn inspects Linux processes, sockets (ps, ss, netstat).",
		SuggestFor: []string{"pssn", "psns", "snp"},
	}
)

func init() {
	command.AddCommand(dsCommand)
	command.AddCommand(nsCommand)
}

func init() {
	cobra.EnablePrefixMatching = true
}

func main() {
	if err := command.Execute(); err != nil {
		fmt.Fprintln(os.Stdout, err)
		os.Exit(1)
	}
}
