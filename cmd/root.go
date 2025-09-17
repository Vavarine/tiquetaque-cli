package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "ttq",
	Short: "CLI para TiqueTaque",
	Long:  `Essa é uma CLI que facilita o registro de pontos no TiqueTaque.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Bem-vindo à CLI do TiqueTaque 🚀")
	},
}

// Função que o main.go chama
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
