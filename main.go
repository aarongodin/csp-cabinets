package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

func main() {
	var rootCmd = &cobra.Command{
		Use: "rlp",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("ok")
		},
	}

	rootCmd.AddCommand(calculateCabinets)

	rootCmd.Execute()
}
