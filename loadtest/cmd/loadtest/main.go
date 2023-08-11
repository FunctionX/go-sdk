package main

import "github.com/functionx/go-sdk/loadtest"

func main() {
	rootCmd := loadtest.NewCmd()
	if err := rootCmd.Execute(); err != nil {
		rootCmd.PrintErrln(err)
	}
}
